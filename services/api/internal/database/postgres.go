package database

import (
	"context"
	"net"
	"net/url"
	"time"
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
