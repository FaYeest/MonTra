package mlclient_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/montra/montra/services/api/internal/config"
	"github.com/montra/montra/services/api/internal/mlclient"
)

func TestCheckHealthAndFallbackWhenUnreachable(t *testing.T) {
	client := mlclient.New(config.MLConfig{
		Enabled: true,
		BaseURL: "http://127.0.0.1:1",
		Timeout: 100 * time.Millisecond,
	})

	st := client.CheckHealth(context.Background())
	if st.Status != "unavailable" {
		t.Fatalf("expected unavailable status when ML server is down, got %+v", st)
	}
}

func TestPredictCategoryModelNotReadyAndConfidencePolicy(t *testing.T) {
	ready := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Internal-Service-Token") != "test-shared-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path == "/internal/v1/health" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"degraded","model_ready":false,"active_models":{"category_classifier":{"loaded":false,"model_version":"untrained-v0"}}}`))
			return
		}
		if r.URL.Path == "/internal/v1/predict/category" {
			if !ready {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"error":{"code":"MODEL_NOT_READY","message":"Model belum dimuat"}}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"predicted_category":"food_and_drink","confidence":0.82,"model_version":"category-v1"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := mlclient.New(config.MLConfig{
		Enabled:                     true,
		BaseURL:                     srv.URL,
		Timeout:                     time.Second,
		CategoryConfidenceThreshold: 0.75,
		InternalServiceToken:        "test-shared-token",
	})

	health := client.CheckHealth(context.Background())
	if health.Status != "degraded" || health.ModelVersion != "untrained-v0" {
		t.Fatalf("expected degraded health with untrained-v0, got %+v", health)
	}

	_, err := client.PredictCategory(context.Background(), mlclient.CategoryPredictionRequest{
		Description:     "Kopi Kenangan",
		TransactionType: "expense",
		Currency:        "IDR",
		Amount:          25000,
	})
	if !errors.Is(err, mlclient.ErrModelNotReady) {
		t.Fatalf("expected ErrModelNotReady before model artifact is loaded, got %v", err)
	}

	ready = true
	res, err := client.PredictCategory(context.Background(), mlclient.CategoryPredictionRequest{
		Description:     "Kopi Kenangan",
		TransactionType: "expense",
		Currency:        "IDR",
		Amount:          25000,
	})
	if err != nil {
		t.Fatalf("expected prediction to succeed when model is ready, got %v", err)
	}
	if !res.MeetsThreshold || res.PredictedCategory != "food_and_drink" {
		t.Fatalf("unexpected prediction result: %+v", res)
	}
}
