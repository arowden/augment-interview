package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	cfg Config
	log *slog.Logger
}

func New(ctx context.Context, cfg Config, log *slog.Logger) (*Pool, error) {
	dsn := cfg.DSN()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to parse DSN: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	log.Info("connecting to database",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("database", cfg.DBName),
		slog.String("user", cfg.User),
		slog.String("sslmode", cfg.SSLMode),
	)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: failed to ping database: %w", err)
	}

	log.Info("database connection established",
		slog.Int("max_conns", int(cfg.MaxConns)),
		slog.Int("min_conns", int(cfg.MinConns)),
	)

	return &Pool{
		Pool: pool,
		cfg:  cfg,
		log:  log,
	}, nil
}

func (p *Pool) Config() Config {
	return p.cfg
}

func (p *Pool) HealthCheck(ctx context.Context) error {
	return p.Ping(ctx)
}
