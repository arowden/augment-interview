
package http

import (
	"context"
	"testing"

	"github.com/arowden/augment-fund/internal/fund"
	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/postgres"
	"github.com/arowden/augment-fund/internal/transfer"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

	fundStore := fund.NewStore(tc.Pool())
	ownershipStore := ownership.NewStore(tc.Pool())

	fundService, err := fund.NewService(
		fundStore,
		fund.WithPool(tc.Pool()),
		fund.WithOwnershipRepository(ownershipStore),
	)
	require.NoError(t, err)

	ownershipService, err := ownership.NewService(ownership.WithRepository(ownershipStore))
	require.NoError(t, err)

	handler := NewAPIHandler(
		WithFundService(fundService),
		WithOwnershipService(ownershipService),
	)

	t.Run("ListFunds returns empty list initially", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.ListFunds(ctx, ListFundsRequestObject{})
		require.NoError(t, err)

		fundList, ok := resp.(ListFunds200JSONResponse)
		require.True(t, ok)
		assert.Empty(t, fundList.Funds)
		assert.Equal(t, 0, fundList.Total)
		assert.Equal(t, 100, fundList.Limit)
		assert.Equal(t, 0, fundList.Offset)
	})

	t.Run("CreateFund creates and returns fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Test Fund",
				TotalUnits:   1000,
				InitialOwner: "Founder LLC",
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
				Name:         "",
				TotalUnits:   1000,
				InitialOwner: "Founder LLC",
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

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Lookup Fund",
				TotalUnits:   500,
				InitialOwner: "Founder LLC",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

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

	t.Run("GetCapTable returns entries for fund with initial owner", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Initial Owner Fund",
				TotalUnits:   1000,
				InitialOwner: "Founder LLC",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		capTable, ok := resp.(GetCapTable200JSONResponse)
		require.True(t, ok)
		assert.Equal(t, created.Id, capTable.FundId)
		require.Len(t, capTable.Entries, 1)
		assert.Equal(t, 1, capTable.Total)
		assert.Equal(t, "Founder LLC", capTable.Entries[0].OwnerName)
		assert.Equal(t, 1000, capTable.Entries[0].Units)
		assert.InDelta(t, 100.0, capTable.Entries[0].Percentage, 0.01)
	})

	t.Run("GetCapTable respects pagination params", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Paginated Fund",
				TotalUnits:   1500,
				InitialOwner: "Initial Owner",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		for i := 1; i <= 4; i++ {
			entry, _ := ownership.NewCapTableEntry(created.Id, "Owner "+string(rune('A'+i-1)), i*100)
			require.NoError(t, ownershipStore.Create(ctx, entry))
		}

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

		assert.Equal(t, "Initial Owner", capTable.Entries[0].OwnerName)
		assert.Equal(t, "Owner D", capTable.Entries[1].OwnerName)

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

func TestAPIHandler_Transfers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	tc, err := postgres.NewTestContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { tc.Cleanup(ctx) })

	fundStore := fund.NewStore(tc.Pool())
	ownershipStore := ownership.NewStore(tc.Pool())
	transferStore := transfer.NewStore(tc.Pool())

	fundService, err := fund.NewService(
		fundStore,
		fund.WithPool(tc.Pool()),
		fund.WithOwnershipRepository(ownershipStore),
	)
	require.NoError(t, err)

	ownershipService, err := ownership.NewService(ownership.WithRepository(ownershipStore))
	require.NoError(t, err)

	transferService, err := transfer.NewService(
		transfer.WithRepository(transferStore),
		transfer.WithOwnershipRepository(ownershipStore),
		transfer.WithPool(tc.Pool()),
	)
	require.NoError(t, err)

	handler := NewAPIHandler(
		WithFundService(fundService),
		WithOwnershipService(ownershipService),
		WithTransferService(transferService),
	)

	t.Run("ListTransfers returns empty list for fund with no transfers", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Transfer Test Fund",
				TotalUnits:   1000,
				InitialOwner: "Founder",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.ListTransfers(ctx, ListTransfersRequestObject{
			FundId: created.Id,
			Params: ListTransfersParams{},
		})
		require.NoError(t, err)

		transferList, ok := resp.(ListTransfers200JSONResponse)
		require.True(t, ok)
		assert.Empty(t, transferList.Transfers)
		assert.Equal(t, 0, transferList.Total)
		assert.Equal(t, created.Id, transferList.FundId)
	})

	t.Run("ListTransfers returns 404 for non-existent fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.ListTransfers(ctx, ListTransfersRequestObject{
			FundId: uuid.New(),
			Params: ListTransfersParams{},
		})
		require.NoError(t, err)

		_, ok := resp.(ListTransfers404JSONResponse)
		assert.True(t, ok)
	})

	t.Run("CreateTransfer successfully transfers units", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Transfer Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     200,
			},
		})
		require.NoError(t, err)

		transferResp, ok := resp.(CreateTransfer201JSONResponse)
		require.True(t, ok)
		assert.Equal(t, created.Id, transferResp.FundId)
		assert.Equal(t, "Alice", transferResp.FromOwner)
		assert.Equal(t, "Bob", transferResp.ToOwner)
		assert.Equal(t, 200, transferResp.Units)

		capResp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		capTable := capResp.(GetCapTable200JSONResponse)
		assert.Len(t, capTable.Entries, 2)

		var aliceUnits, bobUnits int
		for _, e := range capTable.Entries {
			if e.OwnerName == "Alice" {
				aliceUnits = e.Units
			} else if e.OwnerName == "Bob" {
				bobUnits = e.Units
			}
		}
		assert.Equal(t, 800, aliceUnits)
		assert.Equal(t, 200, bobUnits)
	})

	t.Run("CreateTransfer returns 404 for non-existent fund", func(t *testing.T) {
		tc.Reset(ctx)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: uuid.New(),
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     100,
			},
		})
		require.NoError(t, err)

		_, ok := resp.(CreateTransfer404JSONResponse)
		assert.True(t, ok)
	})

	t.Run("CreateTransfer returns 404 for non-existent owner", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Owner Test Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "NonExistent",
				ToOwner:   "Bob",
				Units:     100,
			},
		})
		require.NoError(t, err)

		errResp, ok := resp.(CreateTransfer404JSONResponse)
		require.True(t, ok)
		assert.Equal(t, OWNERNOTFOUND, errResp.Code)
	})

	t.Run("CreateTransfer returns 400 for insufficient units", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Insufficient Units Fund",
				TotalUnits:   100,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     500,
			},
		})
		require.NoError(t, err)

		errResp, ok := resp.(CreateTransfer400JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INSUFFICIENTUNITS, errResp.Code)
	})

	t.Run("CreateTransfer returns 400 for self transfer", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Self Transfer Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Alice",
				Units:     100,
			},
		})
		require.NoError(t, err)

		errResp, ok := resp.(CreateTransfer400JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INVALIDREQUEST, errResp.Code)
	})

	t.Run("CreateTransfer returns 400 for invalid units", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Invalid Units Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     0,
			},
		})
		require.NoError(t, err)

		errResp, ok := resp.(CreateTransfer400JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INVALIDREQUEST, errResp.Code)
	})

	t.Run("CreateTransfer returns 400 for invalid owner name", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Invalid Owner Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "",
				ToOwner:   "Bob",
				Units:     100,
			},
		})
		require.NoError(t, err)

		errResp, ok := resp.(CreateTransfer400JSONResponse)
		require.True(t, ok)
		assert.Equal(t, INVALIDREQUEST, errResp.Code)
	})

	t.Run("CreateTransfer with idempotency key returns same transfer", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Idempotency Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		idempotencyKey := uuid.New()

		resp1, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner:      "Alice",
				ToOwner:        "Bob",
				Units:          100,
				IdempotencyKey: (*openapi_types.UUID)(&idempotencyKey),
			},
		})
		require.NoError(t, err)
		transfer1 := resp1.(CreateTransfer201JSONResponse)

		resp2, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner:      "Alice",
				ToOwner:        "Bob",
				Units:          100,
				IdempotencyKey: (*openapi_types.UUID)(&idempotencyKey),
			},
		})
		require.NoError(t, err)
		transfer2 := resp2.(CreateTransfer201JSONResponse)

		assert.Equal(t, transfer1.Id, transfer2.Id)

		capResp, err := handler.GetCapTable(ctx, GetCapTableRequestObject{
			FundId: created.Id,
			Params: GetCapTableParams{},
		})
		require.NoError(t, err)

		capTable := capResp.(GetCapTable200JSONResponse)
		var bobUnits int
		for _, e := range capTable.Entries {
			if e.OwnerName == "Bob" {
				bobUnits = e.Units
			}
		}
		assert.Equal(t, 100, bobUnits)
	})

	t.Run("CreateTransfer returns 409 for duplicate idempotency key with different data", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Duplicate Key Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		idempotencyKey := uuid.New()

		_, err = handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner:      "Alice",
				ToOwner:        "Bob",
				Units:          100,
				IdempotencyKey: (*openapi_types.UUID)(&idempotencyKey),
			},
		})
		require.NoError(t, err)

		resp, err := handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner:      "Alice",
				ToOwner:        "Bob",
				Units:          200,
				IdempotencyKey: (*openapi_types.UUID)(&idempotencyKey),
			},
		})
		require.NoError(t, err)

		errResp, ok := resp.(CreateTransfer409JSONResponse)
		require.True(t, ok)
		assert.Equal(t, DUPLICATETRANSFER, errResp.Code)
	})

	t.Run("ListTransfers returns transfers after creation", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "List Transfers Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		_, err = handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Bob",
				Units:     100,
			},
		})
		require.NoError(t, err)

		_, err = handler.CreateTransfer(ctx, CreateTransferRequestObject{
			FundId: created.Id,
			Body: &CreateTransferJSONRequestBody{
				FromOwner: "Alice",
				ToOwner:   "Charlie",
				Units:     200,
			},
		})
		require.NoError(t, err)

		resp, err := handler.ListTransfers(ctx, ListTransfersRequestObject{
			FundId: created.Id,
			Params: ListTransfersParams{},
		})
		require.NoError(t, err)

		transferList := resp.(ListTransfers200JSONResponse)
		assert.Equal(t, 2, transferList.Total)
		assert.Len(t, transferList.Transfers, 2)
	})

	t.Run("ListTransfers respects pagination", func(t *testing.T) {
		tc.Reset(ctx)

		createResp, err := handler.CreateFund(ctx, CreateFundRequestObject{
			Body: &CreateFundJSONRequestBody{
				Name:         "Pagination Fund",
				TotalUnits:   1000,
				InitialOwner: "Alice",
			},
		})
		require.NoError(t, err)
		created := createResp.(CreateFund201JSONResponse)

		for i := 0; i < 3; i++ {
			_, err = handler.CreateTransfer(ctx, CreateTransferRequestObject{
				FundId: created.Id,
				Body: &CreateTransferJSONRequestBody{
					FromOwner: "Alice",
					ToOwner:   "Bob",
					Units:     10,
				},
			})
			require.NoError(t, err)
		}

		limit := 2
		resp, err := handler.ListTransfers(ctx, ListTransfersRequestObject{
			FundId: created.Id,
			Params: ListTransfersParams{Limit: &limit},
		})
		require.NoError(t, err)

		transferList := resp.(ListTransfers200JSONResponse)
		assert.Equal(t, 3, transferList.Total)
		assert.Len(t, transferList.Transfers, 2)
		assert.Equal(t, 2, transferList.Limit)
		assert.Equal(t, 0, transferList.Offset)

		offset := 2
		resp, err = handler.ListTransfers(ctx, ListTransfersRequestObject{
			FundId: created.Id,
			Params: ListTransfersParams{Limit: &limit, Offset: &offset},
		})
		require.NoError(t, err)

		transferList = resp.(ListTransfers200JSONResponse)
		assert.Len(t, transferList.Transfers, 1)
		assert.Equal(t, 2, transferList.Offset)
	})
}
