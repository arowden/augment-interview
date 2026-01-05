package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(pool *pgxpool.Pool) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("migrate: failed to create source: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate: failed to create driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate: failed to create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: failed to run migrations: %w", err)
	}

	return nil
}

func MigrateDown(pool *pgxpool.Pool) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("migrate: failed to create source: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate: failed to create driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate: failed to create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: failed to rollback migrations: %w", err)
	}

	return nil
}

func MigrateVersion(pool *pgxpool.Pool) (uint, bool, error) {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return 0, false, fmt.Errorf("migrate: failed to create source: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("migrate: failed to create driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return 0, false, fmt.Errorf("migrate: failed to create migrator: %w", err)
	}
	defer m.Close()

	return m.Version()
}
