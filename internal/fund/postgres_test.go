package fund

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arowden/augment-fund/internal/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	tc, err := postgres.NewTestContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { tc.Cleanup(ctx) })

	repo := NewPostgresRepository(tc.Pool())

	t.Run("Create persists fund to database", func(t *testing.T) {
		tc.Reset(ctx)
		fund, err := NewFund("Test Fund", 1000)
		require.NoError(t, err)

		err = repo.Create(ctx, fund)
		require.NoError(t, err)

		// Verify it was persisted
		found, err := repo.FindByID(ctx, fund.ID)
		require.NoError(t, err)
		assert.Equal(t, fund.ID, found.ID)
		assert.Equal(t, fund.Name, found.Name)
		assert.Equal(t, fund.TotalUnits, found.TotalUnits)
	})

	t.Run("Create returns ErrNilFund for nil fund", func(t *testing.T) {
		tc.Reset(ctx)

		err := repo.Create(ctx, nil)
		assert.ErrorIs(t, err, ErrNilFund)
	})

	t.Run("Create returns ErrDuplicateFundName for duplicate names", func(t *testing.T) {
		tc.Reset(ctx)
		fund1, err := NewFund("Unique Fund", 1000)
		require.NoError(t, err)
		require.NoError(t, repo.Create(ctx, fund1))

		fund2, err := NewFund("Unique Fund", 2000)
		require.NoError(t, err)
		err = repo.Create(ctx, fund2)
		assert.ErrorIs(t, err, ErrDuplicateFundName)
	})

	t.Run("CreateTx uses provided transaction", func(t *testing.T) {
		tc.Reset(ctx)
		fund, err := NewFund("Tx Fund", 500)
		require.NoError(t, err)

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = repo.CreateTx(ctx, tx, fund)
		require.NoError(t, err)

		// Rollback - fund should not be persisted
		err = tx.Rollback(ctx)
		require.NoError(t, err)

		// Verify not found
		_, err = repo.FindByID(ctx, fund.ID)
		assert.True(t, errors.Is(err, ErrNotFound))
	})

	t.Run("CreateTx commits successfully", func(t *testing.T) {
		tc.Reset(ctx)
		fund, err := NewFund("Committed Fund", 750)
		require.NoError(t, err)

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = repo.CreateTx(ctx, tx, fund)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		// Verify found
		found, err := repo.FindByID(ctx, fund.ID)
		require.NoError(t, err)
		assert.Equal(t, fund.ID, found.ID)
	})

	t.Run("FindByID returns ErrNotFound for missing fund", func(t *testing.T) {
		tc.Reset(ctx)
		nonExistentID := uuid.New()

		_, err := repo.FindByID(ctx, nonExistentID)
		assert.True(t, errors.Is(err, ErrNotFound))
	})

	t.Run("List returns empty result when no funds", func(t *testing.T) {
		tc.Reset(ctx)

		result, err := repo.List(ctx, ListParams{})
		require.NoError(t, err)
		assert.NotNil(t, result.Items)
		assert.Empty(t, result.Items)
		assert.Equal(t, 0, result.Total)
		assert.Equal(t, DefaultListLimit, result.Limit)
		assert.Equal(t, 0, result.Offset)
	})

	t.Run("List returns all funds with correct total", func(t *testing.T) {
		tc.Reset(ctx)

		// Create funds with delays to ensure distinct timestamps.
		// NewFund sets CreatedAt to time.Now(), so we must sleep before each call.
		fund1, _ := NewFund("First Fund", 100)
		require.NoError(t, repo.Create(ctx, fund1))

		time.Sleep(50 * time.Millisecond)
		fund2, _ := NewFund("Second Fund", 200)
		require.NoError(t, repo.Create(ctx, fund2))

		time.Sleep(50 * time.Millisecond)
		fund3, _ := NewFund("Third Fund", 300)
		require.NoError(t, repo.Create(ctx, fund3))

		result, err := repo.List(ctx, ListParams{Limit: 10})
		require.NoError(t, err)
		require.Len(t, result.Items, 3)
		assert.Equal(t, 3, result.Total)

		// Verify all funds are present (order is newest first)
		ids := make([]uuid.UUID, len(result.Items))
		for i, f := range result.Items {
			ids[i] = f.ID
		}
		assert.Contains(t, ids, fund1.ID)
		assert.Contains(t, ids, fund2.ID)
		assert.Contains(t, ids, fund3.ID)

		// Verify ordering: most recent first
		assert.Equal(t, fund3.ID, result.Items[0].ID)
		assert.Equal(t, fund2.ID, result.Items[1].ID)
		assert.Equal(t, fund1.ID, result.Items[2].ID)
	})

	t.Run("List respects limit parameter", func(t *testing.T) {
		tc.Reset(ctx)

		// Create 5 funds with distinct timestamps.
		for i := 1; i <= 5; i++ {
			fund, _ := NewFund("Fund "+string(rune('A'+i-1)), i*100)
			require.NoError(t, repo.Create(ctx, fund))
			time.Sleep(20 * time.Millisecond)
		}

		result, err := repo.List(ctx, ListParams{Limit: 2})
		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, 5, result.Total)
		assert.Equal(t, 2, result.Limit)
	})

	t.Run("List respects offset parameter", func(t *testing.T) {
		tc.Reset(ctx)

		// Create 5 funds with distinct timestamps.
		funds := make([]*Fund, 5)
		for i := 0; i < 5; i++ {
			fund, _ := NewFund("Fund "+string(rune('A'+i)), (i+1)*100)
			funds[i] = fund
			require.NoError(t, repo.Create(ctx, fund))
			time.Sleep(20 * time.Millisecond)
		}

		// Skip first 2 (newest), get next 2.
		result, err := repo.List(ctx, ListParams{Limit: 2, Offset: 2})
		require.NoError(t, err)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, 5, result.Total)
		assert.Equal(t, 2, result.Offset)

		// Should get funds[2] and funds[1] (3rd and 2nd newest).
		assert.Equal(t, funds[2].ID, result.Items[0].ID)
		assert.Equal(t, funds[1].ID, result.Items[1].ID)
	})

	t.Run("List returns empty when offset exceeds total", func(t *testing.T) {
		tc.Reset(ctx)

		fund, _ := NewFund("Only Fund", 100)
		require.NoError(t, repo.Create(ctx, fund))

		result, err := repo.List(ctx, ListParams{Offset: 100})
		require.NoError(t, err)
		assert.Empty(t, result.Items)
		assert.Equal(t, 1, result.Total)
	})

	t.Run("List normalizes invalid params", func(t *testing.T) {
		tc.Reset(ctx)

		fund, _ := NewFund("Test Fund", 100)
		require.NoError(t, repo.Create(ctx, fund))

		// Negative values should be normalized
		result, err := repo.List(ctx, ListParams{Limit: -1, Offset: -5})
		require.NoError(t, err)
		assert.Equal(t, DefaultListLimit, result.Limit)
		assert.Equal(t, 0, result.Offset)
	})

	t.Run("List enforces max limit", func(t *testing.T) {
		tc.Reset(ctx)

		fund, _ := NewFund("Test Fund", 100)
		require.NoError(t, repo.Create(ctx, fund))

		result, err := repo.List(ctx, ListParams{Limit: 9999})
		require.NoError(t, err)
		assert.Equal(t, MaxListLimit, result.Limit)
	})
}
