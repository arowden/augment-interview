package transfer

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/arowden/augment-fund/internal/fund"
	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/postgres"
	"github.com/arowden/augment-fund/internal/validation"
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

	store := NewStore(tc.Pool())
	fundStore := fund.NewStore(tc.Pool())
	ownershipStore := ownership.NewStore(tc.Pool())

	// Helper to create a fund for testing.
	createTestFund := func(t *testing.T, name string, units int) *fund.Fund {
		f, err := fund.NewFund(name, units)
		require.NoError(t, err)
		require.NoError(t, fundStore.Create(ctx, f))
		return f
	}

	// Helper to create an ownership entry.
	createOwnership := func(t *testing.T, fundID uuid.UUID, owner string, units int) *ownership.Entry {
		entry, err := ownership.NewCapTableEntry(fundID, owner, units)
		require.NoError(t, err)
		require.NoError(t, ownershipStore.Create(ctx, entry))
		return entry
	}

	t.Run("Create persists transfer to database", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 500)

		transfer := &Transfer{
			ID:            uuid.New(),
			FundID:        testFund.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}

		err := store.Create(ctx, transfer)
		require.NoError(t, err)

		// Verify it was persisted.
		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: 10})
		require.NoError(t, err)
		require.Len(t, list.Transfers, 1)
		assert.Equal(t, transfer.ID, list.Transfers[0].ID)
		assert.Equal(t, "Alice", list.Transfers[0].FromOwner)
		assert.Equal(t, "Bob", list.Transfers[0].ToOwner)
		assert.Equal(t, 100, list.Transfers[0].Units)
	})

	t.Run("Create returns ErrNilTransfer for nil transfer", func(t *testing.T) {
		tc.Reset(ctx)

		err := store.Create(ctx, nil)
		assert.ErrorIs(t, err, ErrNilTransfer)
	})

	t.Run("Create with idempotency key", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		idempotencyKey := uuid.New()
		transfer := &Transfer{
			ID:             uuid.New(),
			FundID:         testFund.ID,
			FromOwner:      "Alice",
			ToOwner:        "Bob",
			Units:          100,
			IdempotencyKey: &idempotencyKey,
			TransferredAt:  time.Now(),
		}

		err := store.Create(ctx, transfer)
		require.NoError(t, err)

		// Verify idempotency key was stored.
		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{})
		require.NoError(t, err)
		require.Len(t, list.Transfers, 1)
		require.NotNil(t, list.Transfers[0].IdempotencyKey)
		assert.Equal(t, idempotencyKey, *list.Transfers[0].IdempotencyKey)
	})

	t.Run("CreateTx uses provided transaction", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		transfer := &Transfer{
			ID:            uuid.New(),
			FundID:        testFund.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = store.CreateTx(ctx, tx, transfer)
		require.NoError(t, err)

		// Rollback - transfer should not be persisted.
		err = tx.Rollback(ctx)
		require.NoError(t, err)

		// Verify not found.
		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{})
		require.NoError(t, err)
		assert.Empty(t, list.Transfers)
	})

	t.Run("CreateTx commits successfully", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		transfer := &Transfer{
			ID:            uuid.New(),
			FundID:        testFund.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)

		err = store.CreateTx(ctx, tx, transfer)
		require.NoError(t, err)

		err = tx.Commit(ctx)
		require.NoError(t, err)

		// Verify found.
		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{})
		require.NoError(t, err)
		require.Len(t, list.Transfers, 1)
		assert.Equal(t, transfer.ID, list.Transfers[0].ID)
	})

	t.Run("FindByFundID returns empty for missing fund", func(t *testing.T) {
		tc.Reset(ctx)
		nonExistentFundID := uuid.New()

		list, err := store.FindByFundID(ctx, nonExistentFundID, ListParams{})
		require.NoError(t, err)
		assert.NotNil(t, list.Transfers)
		assert.Empty(t, list.Transfers)
		assert.Equal(t, 0, list.TotalCount)
	})

	t.Run("FindByFundID returns transfers ordered by time and id ascending", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 300)
		createOwnership(t, testFund.ID, "Charlie", 200)

		// Create transfers - all will have same timestamp from database NOW().
		// Order is deterministic by (transferred_at, id).
		transfers := []*Transfer{
			{ID: uuid.New(), FundID: testFund.ID, FromOwner: "Alice", ToOwner: "Bob", Units: 100},
			{ID: uuid.New(), FundID: testFund.ID, FromOwner: "Bob", ToOwner: "Charlie", Units: 50},
			{ID: uuid.New(), FundID: testFund.ID, FromOwner: "Charlie", ToOwner: "Alice", Units: 25},
		}
		for _, tr := range transfers {
			require.NoError(t, store.Create(ctx, tr))
		}

		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: 10})
		require.NoError(t, err)
		require.Len(t, list.Transfers, 3)
		assert.Equal(t, 3, list.TotalCount)

		// Verify all transfers are returned (order is by timestamp then ID).
		returnedIDs := make(map[uuid.UUID]bool)
		for _, tr := range list.Transfers {
			returnedIDs[tr.ID] = true
		}
		for _, tr := range transfers {
			assert.True(t, returnedIDs[tr.ID], "transfer %s should be returned", tr.ID)
		}
	})

	t.Run("FindByFundID pagination - limit", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		// Create 5 transfers. Timestamps are database-generated (NOW()).
		for i := 0; i < 5; i++ {
			transfer := &Transfer{
				ID:        uuid.New(),
				FundID:    testFund.ID,
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     (i + 1) * 10,
			}
			require.NoError(t, store.Create(ctx, transfer))
		}

		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: 2})
		require.NoError(t, err)
		assert.Len(t, list.Transfers, 2)
		assert.Equal(t, 5, list.TotalCount)
		assert.Equal(t, 2, list.Limit)
	})

	t.Run("FindByFundID pagination - offset", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		// Create 5 transfers. Timestamps are database-generated (NOW()).
		transfers := make([]*Transfer, 5)
		for i := 0; i < 5; i++ {
			transfers[i] = &Transfer{
				ID:        uuid.New(),
				FundID:    testFund.ID,
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     (i + 1) * 10,
			}
			require.NoError(t, store.Create(ctx, transfers[i]))
		}

		// Get all to know the actual order.
		allList, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: 10})
		require.NoError(t, err)
		require.Len(t, allList.Transfers, 5)

		// Skip first 2, get next 2.
		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: 2, Offset: 2})
		require.NoError(t, err)
		assert.Len(t, list.Transfers, 2)
		assert.Equal(t, 5, list.TotalCount)
		assert.Equal(t, 2, list.Offset)

		// Should get 3rd and 4th entries from the full list.
		assert.Equal(t, allList.Transfers[2].ID, list.Transfers[0].ID)
		assert.Equal(t, allList.Transfers[3].ID, list.Transfers[1].ID)
	})

	t.Run("FindByFundID TotalCount correct when offset exceeds total", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		transfer := &Transfer{
			ID:            uuid.New(),
			FundID:        testFund.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}
		require.NoError(t, store.Create(ctx, transfer))

		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Offset: 100})
		require.NoError(t, err)
		assert.Empty(t, list.Transfers)
		assert.Equal(t, 1, list.TotalCount)
	})

	t.Run("FindByFundID normalizes invalid params", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		transfer := &Transfer{
			ID:            uuid.New(),
			FundID:        testFund.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}
		require.NoError(t, store.Create(ctx, transfer))

		// Negative values should be normalized.
		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: -1, Offset: -5})
		require.NoError(t, err)
		assert.Equal(t, validation.DefaultLimit, list.Limit)
		assert.Equal(t, 0, list.Offset)
	})

	t.Run("FindByFundID enforces max limit", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)

		transfer := &Transfer{
			ID:            uuid.New(),
			FundID:        testFund.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}
		require.NoError(t, store.Create(ctx, transfer))

		list, err := store.FindByFundID(ctx, testFund.ID, ListParams{Limit: 9999})
		require.NoError(t, err)
		assert.Equal(t, validation.MaxLimit, list.Limit)
	})

	t.Run("FindByIdempotencyKey returns nil when not found", func(t *testing.T) {
		tc.Reset(ctx)
		nonExistentKey := uuid.New()

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)
		defer tx.Rollback(ctx) //nolint:errcheck

		found, err := store.FindByIdempotencyKey(ctx, tx, nonExistentKey)
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("FindByIdempotencyKey returns existing transfer", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 100)
		idempotencyKey := uuid.New()

		transfer := &Transfer{
			ID:             uuid.New(),
			FundID:         testFund.ID,
			FromOwner:      "Alice",
			ToOwner:        "Bob",
			Units:          100,
			IdempotencyKey: &idempotencyKey,
			TransferredAt:  time.Now(),
		}
		require.NoError(t, store.Create(ctx, transfer))

		tx, err := tc.Pool().Begin(ctx)
		require.NoError(t, err)
		defer tx.Rollback(ctx) //nolint:errcheck

		found, err := store.FindByIdempotencyKey(ctx, tx, idempotencyKey)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, transfer.ID, found.ID)
		assert.Equal(t, transfer.FromOwner, found.FromOwner)
		assert.Equal(t, transfer.ToOwner, found.ToOwner)
		assert.Equal(t, transfer.Units, found.Units)
	})

	t.Run("NewStore returns nil for nil db", func(t *testing.T) {
		store := NewStore(nil)
		assert.Nil(t, store)
	})

	t.Run("transfers from different funds are isolated", func(t *testing.T) {
		tc.Reset(ctx)
		fund1 := createTestFund(t, "Fund One", 1000)
		fund2 := createTestFund(t, "Fund Two", 2000)
		createOwnership(t, fund1.ID, "Alice", 500)
		createOwnership(t, fund1.ID, "Bob", 100)
		createOwnership(t, fund2.ID, "Charlie", 500)
		createOwnership(t, fund2.ID, "Dave", 100)

		transfer1 := &Transfer{
			ID:            uuid.New(),
			FundID:        fund1.ID,
			FromOwner:     "Alice",
			ToOwner:       "Bob",
			Units:         100,
			TransferredAt: time.Now(),
		}
		require.NoError(t, store.Create(ctx, transfer1))

		transfer2 := &Transfer{
			ID:            uuid.New(),
			FundID:        fund2.ID,
			FromOwner:     "Charlie",
			ToOwner:       "Dave",
			Units:         200,
			TransferredAt: time.Now(),
		}
		require.NoError(t, store.Create(ctx, transfer2))

		// Each fund should only see its own transfers.
		list1, err := store.FindByFundID(ctx, fund1.ID, ListParams{})
		require.NoError(t, err)
		require.Len(t, list1.Transfers, 1)
		assert.Equal(t, transfer1.ID, list1.Transfers[0].ID)

		list2, err := store.FindByFundID(ctx, fund2.ID, ListParams{})
		require.NoError(t, err)
		require.Len(t, list2.Transfers, 1)
		assert.Equal(t, transfer2.ID, list2.Transfers[0].ID)
	})
}

