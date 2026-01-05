package fund

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

type mockRepository struct {
	createFunc   func(ctx context.Context, fund *Fund) error
	createTxFunc func(ctx context.Context, tx pgx.Tx, fund *Fund) error
	findByIDFunc func(ctx context.Context, id uuid.UUID) (*Fund, error)
	listFunc     func(ctx context.Context, params ListParams) (*ListResult, error)
}

func (m *mockRepository) Create(ctx context.Context, fund *Fund) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, fund)
	}
	return nil
}

func (m *mockRepository) CreateTx(ctx context.Context, tx pgx.Tx, fund *Fund) error {
	if m.createTxFunc != nil {
		return m.createTxFunc(ctx, tx, fund)
	}
	return nil
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*Fund, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, NotFoundError(id)
}

func (m *mockRepository) List(ctx context.Context, params ListParams) (*ListResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, params)
	}
	return &ListResult{Items: []*Fund{}}, nil
}

type mockOwnershipRepository struct {
	createTxFunc func(ctx context.Context, tx pgx.Tx, entry *ownership.Entry) error
}

func (m *mockOwnershipRepository) Create(_ context.Context, _ *ownership.Entry) error {
	return nil
}

func (m *mockOwnershipRepository) CreateTx(ctx context.Context, tx pgx.Tx, entry *ownership.Entry) error {
	if m.createTxFunc != nil {
		return m.createTxFunc(ctx, tx, entry)
	}
	return nil
}

func (m *mockOwnershipRepository) FindByFundID(_ context.Context, _ uuid.UUID, _ ownership.ListParams) (*ownership.CapTableView, error) {
	return nil, nil
}

func (m *mockOwnershipRepository) FindByFundAndOwner(_ context.Context, _ uuid.UUID, _ string) (*ownership.Entry, error) {
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
	t.Run("returns error when repo is nil", func(t *testing.T) {
		svc, err := NewService(nil)
		assert.Nil(t, svc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository is required")
	})

	t.Run("creates service with valid repo", func(t *testing.T) {
		repo := &mockRepository{}
		svc, err := NewService(repo)
		require.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("applies options", func(t *testing.T) {
		repo := &mockRepository{}
		ownershipRepo := &mockOwnershipRepository{}

		svc, err := NewService(repo, WithOwnershipRepository(ownershipRepo))
		require.NoError(t, err)
		assert.NotNil(t, svc)
		assert.Equal(t, ownershipRepo, svc.ownershipRepo)
	})
}

func TestService_GetFund(t *testing.T) {
	t.Run("returns fund when found", func(t *testing.T) {
		expectedFund := &Fund{
			ID:         uuid.New(),
			Name:       "Test Fund",
			TotalUnits: 1000,
		}
		repo := &mockRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*Fund, error) {
				return expectedFund, nil
			},
		}

		svc, err := NewService(repo)
		require.NoError(t, err)

		fund, err := svc.GetFund(context.Background(), expectedFund.ID)
		require.NoError(t, err)
		assert.Equal(t, expectedFund, fund)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*Fund, error) {
				return nil, NotFoundError(id)
			},
		}

		svc, err := NewService(repo)
		require.NoError(t, err)

		fund, err := svc.GetFund(context.Background(), uuid.New())
		assert.Nil(t, fund)
		assert.True(t, errors.Is(err, ErrNotFound))
	})
}

func TestService_ListFunds(t *testing.T) {
	t.Run("returns funds list", func(t *testing.T) {
		expectedResult := &ListResult{
			Items: []*Fund{
				{ID: uuid.New(), Name: "Fund 1", TotalUnits: 100},
				{ID: uuid.New(), Name: "Fund 2", TotalUnits: 200},
			},
			Total:  2,
			Limit:  100,
			Offset: 0,
		}
		repo := &mockRepository{
			listFunc: func(ctx context.Context, params ListParams) (*ListResult, error) {
				return expectedResult, nil
			},
		}

		svc, err := NewService(repo)
		require.NoError(t, err)

		result, err := svc.ListFunds(context.Background(), ListParams{})
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("passes params to repository", func(t *testing.T) {
		var receivedParams ListParams
		repo := &mockRepository{
			listFunc: func(ctx context.Context, params ListParams) (*ListResult, error) {
				receivedParams = params
				return &ListResult{Items: []*Fund{}}, nil
			},
		}

		svc, err := NewService(repo)
		require.NoError(t, err)

		params := ListParams{Limit: 50, Offset: 10}
		_, err = svc.ListFunds(context.Background(), params)
		require.NoError(t, err)
		assert.Equal(t, params, receivedParams)
	})
}

func TestService_CreateFundWithInitialOwner(t *testing.T) {
	t.Run("returns error when pool is nil", func(t *testing.T) {
		repo := &mockRepository{}
		ownershipRepo := &mockOwnershipRepository{}

		svc, err := NewService(repo, WithOwnershipRepository(ownershipRepo))
		require.NoError(t, err)

		fund, err := svc.CreateFundWithInitialOwner(context.Background(), "Test Fund", 1000, "Owner")
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrPoolRequired)
	})

	t.Run("returns error when ownership repo is nil", func(t *testing.T) {
		repo := &mockRepository{}

		svc, err := NewService(repo)
		require.NoError(t, err)

		fund, err := svc.CreateFundWithInitialOwner(context.Background(), "Test Fund", 1000, "Owner")
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrPoolRequired)
	})

	t.Run("returns validation error for invalid fund name", func(t *testing.T) {
		repo := &mockRepository{}

		svc, err := NewService(repo)
		require.NoError(t, err)

		fund, err := svc.CreateFundWithInitialOwner(context.Background(), "", 1000, "Owner")
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("returns validation error for invalid units", func(t *testing.T) {
		repo := &mockRepository{}

		svc, err := NewService(repo)
		require.NoError(t, err)

		fund, err := svc.CreateFundWithInitialOwner(context.Background(), "Test Fund", 0, "Owner")
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("returns validation error for invalid owner name", func(t *testing.T) {
		repo := &mockRepository{}

		svc, err := NewService(repo)
		require.NoError(t, err)

		fund, err := svc.CreateFundWithInitialOwner(context.Background(), "Test Fund", 1000, "")
		assert.Nil(t, fund)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid initial owner")
	})
}
