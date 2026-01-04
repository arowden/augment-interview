package transfer

import (
	"context"
	"errors"
	"testing"

	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepository implements Repository for testing.
type mockRepository struct {
	findByFundIDFunc        func(ctx context.Context, fundID uuid.UUID, params ListParams) (*TransferList, error)
	findByIdempotencyKeyFunc func(ctx context.Context, tx pgx.Tx, key uuid.UUID) (*Transfer, error)
	createFunc              func(ctx context.Context, t *Transfer) error
	createTxFunc            func(ctx context.Context, tx pgx.Tx, t *Transfer) error
}

func (m *mockRepository) FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*TransferList, error) {
	if m.findByFundIDFunc != nil {
		return m.findByFundIDFunc(ctx, fundID, params)
	}
	return &TransferList{Transfers: []*Transfer{}}, nil
}

func (m *mockRepository) FindByIdempotencyKey(ctx context.Context, tx pgx.Tx, key uuid.UUID) (*Transfer, error) {
	if m.findByIdempotencyKeyFunc != nil {
		return m.findByIdempotencyKeyFunc(ctx, tx, key)
	}
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, t *Transfer) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, t)
	}
	return nil
}

func (m *mockRepository) CreateTx(ctx context.Context, tx pgx.Tx, t *Transfer) error {
	if m.createTxFunc != nil {
		return m.createTxFunc(ctx, tx, t)
	}
	return nil
}

// mockOwnershipRepository implements ownership.Repository for testing.
type mockOwnershipRepository struct{}

func (m *mockOwnershipRepository) Create(ctx context.Context, entry *ownership.Entry) error {
	return nil
}

func (m *mockOwnershipRepository) CreateTx(ctx context.Context, tx pgx.Tx, entry *ownership.Entry) error {
	return nil
}

func (m *mockOwnershipRepository) FindByFundID(ctx context.Context, fundID uuid.UUID, params ownership.ListParams) (*ownership.CapTableView, error) {
	return nil, nil
}

func (m *mockOwnershipRepository) FindByFundAndOwner(ctx context.Context, fundID uuid.UUID, ownerName string) (*ownership.Entry, error) {
	return nil, nil
}

func (m *mockOwnershipRepository) FindByFundAndOwnerForUpdateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string) (*ownership.Entry, error) {
	return nil, nil
}

func (m *mockOwnershipRepository) DecrementUnitsTx(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, units int) error {
	return nil
}

func (m *mockOwnershipRepository) IncrementOrCreateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string, units int) error {
	return nil
}

func (m *mockOwnershipRepository) Upsert(ctx context.Context, entry *ownership.Entry) error {
	return nil
}

func (m *mockOwnershipRepository) UpsertTx(ctx context.Context, tx pgx.Tx, entry *ownership.Entry) error {
	return nil
}

func TestNewService(t *testing.T) {
	t.Run("returns error when repository is nil", func(t *testing.T) {
		svc, err := NewService()
		assert.Nil(t, svc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository is required")
	})

	t.Run("returns error when ownership repository is nil", func(t *testing.T) {
		repo := &mockRepository{}
		svc, err := NewService(WithRepository(repo))
		assert.Nil(t, svc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ownership repository is required")
	})

	t.Run("returns error when pool is nil", func(t *testing.T) {
		repo := &mockRepository{}
		ownershipRepo := &mockOwnershipRepository{}
		svc, err := NewService(
			WithRepository(repo),
			WithOwnershipRepository(ownershipRepo),
		)
		assert.Nil(t, svc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pool is required")
	})
}

func TestService_ExecuteTransfer_Validation(t *testing.T) {
	// Note: We can't fully test ExecuteTransfer without a real pool,
	// but we can test that validation fails before pool usage.

	t.Run("returns error for empty from_owner", func(t *testing.T) {
		// We need a minimal service with just the validator
		svc := &Service{validator: NewValidator()}

		req := Request{
			FundID:    uuid.New(),
			FromOwner: "",
			ToOwner:   "Bob",
			Units:     100,
		}

		_, err := svc.ExecuteTransfer(context.Background(), req)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("returns error for empty to_owner", func(t *testing.T) {
		svc := &Service{validator: NewValidator()}

		req := Request{
			FundID:    uuid.New(),
			FromOwner: "Alice",
			ToOwner:   "",
			Units:     100,
		}

		_, err := svc.ExecuteTransfer(context.Background(), req)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("returns error for self transfer", func(t *testing.T) {
		svc := &Service{validator: NewValidator()}

		req := Request{
			FundID:    uuid.New(),
			FromOwner: "Alice",
			ToOwner:   "Alice",
			Units:     100,
		}

		_, err := svc.ExecuteTransfer(context.Background(), req)
		assert.ErrorIs(t, err, ErrSelfTransfer)
	})

	t.Run("returns error for zero units", func(t *testing.T) {
		svc := &Service{validator: NewValidator()}

		req := Request{
			FundID:    uuid.New(),
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     0,
		}

		_, err := svc.ExecuteTransfer(context.Background(), req)
		assert.ErrorIs(t, err, ErrInvalidUnits)
	})

	t.Run("returns error for negative units", func(t *testing.T) {
		svc := &Service{validator: NewValidator()}

		req := Request{
			FundID:    uuid.New(),
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     -100,
		}

		_, err := svc.ExecuteTransfer(context.Background(), req)
		assert.ErrorIs(t, err, ErrInvalidUnits)
	})
}

func TestService_ListTransfers(t *testing.T) {
	// This test requires a service with at least a repo configured.
	// We'll test that ListTransfers correctly delegates to the repository.

	fundID := uuid.New()
	expectedList := &TransferList{
		Transfers: []*Transfer{
			{ID: uuid.New(), FundID: fundID, FromOwner: "Alice", ToOwner: "Bob", Units: 100},
		},
		TotalCount: 1,
		Limit:      100,
		Offset:     0,
	}

	t.Run("returns transfer list from repository", func(t *testing.T) {
		repo := &mockRepository{
			findByFundIDFunc: func(ctx context.Context, fID uuid.UUID, params ListParams) (*TransferList, error) {
				return expectedList, nil
			},
		}

		// Create minimal service with just the repo (bypassing constructor validation for test)
		svc := &Service{repo: repo, validator: NewValidator()}

		result, err := svc.ListTransfers(context.Background(), fundID, ListParams{})
		require.NoError(t, err)
		assert.Equal(t, expectedList, result)
	})

	t.Run("passes params to repository", func(t *testing.T) {
		var receivedParams ListParams
		repo := &mockRepository{
			findByFundIDFunc: func(ctx context.Context, fID uuid.UUID, params ListParams) (*TransferList, error) {
				receivedParams = params
				return &TransferList{Transfers: []*Transfer{}}, nil
			},
		}

		svc := &Service{repo: repo, validator: NewValidator()}
		params := ListParams{Limit: 50, Offset: 10}

		_, err := svc.ListTransfers(context.Background(), fundID, params)
		require.NoError(t, err)
		assert.Equal(t, params, receivedParams)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		repo := &mockRepository{
			findByFundIDFunc: func(ctx context.Context, fID uuid.UUID, params ListParams) (*TransferList, error) {
				return nil, repoErr
			},
		}

		svc := &Service{repo: repo, validator: NewValidator()}

		result, err := svc.ListTransfers(context.Background(), fundID, ListParams{})
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})
}