func TestTransferService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	tc, err := postgres.NewTestContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { tc.Cleanup(ctx) })

	fundStore := fund.NewStore(tc.Pool())
	ownershipStore := ownership.NewStore(tc.Pool())
	transferStore := NewStore(tc.Pool())

	// Helper to create a fund for testing.
	createTestFund := func(t *testing.T, name string, units int) *fund.Fund {
		f, err := fund.NewFund(name, units)
		require.NoError(t, err)
		require.NoError(t, fundStore.Create(ctx, f))
		return f
	}

	// Helper to create an ownership entry.
	createOwnership := func(t *testing.T, fundID uuid.UUID, owner string, units int) *ownership.Entry {
		entry, err := ownership.NewCapTableEntry(fundID, owner, units)
		require.NoError(t, err)
		require.NoError(t, ownershipStore.Create(ctx, entry))
		return entry
	}

	t.Run("ExecuteTransfer successful transfer", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		req := Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}

		transfer, err := svc.ExecuteTransfer(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, transfer)
		assert.Equal(t, "Alice", transfer.FromOwner)
		assert.Equal(t, "Bob", transfer.ToOwner)
		assert.Equal(t, 100, transfer.Units)

		// Verify Alice's units decreased.
		aliceEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, 400, aliceEntry.Units)

		// Verify Bob was created with correct units.
		bobEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Bob")
		require.NoError(t, err)
		assert.Equal(t, 100, bobEntry.Units)
	})

	t.Run("ExecuteTransfer to existing owner adds units", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)
		createOwnership(t, testFund.ID, "Bob", 200)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		req := Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}

		transfer, err := svc.ExecuteTransfer(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, transfer)

		// Verify Alice's units decreased.
		aliceEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, 400, aliceEntry.Units)

		// Verify Bob's units increased.
		bobEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Bob")
		require.NoError(t, err)
		assert.Equal(t, 300, bobEntry.Units)
	})

	t.Run("ExecuteTransfer fails for non-existent from_owner", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		req := Request{
			FundID:    testFund.ID,
			FromOwner: "NonExistent",
			ToOwner:   "Bob",
			Units:     100,
		}

		_, err = svc.ExecuteTransfer(ctx, req)
		assert.ErrorIs(t, err, ErrOwnerNotFound)
	})

	t.Run("ExecuteTransfer fails for insufficient units", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 50)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		req := Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}

		_, err = svc.ExecuteTransfer(ctx, req)
		assert.ErrorIs(t, err, ErrInsufficientUnits)

		// Verify Alice's units unchanged.
		aliceEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, 50, aliceEntry.Units)
	})

	t.Run("ExecuteTransfer idempotency returns existing transfer", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		idempotencyKey := uuid.New()
		req := Request{
			FundID:         testFund.ID,
			FromOwner:      "Alice",
			ToOwner:        "Bob",
			Units:          100,
			IdempotencyKey: &idempotencyKey,
		}

		// First transfer.
		transfer1, err := svc.ExecuteTransfer(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, transfer1)

		// Second transfer with same idempotency key.
		transfer2, err := svc.ExecuteTransfer(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, transfer2)

		// Should return the same transfer.
		assert.Equal(t, transfer1.ID, transfer2.ID)

		// Alice should only have lost 100 units total.
		aliceEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, 400, aliceEntry.Units)

		// Bob should only have 100 units.
		bobEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Bob")
		require.NoError(t, err)
		assert.Equal(t, 100, bobEntry.Units)
	})

	t.Run("ExecuteTransfer validation errors", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		// Self-transfer.
		_, err = svc.ExecuteTransfer(ctx, Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Alice",
			Units:     100,
		})
		assert.ErrorIs(t, err, ErrSelfTransfer)

		// Zero units.
		_, err = svc.ExecuteTransfer(ctx, Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     0,
		})
		assert.ErrorIs(t, err, ErrInvalidUnits)

		// Empty owner.
		_, err = svc.ExecuteTransfer(ctx, Request{
			FundID:    testFund.ID,
			FromOwner: "",
			ToOwner:   "Bob",
			Units:     100,
		})
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("ExecuteTransfer concurrent transfers use pessimistic locking", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 10000)
		createOwnership(t, testFund.ID, "Alice", 1000)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		const numGoroutines = 10
		const unitsPerTransfer = 100

		var wg sync.WaitGroup
		successCount := 0
		var mu sync.Mutex

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				req := Request{
					FundID:    testFund.ID,
					FromOwner: "Alice",
					ToOwner:   "Bob",
					Units:     unitsPerTransfer,
				}
				_, err := svc.ExecuteTransfer(ctx, req)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// All 10 transfers should succeed since Alice has 1000 units.
		assert.Equal(t, numGoroutines, successCount)

		// Verify final balances.
		aliceEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, 0, aliceEntry.Units)

		bobEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Bob")
		require.NoError(t, err)
		assert.Equal(t, 1000, bobEntry.Units)
	})

	t.Run("ExecuteTransfer concurrent overdraft protection", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 10000)
		createOwnership(t, testFund.ID, "Alice", 100) // Only 100 units

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		const numGoroutines = 5
		const unitsPerTransfer = 50 // Each wants 50, but only 100 available

		var wg sync.WaitGroup
		successCount := 0
		var mu sync.Mutex

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				req := Request{
					FundID:    testFund.ID,
					FromOwner: "Alice",
					ToOwner:   "Bob",
					Units:     unitsPerTransfer,
				}
				_, err := svc.ExecuteTransfer(ctx, req)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Only 2 transfers should succeed (100 / 50 = 2).
		assert.Equal(t, 2, successCount)

		// Verify final balances.
		aliceEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, 0, aliceEntry.Units)

		bobEntry, err := ownershipStore.FindByFundAndOwner(ctx, testFund.ID, "Bob")
		require.NoError(t, err)
		assert.Equal(t, 100, bobEntry.Units)
	})

	t.Run("ListTransfers returns transfer history", func(t *testing.T) {
		tc.Reset(ctx)
		testFund := createTestFund(t, "Test Fund", 1000)
		createOwnership(t, testFund.ID, "Alice", 500)

		svc, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		require.NoError(t, err)

		// Execute multiple transfers.
		_, err = svc.ExecuteTransfer(ctx, Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		})
		require.NoError(t, err)

		_, err = svc.ExecuteTransfer(ctx, Request{
			FundID:    testFund.ID,
			FromOwner: "Alice",
			ToOwner:   "Charlie",
			Units:     50,
		})
		require.NoError(t, err)

		// List transfers.
		list, err := svc.ListTransfers(ctx, testFund.ID, ListParams{})
		require.NoError(t, err)
		assert.Len(t, list.Transfers, 2)
		assert.Equal(t, 2, list.TotalCount)
	})

	t.Run("NewService fails without repository", func(t *testing.T) {
		_, err := NewService(
			WithOwnershipRepository(ownershipStore),
			WithPool(tc.Pool()),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository is required")
	})

	t.Run("NewService fails without ownership repository", func(t *testing.T) {
		_, err := NewService(
			WithRepository(transferStore),
			WithPool(tc.Pool()),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ownership repository is required")
	})

	t.Run("NewService fails without pool", func(t *testing.T) {
		_, err := NewService(
			WithRepository(transferStore),
			WithOwnershipRepository(ownershipStore),
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pool is required")
	})
}
