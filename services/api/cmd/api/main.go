package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/montra/montra/services/api/internal/config"
	"github.com/montra/montra/services/api/internal/database"
	"github.com/montra/montra/services/api/internal/httpserver"
	"github.com/montra/montra/services/api/internal/mlclient"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	dbRequired := cfg.AppEnv == "production" || cfg.AppEnv == "staging"
	dbProbe := database.NewProbe(cfg.Database.DSN(), dbRequired, 2*time.Second)
	mlClient := mlclient.New(cfg.ML)

	handler := httpserver.NewHandler(httpserver.Dependencies{
		Config:   cfg,
		Logger:   logger,
		Database: dbProbe,
		MLClient: mlClient,
	})

	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting MonTra Go API server",
			slog.String("address", cfg.Address()),
			slog.String("environment", cfg.AppEnv),
			slog.Bool("ml_enabled", cfg.ML.Enabled),
		)
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Error("HTTP server stopped unexpectedly", slog.String("error", serveErr.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down MonTra Go API server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
