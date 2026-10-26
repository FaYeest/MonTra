package database

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Status represents the operational health of the PostgreSQL connection.
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HealthChecker defines the contract used by HTTP health endpoints to probe database readiness.
type HealthChecker interface {
	CheckHealth(ctx context.Context) Status
	Required() bool
}

// Probe implements HealthChecker by verifying database configuration and TCP/socket reachability
// during the foundation phase (and wrapping pgxpool/sql.DB in Phase 1).
type Probe struct {
	dsn      string
	required bool
	timeout  time.Duration
	pingFunc func(ctx context.Context) error
}

// NewProbe creates a database health probe from the configured DSN.
func NewProbe(dsn string, required bool, timeout time.Duration) *Probe {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	p := &Probe{
		dsn:      dsn,
		required: required,
		timeout:  timeout,
	}
	p.pingFunc = p.tcpDialPing
	return p
}

// NewProbeWithPinger creates a Probe backed by an explicit ping function (such as sql.DB.PingContext or tests).
func NewProbeWithPinger(required bool, pingFunc func(ctx context.Context) error) *Probe {
	return &Probe{
		dsn:      "configured",
		required: required,
		timeout:  2 * time.Second,
		pingFunc: pingFunc,
	}
}

// NewPool initializes a pgxpool.Pool with sensible defaults.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, errors.New("empty database DSN")
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 25
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	return pgxpool.NewWithConfig(ctx, cfg)
}

// NewProbeWithPool creates a database health probe backed by a live pgxpool.Pool.
func NewProbeWithPool(pool *pgxpool.Pool, required bool) *Probe {
	return NewProbeWithPinger(required, func(ctx context.Context) error {
		if pool == nil {
			return errors.New("database pool is nil")
		}
		return pool.Ping(ctx)
	})
}

// Required reports whether a live database connection is mandatory in the current runtime mode.
func (p *Probe) Required() bool {
	return p != nil && p.required
}

// CheckHealth probes the database within a bounded context timeout.
func (p *Probe) CheckHealth(ctx context.Context) Status {
	if p == nil || p.dsn == "" {
		return Status{
			Status:  "unconfigured",
			Message: "Database connection parameters are not configured.",
		}
	}
	probeCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if err := p.pingFunc(probeCtx); err != nil {
		return Status{
			Status:  "unavailable",
			Message: "Database connection check failed.",
		}
	}
	return Status{
		Status: "ok",
	}
}

func (p *Probe) tcpDialPing(ctx context.Context) error {
	u, err := url.Parse(p.dsn)
	if err != nil || u.Host == "" {
		return err
	}
	host := u.Host
	if _, _, splitErr := net.SplitHostPort(host); splitErr != nil {
		host = net.JoinHostPort(host, "5432")
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", host)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
