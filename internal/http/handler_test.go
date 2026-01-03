package http

import (
	"context"
	"testing"

	"github.com/arowden/augment-fund/internal/fund"
	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	tc, err := postgres.NewTestContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { tc.Cleanup(ctx) })

	fundRepo := fund.NewPostgresRepository(tc.Pool())
	fundService, err := fund.NewService(fundRepo)
	require.NoError(t, err)

	ownershipRepo := ownership.NewPostgresRepository(tc.Pool())
	ownershipService, err := ownership.NewService(ownership.WithRepository(ownershipRepo))
	require.NoError(t, err)

	handler := NewAPIHandler(
		WithFundService(fundService),
		WithOwnershipService(ownershipService),
	)

	t.Run("ListFunds returns empty list initially", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.ListFunds(ctx, ListFundsRequestObject{})
		require.NoError(t, err)

		funds, ok := resp.(ListFunds200JSONResponse)
		require.True(t, ok)
		assert.Empty(t, funds)
	})

	t.Run("CreateFund creates and returns fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:       "Test Fund",
				TotalUnits: 1000,
			},
		})
		require.NoError(t, err)

		created, ok := resp.(CreateFund201JSONResponse)
		require.True(t, ok)
		assert.Equal(t, "Test Fund", created.Name)
		assert.Equal(t, 1000, created.TotalUnits)
		assert.NotEqual(t, uuid.Nil, created.Id)
	})

	t.Run("CreateFund returns 400 for invalid fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:       "",
				TotalUnits: 1000,
			},
		})
		require.NoError(t, err)

		_, ok := resp.(CreateFund400JSONResponse)
		assert.True(t, ok)
	})

	t.Run("CreateFund returns 400 for nil body", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: nil,
		})
		require.NoError(t, err)

		badReq, ok := resp.(CreateFund400JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INVALIDREQUEST, badReq.Code)
	})

	t.Run("GetFund returns 404 for non-existent fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.GetFund(ctx, GetFundRequestObject{
			FundId: uuid.New(),
		})
		require.NoError(t, err)

		_, ok := resp.(GetFund404JSONResponse)
		assert.True(t, ok)
	})

	t.Run("GetFund returns fund by ID", func(t *testing.T) {
		tc.Reset(ctx)

		// Create fund first.
		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:       "Lookup Fund",
				TotalUnits: 500,
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		// Get fund.
		resp, err := handler.GetFund(ctx, GetFundRequestObject{
			FundId: created.Id,
		})
		require.NoError(t, err)

		found, ok := resp.(GetFund200JSONResponse)
		require.True(t, ok)
		assert.Equal(t, created.Id, found.Id)
		assert.Equal(t, "Lookup Fund", found.Name)
	})

	t.Run("GetCapTable returns 404 for non-existent fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: uuid.New(),
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		_, ok := resp.(GetCapTable404JSONResponse)
		assert.True(t, ok)
	})

	t.Run("GetCapTable returns empty entries for fund with no owners", func(t *testing.T) {
		tc.Reset(ctx)

		// Create fund.
		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:       "Empty Fund",
				TotalUnits: 1000,
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		// Get cap table.
		resp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		capTable, ok := resp.(GetCapTable200JSONResponse)
		require.True(t, ok)
		assert.Equal(t, created.Id, capTable.FundId)
		assert.Empty(t, capTable.Entries)
		assert.Equal(t, 0, capTable.Total)
	})

	t.Run("GetCapTable returns entries with percentages", func(t *testing.T) {
		tc.Reset(ctx)

		// Create fund.
		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:       "Cap Table Fund",
				TotalUnits: 1000,
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		// Add ownership entries directly via repository.
		entry1, _ := ownership.NewCapTableEntry(created.Id, "Owner A", 600)
		require.NoError(t, ownershipRepo.Create(ctx, entry1))

		entry2, _ := ownership.NewCapTableEntry(created.Id, "Owner B", 400)
		require.NoError(t, ownershipRepo.Create(ctx, entry2))

		// Get cap table.
		resp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		capTable, ok := resp.(GetCapTable200JSONResponse)
		require.True(t, ok)
		assert.Equal(t, 2, capTable.Total)
		require.Len(t, capTable.Entries, 2)

		// Entries should be ordered by units descending.
		assert.Equal(t, "Owner A", capTable.Entries[0].OwnerName)
		assert.Equal(t, 600, capTable.Entries[0].Units)
		assert.InDelta(t, 60.0, capTable.Entries[0].Percentage, 0.01)

		assert.Equal(t, "Owner B", capTable.Entries[1].OwnerName)
		assert.Equal(t, 400, capTable.Entries[1].Units)
		assert.InDelta(t, 40.0, capTable.Entries[1].Percentage, 0.01)
	})

	t.Run("GetCapTable respects pagination params", func(t *testing.T) {
		tc.Reset(ctx)

		// Create fund.
		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:       "Paginated Fund",
				TotalUnits: 1000,
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		// Add 5 ownership entries.
		for i := 1; i <= 5; i++ {
			entry, _ := ownership.NewCapTableEntry(created.Id, "Owner "+string(rune('A'+i-1)), i*100)
			require.NoError(t, ownershipRepo.Create(ctx, entry))
		}

		// Get first page (limit 2).
		limit := 2
		resp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{
				Limit: &limit,
			},
		})
		require.NoError(t, err)

		capTable, ok := resp.(GetCapTable200JSONResponse)
		require.True(t, ok)
		assert.Equal(t, 5, capTable.Total)
		assert.Len(t, capTable.Entries, 2)
		assert.Equal(t, 2, capTable.Limit)
		assert.Equal(t, 0, capTable.Offset)

		// First two should be highest unit owners (E=500, D=400).
		assert.Equal(t, "Owner E", capTable.Entries[0].OwnerName)
		assert.Equal(t, "Owner D", capTable.Entries[1].OwnerName)

		// Get second page.
		offset := 2
		resp, err = handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{
				Limit:  &limit,
				Offset: &offset,
			},
		})
		require.NoError(t, err)

		capTable, ok = resp.(GetCapTable200JSONResponse)
		require.True(t, ok)
		assert.Len(t, capTable.Entries, 2)
		assert.Equal(t, 2, capTable.Offset)

		// Next two (C=300, B=200).
		assert.Equal(t, "Owner C", capTable.Entries[0].OwnerName)
		assert.Equal(t, "Owner B", capTable.Entries[1].OwnerName)
	})

	t.Run("GetCapTable returns 500 when ownership service not configured", func(t *testing.T) {
		tc.Reset(ctx)

		handlerNoOwnership := NewAPIHandler(
			WithFundService(fundService),
		)

		resp, err := handlerNoOwnership.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: uuid.New(),
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		errResp, ok := resp.(GetCapTable500JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INTERNALERROR, errResp.Code)
	})

	t.Run("ListFunds returns 500 when fund service not configured", func(t *testing.T) {
		tc.Reset(ctx)

		handlerNoFund := NewAPIHandler(
			WithOwnershipService(ownershipService),
		)

		resp, err := handlerNoFund.ListFunds(ctx, ListFundsRequestObject{})
		require.NoError(t, err)

		errResp, ok := resp.(ListFunds500JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INTERNALERROR, errResp.Code)
	})
}
