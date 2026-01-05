package ownership

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	createFunc                      func(ctx context.Context, entry *Entry) error
	createTxFunc                    func(ctx context.Context, tx pgx.Tx, entry *Entry) error
	findByFundIDFunc                func(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error)
	findByFundAndOwnerFunc          func(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error)
	findByFundAndOwnerForUpdateFunc func(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string) (*Entry, error)
	decrementUnitsTxFunc            func(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, units int) error
	incrementOrCreateTxFunc         func(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string, units int) error
	upsertFunc                      func(ctx context.Context, entry *Entry) error
	upsertTxFunc                    func(ctx context.Context, tx pgx.Tx, entry *Entry) error
}

func (m *mockRepository) Create(ctx context.Context, entry *Entry) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, entry)
	}
	return nil
}

func (m *mockRepository) CreateTx(ctx context.Context, tx pgx.Tx, entry *Entry) error {
	if m.createTxFunc != nil {
		return m.createTxFunc(ctx, tx, entry)
	}
	return nil
}

func (m *mockRepository) FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error) {
	if m.findByFundIDFunc != nil {
		return m.findByFundIDFunc(ctx, fundID, params)
	}
	return &CapTableView{FundID: fundID, Entries: []*Entry{}}, nil
}

func (m *mockRepository) FindByFundAndOwner(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error) {
	if m.findByFundAndOwnerFunc != nil {
		return m.findByFundAndOwnerFunc(ctx, fundID, ownerName)
	}
	return nil, OwnerNotFoundError(fundID, ownerName)
}

func (m *mockRepository) FindByFundAndOwnerForUpdateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string) (*Entry, error) {
	if m.findByFundAndOwnerForUpdateFunc != nil {
		return m.findByFundAndOwnerForUpdateFunc(ctx, tx, fundID, ownerName)
	}
	return nil, OwnerNotFoundError(fundID, ownerName)
}

func (m *mockRepository) DecrementUnitsTx(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, units int) error {
	if m.decrementUnitsTxFunc != nil {
		return m.decrementUnitsTxFunc(ctx, tx, entryID, units)
	}
	return nil
}

func (m *mockRepository) IncrementOrCreateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string, units int) error {
	if m.incrementOrCreateTxFunc != nil {
		return m.incrementOrCreateTxFunc(ctx, tx, fundID, ownerName, units)
	}
	return nil
}

func (m *mockRepository) Upsert(ctx context.Context, entry *Entry) error {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, entry)
	}
	return nil
}

func (m *mockRepository) UpsertTx(ctx context.Context, tx pgx.Tx, entry *Entry) error {
	if m.upsertTxFunc != nil {
		return m.upsertTxFunc(ctx, tx, entry)
	}
	return nil
}

func TestNewService(t *testing.T) {
	t.Run("returns error when no repository is configured", func(t *testing.T) {
		svc, err := NewService()
		assert.Nil(t, svc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository is required")
	})

	t.Run("creates service with repository", func(t *testing.T) {
		repo := &mockRepository{}
		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)
		assert.NotNil(t, svc)
	})
}

func TestService_GetCapTable(t *testing.T) {
	fundID := uuid.New()

	t.Run("returns cap table entries", func(t *testing.T) {
		expectedView := &CapTableView{
			FundID: fundID,
			Entries: []*Entry{
				{ID: uuid.New(), FundID: fundID, OwnerName: "Alice", Units: 500},
				{ID: uuid.New(), FundID: fundID, OwnerName: "Bob", Units: 300},
			},
			TotalCount: 2,
			Limit:      100,
			Offset:     0,
		}
		repo := &mockRepository{
			findByFundIDFunc: func(ctx context.Context, fID uuid.UUID, params ListParams) (*CapTableView, error) {
				return expectedView, nil
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		view, err := svc.GetCapTable(context.Background(), fundID, ListParams{})
		require.NoError(t, err)
		assert.Equal(t, expectedView, view)
	})

	t.Run("passes params to repository", func(t *testing.T) {
		var receivedParams ListParams
		repo := &mockRepository{
			findByFundIDFunc: func(ctx context.Context, fID uuid.UUID, params ListParams) (*CapTableView, error) {
				receivedParams = params
				return &CapTableView{FundID: fID, Entries: []*Entry{}}, nil
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		params := ListParams{Limit: 50, Offset: 10}
		_, err = svc.GetCapTable(context.Background(), fundID, params)
		require.NoError(t, err)
		assert.Equal(t, params, receivedParams)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		repo := &mockRepository{
			findByFundIDFunc: func(ctx context.Context, fID uuid.UUID, params ListParams) (*CapTableView, error) {
				return nil, repoErr
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		view, err := svc.GetCapTable(context.Background(), fundID, ListParams{})
		assert.Nil(t, view)
		assert.Equal(t, repoErr, err)
	})
}

func TestService_GetOwnership(t *testing.T) {
	fundID := uuid.New()

	t.Run("returns entry when found", func(t *testing.T) {
		expectedEntry := &Entry{
			ID:        uuid.New(),
			FundID:    fundID,
			OwnerName: "Alice",
			Units:     500,
		}
		repo := &mockRepository{
			findByFundAndOwnerFunc: func(ctx context.Context, fID uuid.UUID, ownerName string) (*Entry, error) {
				return expectedEntry, nil
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		entry, err := svc.GetOwnership(context.Background(), fundID, "Alice")
		require.NoError(t, err)
		assert.Equal(t, expectedEntry, entry)
	})

	t.Run("returns error when owner not found", func(t *testing.T) {
		repo := &mockRepository{
			findByFundAndOwnerFunc: func(ctx context.Context, fID uuid.UUID, ownerName string) (*Entry, error) {
				return nil, OwnerNotFoundError(fID, ownerName)
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		entry, err := svc.GetOwnership(context.Background(), fundID, "NonExistent")
		assert.Nil(t, entry)
		assert.True(t, errors.Is(err, ErrOwnerNotFound))
	})
}

func TestService_CreateEntry(t *testing.T) {
	fundID := uuid.New()

	t.Run("creates entry with valid inputs", func(t *testing.T) {
		var createdEntry *Entry
		repo := &mockRepository{
			createFunc: func(ctx context.Context, entry *Entry) error {
				createdEntry = entry
				return nil
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		entry, err := svc.CreateEntry(context.Background(), fundID, "Alice", 500)
		require.NoError(t, err)
		assert.NotNil(t, entry)
		assert.Equal(t, fundID, entry.FundID)
		assert.Equal(t, "Alice", entry.OwnerName)
		assert.Equal(t, 500, entry.Units)
		assert.Equal(t, createdEntry, entry)
	})

	t.Run("returns error for invalid owner name", func(t *testing.T) {
		repo := &mockRepository{}
		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		entry, err := svc.CreateEntry(context.Background(), fundID, "", 500)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("returns error for invalid units", func(t *testing.T) {
		repo := &mockRepository{}
		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		entry, err := svc.CreateEntry(context.Background(), fundID, "Alice", -1)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidUnits)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		repo := &mockRepository{
			createFunc: func(ctx context.Context, entry *Entry) error {
				return repoErr
			},
		}

		svc, err := NewService(WithRepository(repo))
		require.NoError(t, err)

		entry, err := svc.CreateEntry(context.Background(), fundID, "Alice", 500)
		assert.Nil(t, entry)
		assert.Equal(t, repoErr, err)
	})
}
