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

	// Create stores.
	fundStore := fund.NewStore(pool.Pool)
	ownershipStore := ownership.NewStore(pool.Pool)
	transferStore := transfer.NewStore(pool.Pool)

	// Create services.
	// Fund service has pool and ownership repo for transactional fund creation.
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

	// Create API handler with strict validation (fail fast if services missing).
	apiHandler, err := apihttp.NewAPIHandlerStrict(
		apihttp.WithFundService(fundSvc),
		apihttp.WithOwnershipService(ownershipSvc),
		apihttp.WithTransferService(transferSvc),
	)
	if err != nil {
		return fmt.Errorf("creating API handler: %w", err)
	}

	// Create strict handler wrapper.
	strictHandler := apihttp.NewStrictHandler(apiHandler, nil)

	// Create router.
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint.
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Mount API routes.
	apihttp.HandlerFromMux(strictHandler, r)

	// Create HTTP server.
	addr := net.JoinHostPort(cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port))
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine.
	errCh := make(chan error, 1)
	go func() {
		log.Info("starting HTTP server", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or error.
	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Info("shutting down server")
	}

	// Graceful shutdown.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Info("server stopped")
	return nil
}
