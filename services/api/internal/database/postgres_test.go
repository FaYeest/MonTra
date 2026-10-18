package database_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/montra/montra/services/api/internal/database"
)

func TestProbeUnconfigured(t *testing.T) {
	probe := database.NewProbe("", false, time.Second)
	status := probe.CheckHealth(context.Background())
	if status.Status != "unconfigured" {
		t.Fatalf("expected unconfigured status, got %q", status.Status)
	}
}

func TestProbeWithPinger(t *testing.T) {
	okProbe := database.NewProbeWithPinger(true, func(ctx context.Context) error {
		return nil
	})
	if st := okProbe.CheckHealth(context.Background()); st.Status != "ok" {
		t.Fatalf("expected ok status, got %+v", st)
	}

	failProbe := database.NewProbeWithPinger(true, func(ctx context.Context) error {
		return errors.New("connection refused")
	})
	if st := failProbe.CheckHealth(context.Background()); st.Status != "unavailable" {
		t.Fatalf("expected unavailable status, got %+v", st)
	}
}
