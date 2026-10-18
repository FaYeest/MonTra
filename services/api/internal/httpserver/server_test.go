package httpserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/montra/montra/services/api/internal/config"
	"github.com/montra/montra/services/api/internal/database"
	"github.com/montra/montra/services/api/internal/httpserver"
	"github.com/montra/montra/services/api/internal/mlclient"
)

type stubMLClient struct {
	status mlclient.HealthStatus
}

func (s stubMLClient) CheckHealth(_ context.Context) mlclient.HealthStatus {
	return s.status
}

func TestHealthEndpointRemainsCoreReadyWhenMLServiceIsDown(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	dbProbe := database.NewProbeWithPinger(true, func(ctx context.Context) error {
		return nil
	})
	mlStub := stubMLClient{
		status: mlclient.HealthStatus{
			Status:  "unavailable",
			Message: "ML container unreachable",
		},
	}

	handler := httpserver.NewHandler(httpserver.Dependencies{
		Config:   config.Config{AppEnv: "test"},
		Logger:   logger,
		Database: dbProbe,
		MLClient: mlStub,
	})

	for _, path := range []string{"/healthz", "/api/v1/health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected HTTP 200 when DB is healthy and ML is down, got %d", path, rec.Code)
		}

		var body httpserver.HealthResponse
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("%s: failed to decode response: %v", path, err)
		}
		if !body.CoreReady {
			t.Errorf("%s: expected CoreReady=true even when ML service is unavailable", path)
		}
		if body.Status != "degraded" {
			t.Errorf("%s: expected overall status 'degraded', got %q", path, body.Status)
		}
		if body.Components["database"].Status != "ok" {
			t.Errorf("%s: expected database status 'ok', got %q", path, body.Components["database"].Status)
		}
		if body.Components["ml_service"].Status != "unavailable" {
			t.Errorf("%s: expected ml_service status 'unavailable', got %q", path, body.Components["ml_service"].Status)
		}
	}
}

func TestHealthEndpointReturns503WhenRequiredDatabaseIsDown(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	dbProbe := database.NewProbeWithPinger(true, func(ctx context.Context) error {
		return errors.New("db down")
	})

	handler := httpserver.NewHandler(httpserver.Dependencies{
		Config:   config.Config{AppEnv: "production"},
		Logger:   logger,
		Database: dbProbe,
		MLClient: stubMLClient{status: mlclient.HealthStatus{Status: "ok", ModelReady: true}},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected HTTP 503 when required DB is down, got %d", rec.Code)
	}
}

func TestUnknownRouteReturnsJSONErrorEnvelope(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	handler := httpserver.NewHandler(httpserver.Dependencies{
		Config: config.Config{AppEnv: "test"},
		Logger: logger,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/does-not-exist", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected HTTP 404, got %d", rec.Code)
	}

	var errResp httpserver.ErrorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error envelope: %v", err)
	}
	if errResp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected error code NOT_FOUND, got %q", errResp.Error.Code)
	}
}
