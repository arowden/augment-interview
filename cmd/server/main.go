package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/arowden/augment-fund/internal/config"
	"github.com/arowden/augment-fund/internal/fund"
	apihttp "github.com/arowden/augment-fund/internal/http"
	"github.com/arowden/augment-fund/internal/otel"
	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/postgres"
	"github.com/arowden/augment-fund/internal/transfer"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if err := run(log); err != nil {
		log.Error("fatal error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	providers, err := otel.Init(ctx, cfg.Telemetry, log)
	if err != nil {
		return fmt.Errorf("initializing telemetry: %w", err)
	}
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := providers.Shutdown(shutdownCtx); err != nil {
			log.Error("failed to shutdown telemetry", slog.String("error", err.Error()))
		}
	}()

	if err := otel.InitMetrics(); err != nil {
		return fmt.Errorf("initializing metrics: %w", err)
	}

	pool, err := postgres.New(ctx, cfg.Database, log)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := postgres.Migrate(pool.Pool); err != nil {
		return err
	}
	log.Info("migrations complete")

	if err := postgres.RegisterMetrics(pool); err != nil {
		return err
	}

	fundStore := fund.NewStore(pool.Pool)
	ownershipStore := ownership.NewStore(pool.Pool)
	transferStore := transfer.NewStore(pool.Pool)

	fundSvc, err := fund.NewService(
		fundStore,
		fund.WithPool(pool.Pool),
		fund.WithOwnershipRepository(ownershipStore),
	)
	if err != nil {
		return fmt.Errorf("creating fund service: %w", err)
	}

	ownershipSvc, err := ownership.NewService(ownership.WithRepository(ownershipStore))
	if err != nil {
		return fmt.Errorf("creating ownership service: %w", err)
	}

	transferSvc, err := transfer.NewService(
		transfer.WithRepository(transferStore),
		transfer.WithOwnershipRepository(ownershipStore),
		transfer.WithPool(pool.Pool),
	)
	if err != nil {
		return fmt.Errorf("creating transfer service: %w", err)
	}

	apiHandler, err := apihttp.NewAPIHandlerStrict(
		apihttp.WithFundService(fundSvc),
		apihttp.WithOwnershipService(ownershipSvc),
		apihttp.WithTransferService(transferSvc),
		apihttp.WithPool(pool.Pool),
	)
	if err != nil {
		return fmt.Errorf("creating API handler: %w", err)
	}

	strictHandler := apihttp.NewStrictHandler(apiHandler, nil)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	apihttp.HandlerFromMux(strictHandler, r)

	handler := otel.WrapHandler(r, "augment-fund-api")

	addr := net.JoinHostPort(cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port))
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("starting HTTP server", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Info("shutting down server")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Info("server stopped")
	return nil
}
