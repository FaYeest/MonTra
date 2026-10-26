package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds validated runtime configuration for the MonTra Go API server.
type Config struct {
	AppEnv          string
	LogLevel        string
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	DefaultCurrency string
	DefaultTimezone string

	Auth     AuthConfig
	Database DatabaseConfig
	ML       MLConfig
}

// AuthConfig holds authentication and token issuance parameters.
type AuthConfig struct {
	JWTSecret string
	JWTExpiry time.Duration
}

// DatabaseConfig holds PostgreSQL connection parameters.
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	Explicit string
}

// DSN returns a PostgreSQL connection URL suitable for pgx / database/sql.
func (d DatabaseConfig) DSN() string {
	if strings.TrimSpace(d.Explicit) != "" {
		return strings.TrimSpace(d.Explicit)
	}
	if strings.TrimSpace(d.Host) == "" || strings.TrimSpace(d.Name) == "" || strings.TrimSpace(d.User) == "" {
		return ""
	}
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   d.Name,
	}
	q := u.Query()
	sslMode := d.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()
	return u.String()
}

// Configured reports whether database connection details have been provided.
func (d DatabaseConfig) Configured() bool {
	return d.DSN() != ""
}

// MLConfig holds internal Python FastAPI service settings.
type MLConfig struct {
	Enabled                     bool
	BaseURL                     string
	Timeout                     time.Duration
	CategoryConfidenceThreshold float64
	InternalServiceToken        string
}

// LoadFromEnv reads and validates configuration from process environment variables.
func LoadFromEnv() (Config, error) {
	return Load(os.LookupEnv)
}

// Load reads and validates configuration using the provided lookup function.
func Load(lookup func(string) (string, bool)) (Config, error) {
	port, err := parseInt(lookup, "API_PORT", 8080)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, errors.New("API_PORT must be a valid TCP port (1-65535)")
	}

	readTimeoutMs, err := parseInt(lookup, "API_READ_TIMEOUT_MS", 5000)
	if err != nil || readTimeoutMs <= 0 {
		return Config{}, errors.New("API_READ_TIMEOUT_MS must be a positive integer")
	}

	writeTimeoutMs, err := parseInt(lookup, "API_WRITE_TIMEOUT_MS", 10000)
	if err != nil || writeTimeoutMs <= 0 {
		return Config{}, errors.New("API_WRITE_TIMEOUT_MS must be a positive integer")
	}

	shutdownTimeoutMs, err := parseInt(lookup, "API_SHUTDOWN_TIMEOUT_MS", 10000)
	if err != nil || shutdownTimeoutMs <= 0 {
		return Config{}, errors.New("API_SHUTDOWN_TIMEOUT_MS must be a positive integer")
	}

	dbPort, err := parseInt(lookup, "POSTGRES_PORT", 5432)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return Config{}, errors.New("POSTGRES_PORT must be a valid TCP port (1-65535)")
	}

	mlEnabled, err := parseBool(lookup, "ML_ENABLED", true)
	if err != nil {
		return Config{}, errors.New("ML_ENABLED must be a boolean value")
	}

	mlTimeoutMs, err := parseInt(lookup, "ML_CLIENT_TIMEOUT_MS", 1500)
	if err != nil || mlTimeoutMs <= 0 {
		return Config{}, errors.New("ML_CLIENT_TIMEOUT_MS must be a positive integer")
	}

	threshold, err := parseFloat(lookup, "ML_CATEGORY_CONFIDENCE_THRESHOLD", 0.75)
	if err != nil || threshold <= 0.0 || threshold > 1.0 {
		return Config{}, errors.New("ML_CATEGORY_CONFIDENCE_THRESHOLD must be in (0.0, 1.0]")
	}

	currency := strings.ToUpper(getString(lookup, "DEFAULT_CURRENCY", "IDR"))
	if len(currency) != 3 {
		return Config{}, errors.New("DEFAULT_CURRENCY must be a 3-letter ISO 4217 currency code")
	}

	mlBaseURL := strings.TrimRight(getString(lookup, "ML_SERVICE_BASE_URL", "http://localhost:8000"), "/")
	if mlEnabled {
		parsedURL, err := url.Parse(mlBaseURL)
		if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
			return Config{}, errors.New("ML_SERVICE_BASE_URL must be a valid HTTP/HTTPS URL when ML_ENABLED is true")
		}
	}

	appEnv := getString(lookup, "APP_ENV", "development")
	jwtSecret := getString(lookup, "JWT_SECRET", "CHANGE_ME_LOCAL_DEV_JWT_SECRET_32_BYTES_MINIMUM")
	jwtExpiryHours, err := parseInt(lookup, "JWT_EXPIRY_HOURS", 72)
	if err != nil || jwtExpiryHours <= 0 {
		return Config{}, errors.New("JWT_EXPIRY_HOURS must be a positive integer")
	}
	if (appEnv == "production" || appEnv == "staging") && len(jwtSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must be at least 32 characters in production or staging")
	}

	return Config{
		AppEnv:          appEnv,
		LogLevel:        strings.ToUpper(getString(lookup, "LOG_LEVEL", "INFO")),
		Host:            getString(lookup, "API_HOST", "0.0.0.0"),
		Port:            port,
		ReadTimeout:     time.Duration(readTimeoutMs) * time.Millisecond,
		WriteTimeout:    time.Duration(writeTimeoutMs) * time.Millisecond,
		ShutdownTimeout: time.Duration(shutdownTimeoutMs) * time.Millisecond,
		DefaultCurrency: currency,
		DefaultTimezone: getString(lookup, "DEFAULT_TIMEZONE", "Asia/Jakarta"),
		Auth: AuthConfig{
			JWTSecret: jwtSecret,
			JWTExpiry: time.Duration(jwtExpiryHours) * time.Hour,
		},
		Database: DatabaseConfig{
			Host:     getString(lookup, "POSTGRES_HOST", ""),
			Port:     dbPort,
			Name:     getString(lookup, "POSTGRES_DB", ""),
			User:     getString(lookup, "POSTGRES_USER", ""),
			Password: getString(lookup, "POSTGRES_PASSWORD", ""),
			SSLMode:  getString(lookup, "POSTGRES_SSLMODE", "disable"),
			Explicit: getString(lookup, "DATABASE_URL", ""),
		},
		ML: MLConfig{
			Enabled:                     mlEnabled,
			BaseURL:                     mlBaseURL,
			Timeout:                     time.Duration(mlTimeoutMs) * time.Millisecond,
			CategoryConfidenceThreshold: threshold,
			InternalServiceToken:        getString(lookup, "INTERNAL_SERVICE_TOKEN", ""),
		},
	}, nil
}

// Address returns the host:port listen string.
func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getString(lookup func(string) (string, bool), key, fallback string) string {
	val, ok := lookup(key)
	if !ok || strings.TrimSpace(val) == "" {
		return fallback
	}
	return strings.TrimSpace(val)
}

func parseInt(lookup func(string) (string, bool), key string, fallback int) (int, error) {
	raw, ok := lookup(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	return strconv.Atoi(strings.TrimSpace(raw))
}

func parseBool(lookup func(string) (string, bool), key string, fallback bool) (bool, error) {
	raw, ok := lookup(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	return strconv.ParseBool(strings.TrimSpace(raw))
}

func parseFloat(lookup func(string) (string, bool), key string, fallback float64) (float64, error) {
	raw, ok := lookup(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	return strconv.ParseFloat(strings.TrimSpace(raw), 64)
}
