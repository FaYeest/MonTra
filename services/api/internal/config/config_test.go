package config_test

import (
	"testing"
	"time"

	"github.com/montra/montra/services/api/internal/config"
)

func mapLookup(env map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := config.Load(mapLookup(nil))
	if err != nil {
		t.Fatalf("expected no error loading defaults, got: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.DefaultCurrency != "IDR" {
		t.Errorf("expected default currency IDR, got %s", cfg.DefaultCurrency)
	}
	if cfg.ML.Timeout != 1500*time.Millisecond {
		t.Errorf("expected default ML timeout 1500ms, got %v", cfg.ML.Timeout)
	}
	if cfg.Database.Configured() {
		t.Errorf("expected database to be unconfigured when no env vars are set")
	}
}

func TestLoadValidationErrors(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "invalid port",
			env:  map[string]string{"API_PORT": "99999"},
		},
		{
			name: "invalid confidence threshold",
			env:  map[string]string{"ML_CATEGORY_CONFIDENCE_THRESHOLD": "1.5"},
		},
		{
			name: "invalid currency code",
			env:  map[string]string{"DEFAULT_CURRENCY": "RUPIAH"},
		},
		{
			name: "invalid ML URL",
			env:  map[string]string{"ML_SERVICE_BASE_URL": "not-a-url"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := config.Load(mapLookup(tc.env)); err == nil {
				t.Fatalf("expected validation error for case %q, got nil", tc.name)
			}
		})
	}
}

func TestDatabaseDSNGeneration(t *testing.T) {
	cfg, err := config.Load(mapLookup(map[string]string{
		"POSTGRES_HOST":     "postgres",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_DB":       "montra_dev",
		"POSTGRES_USER":     "montra_app",
		"POSTGRES_PASSWORD": "placeholder_password",
		"POSTGRES_SSLMODE":  "disable",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Database.Configured() {
		t.Fatalf("expected database to be configured")
	}
	want := "postgres://montra_app:placeholder_password@postgres:5432/montra_dev?sslmode=disable"
	if got := cfg.Database.DSN(); got != want {
		t.Errorf("DSN mismatch: got %q, want %q", got, want)
	}
}
