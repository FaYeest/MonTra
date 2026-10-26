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

	"github.com/montra/montra/services/api/internal/accounts"
	"github.com/montra/montra/services/api/internal/auth"
	"github.com/montra/montra/services/api/internal/categories"
	"github.com/montra/montra/services/api/internal/config"
	"github.com/montra/montra/services/api/internal/database"
	"github.com/montra/montra/services/api/internal/httpserver"
	"github.com/montra/montra/services/api/internal/mlclient"
	"github.com/montra/montra/services/api/internal/transactions"
	"github.com/montra/montra/services/api/internal/transfers"
	"github.com/montra/montra/services/api/internal/users"
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
	var dbProbe database.HealthChecker
	var authHandler *auth.Handler
	var categoriesHandler *categories.Handler
	var accountsHandler *accounts.Handler
	var transactionsHandler *transactions.Handler
	var transfersHandler *transfers.Handler

	if cfg.Database.Configured() {
		poolCtx, poolCancel := context.WithTimeout(context.Background(), 5*time.Second)
		dbPool, err := database.NewPool(poolCtx, cfg.Database.DSN())
		poolCancel()
		if err != nil {
			logger.Warn("could not initialize database connection pool", slog.String("error", err.Error()))
			dbProbe = database.NewProbe(cfg.Database.DSN(), dbRequired, 2*time.Second)
		} else {
			defer dbPool.Close()
			dbProbe = database.NewProbeWithPool(dbPool, dbRequired)

			userRepo := users.NewPgRepository(dbPool)
			hasher := auth.NewHasher()
			tokenMgr, err := auth.NewTokenManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry)
			if err != nil {
				logger.Error("failed to initialize token manager", slog.String("error", err.Error()))
				os.Exit(1)
			}
			authSvc := auth.NewService(userRepo, hasher, tokenMgr)
			authHandler = auth.NewHandler(authSvc, tokenMgr)

			catRepo := categories.NewPgRepository(dbPool)
			catSvc := categories.NewService(catRepo)
			categoriesHandler = categories.NewHandler(catSvc, tokenMgr)

			accRepo := accounts.NewPgRepository(dbPool)
			accSvc := accounts.NewService(accRepo)
			accountsHandler = accounts.NewHandler(accSvc, tokenMgr)

			txRepo := transactions.NewPgRepository(dbPool)
			txSvc := transactions.NewService(txRepo)
			transactionsHandler = transactions.NewHandler(txSvc, tokenMgr)

			trRepo := transfers.NewPgRepository(dbPool)
			trSvc := transfers.NewService(trRepo)
			transfersHandler = transfers.NewHandler(trSvc, tokenMgr)
		}
	} else {
		dbProbe = database.NewProbe("", dbRequired, 2*time.Second)
	}

	mlClient := mlclient.New(cfg.ML)

	handler := httpserver.NewHandler(httpserver.Dependencies{
		Config:              cfg,
		Logger:              logger,
		Database:            dbProbe,
		MLClient:            mlClient,
		AuthHandler:         authHandler,
		CategoriesHandler:   categoriesHandler,
		AccountsHandler:     accountsHandler,
		TransactionsHandler: transactionsHandler,
		TransfersHandler:    transfersHandler,
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
