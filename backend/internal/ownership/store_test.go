
package ownership_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/arowden/augment-fund/internal/fund"
	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	tc, err := postgres.NewTestContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { tc.Cleanup(ctx) })

	store := ownership.NewStore(tc.Pool())
	fundStore := fund.NewStore(tc.Pool())

	createTestFund := func(t *testing.T, name string, units int) *fund.Fund {
		f, err := fund.NewFund(name, units)
		require.NoError(t, err)
		require.NoError(t, fundStore.Create(ctx, f))
		return f
	}

	t.Run("Create persists entry to database", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, err := ownership.NewCapTableEntry(testFund.ID, "John Doe", 500)
		require.NoError(t, err)

		err = store.Create(ctx, entry)
		require.NoError(t, err)

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "John Doe")
		require.NoError(t, err)
		assert.Equal(t, entry.ID, found.ID)
		assert.Equal(t, entry.OwnerName, found.OwnerName)
		assert.Equal(t, entry.Units, found.Units)
	})

	t.Run("Create returns ErrNilEntry for nil entry", func(t *testing.T) {
		tc.Reset(ctx)

		err := store.Create(ctx, nil)
		assert.ErrorIs(t, err, ownership.ErrNilEntry)
	})

	t.Run("Create returns error for duplicate owner in same fund", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry1, err := ownership.NewCapTableEntry(testFund.ID, "John Doe", 500)
		require.NoError(t, err)
		require.NoError(t, store.Create(ctx, entry1))

		entry2, err := ownership.NewCapTableEntry(testFund.ID, "John Doe", 300)
		require.NoError(t, err)
		err = store.Create(ctx, entry2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "John Doe")
	})

	t.Run("CreateTx uses provided transaction", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, err := ownership.NewCapTableEntry(testFund.ID, "Tx Owner", 250)
		require.NoError(t, err)

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = store.CreateTx(ctx, tx, entry)
		require.NoError(t, err)

		err = tx.Rollback(ctx)
		require.NoError(t, err)

		_, err = store.FindByFundAndOwner(ctx, testFund.ID, "Tx Owner")
		assert.True(t, errors.Is(err, ownership.ErrOwnerNotFound))
	})

	t.Run("CreateTx commits successfully", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, err := ownership.NewCapTableEntry(testFund.ID, "Committed Owner", 750)
		require.NoError(t, err)

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = store.CreateTx(ctx, tx, entry)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "Committed Owner")
		require.NoError(t, err)
		assert.Equal(t, entry.ID, found.ID)
	})

	t.Run("FindByFundID returns empty slice for missing fund", func(t *testing.T) {
		tc.Reset(ctx)
		nonExistentFundID := uuid.New()

		view, err := store.FindByFundID(ctx, nonExistentFundID, ownership.ListParams{})
		require.NoError(t, err)
		assert.NotNil(t, view.Entries)
		assert.Empty(t, view.Entries)
		assert.Equal(t, 0, view.TotalCount)
	})

	t.Run("FindByFundID returns entries ordered by units descending", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry1, _ := ownership.NewCapTableEntry(testFund.ID, "Small Owner", 100)
		require.NoError(t, store.Create(ctx, entry1))

		entry2, _ := ownership.NewCapTableEntry(testFund.ID, "Large Owner", 500)
		require.NoError(t, store.Create(ctx, entry2))

		entry3, _ := ownership.NewCapTableEntry(testFund.ID, "Medium Owner", 300)
		require.NoError(t, store.Create(ctx, entry3))

		view, err := store.FindByFundID(ctx, testFund.ID, ownership.ListParams{Limit: 10})
		require.NoError(t, err)
		require.Len(t, view.Entries, 3)
		assert.Equal(t, 3, view.TotalCount)

		assert.Equal(t, "Large Owner", view.Entries[0].OwnerName)
		assert.Equal(t, 500, view.Entries[0].Units)
		assert.Equal(t, "Medium Owner", view.Entries[1].OwnerName)
		assert.Equal(t, 300, view.Entries[1].Units)
		assert.Equal(t, "Small Owner", view.Entries[2].OwnerName)
		assert.Equal(t, 100, view.Entries[2].Units)
	})

	t.Run("FindByFundID pagination - limit", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		for i := 1; i <= 5; i++ {
			entry, _ := ownership.NewCapTableEntry(testFund.ID, "Owner "+string(rune('A'+i-1)), i*100)
			require.NoError(t, store.Create(ctx, entry))
		}

		view, err := store.FindByFundID(ctx, testFund.ID, ownership.ListParams{Limit: 2})
		require.NoError(t, err)
		assert.Len(t, view.Entries, 2)
		assert.Equal(t, 5, view.TotalCount)
		assert.Equal(t, 2, view.Limit)
	})

	t.Run("FindByFundID pagination - offset", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1500)

		entries := []struct {
			name  string
			units int
		}{
			{"Owner A", 100},
			{"Owner B", 200},
			{"Owner C", 300},
			{"Owner D", 400},
			{"Owner E", 500},
		}
		for _, e := range entries {
			entry, _ := ownership.NewCapTableEntry(testFund.ID, e.name, e.units)
			require.NoError(t, store.Create(ctx, entry))
		}

		view, err := store.FindByFundID(ctx, testFund.ID, ownership.ListParams{Limit: 2, Offset: 2})
		require.NoError(t, err)
		assert.Len(t, view.Entries, 2)
		assert.Equal(t, 5, view.TotalCount)
		assert.Equal(t, 2, view.Offset)

		assert.Equal(t, "Owner C", view.Entries[0].OwnerName)
		assert.Equal(t, "Owner B", view.Entries[1].OwnerName)
	})

	t.Run("FindByFundID TotalCount correct when offset exceeds total", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 100)

		entry, _ := ownership.NewCapTableEntry(testFund.ID, "Only Owner", 100)
		require.NoError(t, store.Create(ctx, entry))

		view, err := store.FindByFundID(ctx, testFund.ID, ownership.ListParams{Offset: 100})
		require.NoError(t, err)
		assert.Empty(t, view.Entries)
		assert.Equal(t, 1, view.TotalCount)
	})

	t.Run("FindByFundAndOwner returns ErrOwnerNotFound", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		_, err := store.FindByFundAndOwner(ctx, testFund.ID, "Nonexistent")
		assert.True(t, errors.Is(err, ownership.ErrOwnerNotFound))
	})

	t.Run("FindByFundAndOwner returns correct entry", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, _ := ownership.NewCapTableEntry(testFund.ID, "Specific Owner", 333)
		require.NoError(t, store.Create(ctx, entry))

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "Specific Owner")
		require.NoError(t, err)
		assert.Equal(t, entry.ID, found.ID)
		assert.Equal(t, "Specific Owner", found.OwnerName)
		assert.Equal(t, 333, found.Units)
	})

	t.Run("Upsert creates new entry", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, err := ownership.NewCapTableEntry(testFund.ID, "New Owner", 400)
		require.NoError(t, err)

		err = store.Upsert(ctx, entry)
		require.NoError(t, err)

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "New Owner")
		require.NoError(t, err)
		assert.Equal(t, 400, found.Units)
	})

	t.Run("Upsert updates existing entry", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		original, _ := ownership.NewCapTableEntry(testFund.ID, "Update Owner", 300)
		require.NoError(t, store.Create(ctx, original))

		time.Sleep(50 * time.Millisecond)

		updated, _ := ownership.NewCapTableEntry(testFund.ID, "Update Owner", 500)
		err := store.Upsert(ctx, updated)
		require.NoError(t, err)

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "Update Owner")
		require.NoError(t, err)
		assert.Equal(t, 500, found.Units)
		assert.Equal(t, original.ID, found.ID)
	})

	t.Run("Upsert preserves acquiredAt on update", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		original, _ := ownership.NewCapTableEntry(testFund.ID, "Acquired Owner", 200)
		require.NoError(t, store.Create(ctx, original))

		created, err := store.FindByFundAndOwner(ctx, testFund.ID, "Acquired Owner")
		require.NoError(t, err)
		originalAcquiredAt := created.AcquiredAt

		time.Sleep(50 * time.Millisecond)

		updated, _ := ownership.NewCapTableEntry(testFund.ID, "Acquired Owner", 400)
		err = store.Upsert(ctx, updated)
		require.NoError(t, err)

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "Acquired Owner")
		require.NoError(t, err)

		assert.Equal(t, originalAcquiredAt.Unix(), found.AcquiredAt.Unix())
		assert.True(t, found.UpdatedAt.After(originalAcquiredAt))
	})

	t.Run("Upsert returns ErrNilEntry for nil entry", func(t *testing.T) {
		tc.Reset(ctx)

		err := store.Upsert(ctx, nil)
		assert.ErrorIs(t, err, ownership.ErrNilEntry)
	})

	t.Run("UpsertTx with transaction rollback", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, err := ownership.NewCapTableEntry(testFund.ID, "Rollback Owner", 600)
		require.NoError(t, err)

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = store.UpsertTx(ctx, tx, entry)
		require.NoError(t, err)

		err = tx.Rollback(ctx)
		require.NoError(t, err)

		_, err = store.FindByFundAndOwner(ctx, testFund.ID, "Rollback Owner")
		assert.True(t, errors.Is(err, ownership.ErrOwnerNotFound))
	})

	t.Run("UpsertTx with transaction commit", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, err := ownership.NewCapTableEntry(testFund.ID, "Commit Owner", 700)
		require.NoError(t, err)

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = store.UpsertTx(ctx, tx, entry)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "Commit Owner")
		require.NoError(t, err)
		assert.Equal(t, 700, found.Units)
	})

	t.Run("concurrent upserts handle race condition", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Concurrent Fund", 10000)

		const numGoroutines = 10
		const unitsPerGoroutine = 100

		var wg sync.WaitGroup
		errChan := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				entry, err := ownership.NewCapTableEntry(testFund.ID, "Concurrent Owner", (id+1)*unitsPerGoroutine)
				if err != nil {
					errChan <- err
					return
				}
				if err := store.Upsert(ctx, entry); err != nil {
					errChan <- err
				}
			}(i)
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			t.Errorf("concurrent upsert error: %v", err)
		}

		found, err := store.FindByFundAndOwner(ctx, testFund.ID, "Concurrent Owner")
		require.NoError(t, err)
		assert.Greater(t, found.Units, 0)
	})

	t.Run("FindByFundID excludes soft-deleted entries", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, _ := ownership.NewCapTableEntry(testFund.ID, "Active Owner", 500)
		require.NoError(t, store.Create(ctx, entry))

		_, err := tc.Pool().Exec(ctx, `UPDATE cap_table_entries SET deleted_at = NOW() WHERE owner_name = $1`, "Active Owner")
		require.NoError(t, err)

		view, err := store.FindByFundID(ctx, testFund.ID, ownership.ListParams{})
		require.NoError(t, err)
		assert.Empty(t, view.Entries)
	})

	t.Run("FindByFundAndOwner excludes soft-deleted entries", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		entry, _ := ownership.NewCapTableEntry(testFund.ID, "Deleted Owner", 500)
		require.NoError(t, store.Create(ctx, entry))

		_, err := tc.Pool().Exec(ctx, `UPDATE cap_table_entries SET deleted_at = NOW() WHERE owner_name = $1`, "Deleted Owner")
		require.NoError(t, err)

		_, err = store.FindByFundAndOwner(ctx, testFund.ID, "Deleted Owner")
		assert.True(t, errors.Is(err, ownership.ErrOwnerNotFound))
	})

	t.Run("same owner name can exist in different funds", func(t *testing.T) {
		tc.Reset(ctx)
		fund1 := createTestFund(t, "Fund One", 1000)
		fund2 := createTestFund(t, "Fund Two", 2000)

		entry1, _ := ownership.NewCapTableEntry(fund1.ID, "Shared Owner", 100)
		require.NoError(t, store.Create(ctx, entry1))

		entry2, _ := ownership.NewCapTableEntry(fund2.ID, "Shared Owner", 200)
		require.NoError(t, store.Create(ctx, entry2))

		found1, err := store.FindByFundAndOwner(ctx, fund1.ID, "Shared Owner")
		require.NoError(t, err)
		assert.Equal(t, 100, found1.Units)

		found2, err := store.FindByFundAndOwner(ctx, fund2.ID, "Shared Owner")
		require.NoError(t, err)
		assert.Equal(t, 200, found2.Units)
	})

	t.Run("NewStore returns nil for nil db", func(t *testing.T) {
		store := ownership.NewStore(nil)
		assert.Nil(t, store)
	})
}
