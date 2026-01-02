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

// TestContainer wraps a PostgreSQL test container.
type TestContainer struct {
	container testcontainers.Container
	pool      *pgxpool.Pool
	cfg       Config
}

// NewTestContainer creates a new PostgreSQL container for testing.
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
		pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		pgContainer.Terminate(ctx)
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
		pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	tc := &TestContainer{
		container: pgContainer,
		pool:      pool,
		cfg:       cfg,
	}

	if err := tc.Migrate(ctx); err != nil {
		tc.Cleanup(ctx)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return tc, nil
}

// Pool returns the connection pool for the test container.
func (tc *TestContainer) Pool() *pgxpool.Pool {
	return tc.pool
}

// Config returns the configuration for the test container.
func (tc *TestContainer) Config() Config {
	return tc.cfg
}

// Migrate runs all database migrations.
func (tc *TestContainer) Migrate(_ context.Context) error {
	return Migrate(tc.pool)
}

// Cleanup terminates the container and closes the pool.
func (tc *TestContainer) Cleanup(ctx context.Context) error {
	if tc.pool != nil {
		tc.pool.Close()
	}
	if tc.container != nil {
		return tc.container.Terminate(ctx)
	}
	return nil
}

// Reset truncates all tables except schema_migrations.
func (tc *TestContainer) Reset(ctx context.Context) error {
	_, err := tc.pool.Exec(ctx, `
		TRUNCATE TABLE transfers, cap_table_entries, funds CASCADE
	`)
	return err
}
