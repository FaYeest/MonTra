package mlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/montra/montra/services/api/internal/config"
)

var (
	// ErrMLDisabled indicates that AI/ML features are disabled by configuration.
	ErrMLDisabled = errors.New("mlclient: ML service is disabled")
	// ErrMLUnavailable indicates that the internal Python ML service is unreachable or returned an error.
	ErrMLUnavailable = errors.New("mlclient: ML service is unavailable")
	// ErrModelNotReady indicates that the ML service is reachable but no evaluated model artifact is loaded.
	ErrModelNotReady = errors.New("mlclient: requested ML model artifact is not loaded")
)

// HealthStatus represents the health and model readiness of the internal Python service.
type HealthStatus struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	ModelVersion string `json:"model_version,omitempty"`
	ModelReady   bool   `json:"model_ready"`
}

// CategoryPredictionRequest is the minimal feature payload sent to POST /internal/v1/predict/category.
// Sensitive user identifiers, emails, and tokens are never included.
type CategoryPredictionRequest struct {
	Description     string `json:"description"`
	TransactionType string `json:"transaction_type"`
	Currency        string `json:"currency"`
	Amount          int64  `json:"amount"`
}

// CategoryPredictionResponse is the response returned by the Python ML service.
type CategoryPredictionResponse struct {
	PredictedCategory string  `json:"predicted_category"`
	Confidence        float64 `json:"confidence"`
	ModelVersion      string  `json:"model_version"`
	MeetsThreshold    bool    `json:"meets_threshold"`
}

// Client communicates with the internal Python FastAPI service over HTTP/JSON with strict timeouts.
type Client struct {
	cfg        config.MLConfig
	httpClient *http.Client
}

// New creates a configured ML Client.
func New(cfg config.MLConfig) *Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 1500 * time.Millisecond
	}
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// CheckHealth probes GET /internal/v1/health on the Python ML service.
// Failures are reported gracefully without panicking or failing core financial endpoints.
func (c *Client) CheckHealth(ctx context.Context) HealthStatus {
	if c == nil || !c.cfg.Enabled {
		return HealthStatus{
			Status:  "disabled",
			Message: "AI/ML service integration is disabled via configuration.",
		}
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/internal/v1/health"
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return HealthStatus{
			Status:  "unavailable",
			Message: "Invalid ML service endpoint configuration.",
		}
	}
	c.attachInternalHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return HealthStatus{
			Status:  "unavailable",
			Message: "Internal ML service is unreachable; core financial features remain active.",
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthStatus{
			Status:  "unavailable",
			Message: fmt.Sprintf("Internal ML service returned HTTP %d.", resp.StatusCode),
		}
	}

	var payload struct {
		Status       string `json:"status"`
		ModelReady   bool   `json:"model_ready"`
		ActiveModels struct {
			CategoryClassifier struct {
				Loaded       bool   `json:"loaded"`
				ModelVersion string `json:"model_version"`
			} `json:"category_classifier"`
		} `json:"active_models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return HealthStatus{
			Status:  "degraded",
			Message: "Internal ML service returned an unreadable health payload.",
		}
	}

	status := "ok"
	msg := ""
	if !payload.ModelReady {
		status = "degraded"
		msg = "ML service reachable; trained model artifacts are not yet loaded."
	}

	return HealthStatus{
		Status:       status,
		Message:      msg,
		ModelVersion: payload.ActiveModels.CategoryClassifier.ModelVersion,
		ModelReady:   payload.ModelReady,
	}
}

// PredictCategory calls POST /internal/v1/predict/category and evaluates the confidence policy.
func (c *Client) PredictCategory(ctx context.Context, input CategoryPredictionRequest) (CategoryPredictionResponse, error) {
	if c == nil || !c.cfg.Enabled {
		return CategoryPredictionResponse{}, ErrMLDisabled
	}

	body, err := json.Marshal(input)
	if err != nil {
		return CategoryPredictionResponse{}, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/internal/v1/predict/category"
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return CategoryPredictionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.attachInternalHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CategoryPredictionResponse{}, fmt.Errorf("%w: %v", ErrMLUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return CategoryPredictionResponse{}, ErrModelNotReady
	}
	if resp.StatusCode != http.StatusOK {
		return CategoryPredictionResponse{}, fmt.Errorf("%w: status %d", ErrMLUnavailable, resp.StatusCode)
	}

	var result CategoryPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CategoryPredictionResponse{}, fmt.Errorf("%w: invalid response JSON", ErrMLUnavailable)
	}
	result.MeetsThreshold = result.Confidence >= c.cfg.CategoryConfidenceThreshold
	return result, nil
}

func (c *Client) attachInternalHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	if strings.TrimSpace(c.cfg.InternalServiceToken) != "" {
		req.Header.Set("X-Internal-Service-Token", c.cfg.InternalServiceToken)
	}
}
