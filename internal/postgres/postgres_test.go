package postgres_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/arowden/augment-fund/internal/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testContainer *postgres.TestContainer

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = postgres.NewTestContainer(ctx)
	if err != nil {
		panic("failed to create test container: " + err.Error())
	}

	code := m.Run()

	testContainer.Cleanup(ctx)

	if code != 0 {
		panic("tests failed")
	}
}

func TestFKConstraintEnforcement(t *testing.T) {
	ctx := context.Background()
	pool := testContainer.Pool()

	err := testContainer.Reset(ctx)
	require.NoError(t, err)

	var fundID uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO funds (name, total_units)
		VALUES ('Test Fund', 1000)
		RETURNING id
	`).Scan(&fundID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO cap_table_entries (fund_id, owner_name, units)
		VALUES ($1, 'Alice', 500), ($1, 'Bob', 500)
	`, fundID)
	require.NoError(t, err)

	t.Run("valid transfer with existing owners succeeds", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units)
			VALUES ($1, 'Alice', 'Bob', 100)
		`, fundID)
		assert.NoError(t, err)
	})

	t.Run("transfer with non-existent from_owner fails", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units)
			VALUES ($1, 'NonExistent', 'Bob', 100)
		`, fundID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "violates foreign key constraint")
	})

	t.Run("transfer with non-existent to_owner fails", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units)
			VALUES ($1, 'Alice', 'NonExistent', 100)
		`, fundID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "violates foreign key constraint")
	})

	t.Run("self-transfer fails with check constraint", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units)
			VALUES ($1, 'Alice', 'Alice', 100)
		`, fundID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chk_different_owners")
	})
}

func TestAutomaticUpdatedAtTrigger(t *testing.T) {
	ctx := context.Background()
	pool := testContainer.Pool()

	err := testContainer.Reset(ctx)
	require.NoError(t, err)

	var fundID uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO funds (name, total_units)
		VALUES ('Trigger Test Fund', 1000)
		RETURNING id
	`).Scan(&fundID)
	require.NoError(t, err)

	var entryID uuid.UUID
	var createdAt, updatedAt time.Time
	err = pool.QueryRow(ctx, `
		INSERT INTO cap_table_entries (fund_id, owner_name, units)
		VALUES ($1, 'Alice', 500)
		RETURNING id, acquired_at, updated_at
	`, fundID).Scan(&entryID, &createdAt, &updatedAt)
	require.NoError(t, err)

	assert.WithinDuration(t, createdAt, updatedAt, time.Second)

	time.Sleep(100 * time.Millisecond)

	var newUpdatedAt time.Time
	err = pool.QueryRow(ctx, `
		UPDATE cap_table_entries
		SET units = 400
		WHERE id = $1
		RETURNING updated_at
	`, entryID).Scan(&newUpdatedAt)
	require.NoError(t, err)

	assert.True(t, newUpdatedAt.After(updatedAt),
		"updated_at should be after original: got %v, want after %v", newUpdatedAt, updatedAt)
}

func TestIdempotencyKeyUniqueness(t *testing.T) {
	ctx := context.Background()
	pool := testContainer.Pool()

	err := testContainer.Reset(ctx)
	require.NoError(t, err)

	var fundID uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO funds (name, total_units)
		VALUES ('Idempotency Test Fund', 1000)
		RETURNING id
	`).Scan(&fundID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO cap_table_entries (fund_id, owner_name, units)
		VALUES ($1, 'Alice', 500), ($1, 'Bob', 500)
	`, fundID)
	require.NoError(t, err)

	idempotencyKey := uuid.New()

	t.Run("first transfer with idempotency key succeeds", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units, idempotency_key)
			VALUES ($1, 'Alice', 'Bob', 100, $2)
		`, fundID, idempotencyKey)
		assert.NoError(t, err)
	})

	t.Run("duplicate idempotency key fails", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units, idempotency_key)
			VALUES ($1, 'Alice', 'Bob', 100, $2)
		`, fundID, idempotencyKey)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
	})

	t.Run("null idempotency keys are allowed multiple times", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units, idempotency_key)
			VALUES ($1, 'Alice', 'Bob', 50, NULL)
		`, fundID)
		assert.NoError(t, err)

		_, err = pool.Exec(ctx, `
			INSERT INTO transfers (fund_id, from_owner, to_owner, units, idempotency_key)
			VALUES ($1, 'Alice', 'Bob', 50, NULL)
		`, fundID)
		assert.NoError(t, err)
	})
}

func TestMigrations(t *testing.T) {
	ctx := context.Background()
	pool := testContainer.Pool()

	t.Run("migrations applied", func(t *testing.T) {
		version, dirty, err := postgres.MigrateVersion(pool)
		require.NoError(t, err)
		assert.False(t, dirty)
		assert.EqualValues(t, 6, version)
	})

	t.Run("funds table exists", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables
				WHERE table_name = 'funds'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("cap_table_entries table has soft delete column", func(t *testing.T) {
		var exists bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name = 'cap_table_entries' AND column_name = 'deleted_at'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("indexes are created", func(t *testing.T) {
		var count int
		err := pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM pg_indexes
			WHERE tablename IN ('cap_table_entries', 'transfers')
			AND indexname LIKE 'idx_%'
		`).Scan(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 7, "expected at least 7 custom indexes")
	})
}

func TestPoolHealthCheck(t *testing.T) {
	ctx := context.Background()
	log := testLogger()

	cfg := testContainer.Config()
	cfg.MaxConns = 5
	cfg.MinConns = 1

	pool, err := postgres.New(ctx, cfg, log)
	require.NoError(t, err)
	defer pool.Close()

	err = pool.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestRegisterMetrics(t *testing.T) {
	ctx := context.Background()
	log := testLogger()

	cfg := testContainer.Config()
	cfg.MaxConns = 5
	cfg.MinConns = 1

	pool, err := postgres.New(ctx, cfg, log)
	require.NoError(t, err)
	defer pool.Close()

	err = postgres.RegisterMetrics(pool)
	assert.NoError(t, err)
}
