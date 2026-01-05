package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestContainer struct {
	container testcontainers.Container
	pool      *pgxpool.Pool
	cfg       Config
}

func NewTestContainer(ctx context.Context) (*TestContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	cfg := Config{
		Host:            host,
		Port:            port.Int(),
		User:            "test",
		Password:        "test",
		DBName:          "test_db",
		SSLMode:         "disable",
		MaxConns:        10,
		MinConns:        1,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 10 * time.Minute,
	}

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	tc := &TestContainer{
		container: pgContainer,
		pool:      pool,
		cfg:       cfg,
	}

	if err := tc.Migrate(ctx); err != nil {
		_ = tc.Cleanup(ctx)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return tc, nil
}

func (tc *TestContainer) Pool() *pgxpool.Pool {
	return tc.pool
}

func (tc *TestContainer) Config() Config {
	return tc.cfg
}

func (tc *TestContainer) Migrate(_ context.Context) error {
	return Migrate(tc.pool)
}

func (tc *TestContainer) Cleanup(ctx context.Context) error {
	if tc.pool != nil {
		tc.pool.Close()
	}
	if tc.container != nil {
		return tc.container.Terminate(ctx)
	}
	return nil
}

func (tc *TestContainer) Reset(ctx context.Context) error {
	_, err := tc.pool.Exec(ctx, `
		TRUNCATE TABLE transfers, cap_table_entries, funds CASCADE
	`)
	return err
}
