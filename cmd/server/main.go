package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/arowden/augment-fund/internal/config"
	"github.com/arowden/augment-fund/internal/postgres"
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

	log.Info("configuration loaded",
		slog.String("db_host", cfg.Database.Host),
		slog.String("db_name", cfg.Database.DBName),
		slog.String("server_host", cfg.Server.Host),
		slog.Int("server_port", cfg.Server.Port),
	)

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

	log.Info("server starting",
		slog.String("host", cfg.Server.Host),
		slog.Int("port", cfg.Server.Port),
	)

	<-ctx.Done()
	log.Info("shutting down")

	return nil
}
