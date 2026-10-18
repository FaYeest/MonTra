package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/montra/montra/services/api/internal/config"
	"github.com/montra/montra/services/api/internal/database"
	"github.com/montra/montra/services/api/internal/mlclient"
)

// MLHealthChecker abstracts the internal ML service health probe.
type MLHealthChecker interface {
	CheckHealth(ctx context.Context) mlclient.HealthStatus
}

// ComponentHealth represents the status of an individual backend dependency in the public health response.
type ComponentHealth struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	ModelVersion string `json:"model_version,omitempty"`
}

// HealthResponse matches the PublicHealthResponse schema in packages/api-spec/openapi.yaml.
type HealthResponse struct {
	Status      string                     `json:"status"`
	Service     string                     `json:"service"`
	Environment string                     `json:"environment"`
	CoreReady   bool                       `json:"core_ready"`
	Timestamp   string                     `json:"timestamp"`
	Components  map[string]ComponentHealth `json:"components"`
}

// ErrorEnvelope is the standard public API error payload.
type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail provides a machine-readable code and user-safe message without leaking internal stack traces.
type ErrorDetail struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// Dependencies bundles runtime collaborators for the HTTP server.
type Dependencies struct {
	Config   config.Config
	Logger   *slog.Logger
	Database database.HealthChecker
	MLClient MLHealthChecker
}

// NewHandler builds the HTTP router and middleware stack for the Go API server.
func NewHandler(deps Dependencies) http.Handler {
	logger := deps.Logger
	if logger == nil {
		logger = slog.Default()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth(deps))
	mux.HandleFunc("GET /api/v1/health", handleHealth(deps))

	// Fallback handler returning a consistent JSON 404 envelope.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, ErrorEnvelope{
			Error: ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "Endpoint tidak ditemukan.",
			},
		})
	})

	return withSecurityAndLogging(logger, mux)
}

func handleHealth(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dbComponent := ComponentHealth{Status: "unconfigured"}
		coreReady := true
		if deps.Database != nil {
			dbStatus := deps.Database.CheckHealth(ctx)
			dbComponent = ComponentHealth{
				Status:  dbStatus.Status,
				Message: dbStatus.Message,
			}
			if deps.Database.Required() && dbStatus.Status != "ok" {
				coreReady = false
			}
		}

		mlComponent := ComponentHealth{Status: "disabled"}
		if deps.MLClient != nil {
			mlStatus := deps.MLClient.CheckHealth(ctx)
			mlComponent = ComponentHealth{
				Status:       mlStatus.Status,
				Message:      mlStatus.Message,
				ModelVersion: mlStatus.ModelVersion,
			}
		}

		// Architectural invariant: ML service unavailability or untrained model status
		// degrades only assistive AI features; it never sets core_ready to false.
		overallStatus := "ok"
		statusCode := http.StatusOK
		if !coreReady {
			overallStatus = "unavailable"
			statusCode = http.StatusServiceUnavailable
		} else if dbComponent.Status == "unavailable" || mlComponent.Status == "unavailable" || mlComponent.Status == "degraded" {
			overallStatus = "degraded"
		}

		resp := HealthResponse{
			Status:      overallStatus,
			Service:     "montra-api",
			Environment: deps.Config.AppEnv,
			CoreReady:   coreReady,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Components: map[string]ComponentHealth{
				"database":   dbComponent,
				"ml_service": mlComponent,
			},
		}

		writeJSON(w, statusCode, resp)
	}
}

func withSecurityAndLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")

		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("recovered from panic in HTTP handler",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
				)
				writeJSON(rec, http.StatusInternalServerError, ErrorEnvelope{
					Error: ErrorDetail{
						Code:    "INTERNAL_ERROR",
						Message: "Terjadi kesalahan pada server.",
					},
				})
			}
			logger.Info("http_request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rec.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		}()

		next.ServeHTTP(rec, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
