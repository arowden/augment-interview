package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/arowden/augment-fund/internal/fund"
	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/transfer"
	"github.com/arowden/augment-fund/internal/validation"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func logError(ctx context.Context, msg string, err error, extraAttrs ...slog.Attr) {
	attrs := make([]slog.Attr, 0, len(extraAttrs)+2)
	attrs = append(attrs, slog.String("error", err.Error()))
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		attrs = append(attrs, slog.String("requestId", reqID))
	}
	attrs = append(attrs, extraAttrs...)
	slog.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

func errorDetails(ctx context.Context, extra map[string]interface{}) *map[string]interface{} {
	details := make(map[string]interface{})

	if reqID := middleware.GetReqID(ctx); reqID != "" {
		details["requestId"] = reqID
	}

	for k, v := range extra {
		details[k] = v
	}

	if len(details) == 0 {
		return nil
	}
	return &details
}

type APIHandler struct {
	fundService      *fund.Service
	ownershipService *ownership.Service
	transferService  *transfer.Service
	pool             *pgxpool.Pool
}

type APIHandlerOption func(*APIHandler)

func WithFundService(svc *fund.Service) APIHandlerOption {
	return func(h *APIHandler) {
		h.fundService = svc
	}
}

func WithOwnershipService(svc *ownership.Service) APIHandlerOption {
	return func(h *APIHandler) {
		h.ownershipService = svc
	}
}

func WithTransferService(svc *transfer.Service) APIHandlerOption {
	return func(h *APIHandler) {
		h.transferService = svc
	}
}

func WithPool(p *pgxpool.Pool) APIHandlerOption {
	return func(h *APIHandler) {
		h.pool = p
	}
}

func NewAPIHandler(opts ...APIHandlerOption) *APIHandler {
	h := &APIHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func NewAPIHandlerStrict(opts ...APIHandlerOption) (*APIHandler, error) {
	h := &APIHandler{}
	for _, opt := range opts {
		opt(h)
	}

	var missing []string
	if h.fundService == nil {
		missing = append(missing, "fund service")
	}
	if h.ownershipService == nil {
		missing = append(missing, "ownership service")
	}
	if h.transferService == nil {
		missing = append(missing, "transfer service")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("http: missing required services: %v", missing)
	}

	return h, nil
}

func (h *APIHandler) ListFunds(ctx context.Context, request ListFundsRequestObject) (ListFundsResponseObject, error) {
	if h.fundService == nil {
		return ListFunds500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "fund service not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	params := fund.ListParams{}
	if request.Params.Limit != nil {
		params.Limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		params.Offset = *request.Params.Offset
	}

	result, err := h.fundService.ListFunds(ctx, params)
	if err != nil {
		logError(ctx, "failed to list funds", err)
		return ListFunds500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to list funds",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	funds := make([]Fund, len(result.Items))
	for i, f := range result.Items {
		funds[i] = Fund{
			Id:         f.ID,
			Name:       f.Name,
			TotalUnits: f.TotalUnits,
			CreatedAt:  f.CreatedAt,
		}
	}

	return ListFunds200JSONResponse(FundList{
		Funds:  funds,
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}), nil
}

func (h *APIHandler) CreateFund(ctx context.Context, request CreateFundRequestObject) (CreateFundResponseObject, error) {
	if h.fundService == nil {
		return CreateFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "fund service not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	if request.Body == nil {
		return CreateFund400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{
				Code:    INVALIDREQUEST,
				Message: "request body is required",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	f, err := h.fundService.CreateFundWithInitialOwner(
		ctx,
		request.Body.Name,
		request.Body.TotalUnits,
		request.Body.InitialOwner,
	)
	if err != nil {
		if errors.Is(err, fund.ErrInvalidFund) {
			return CreateFund400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{
					Code:    INVALIDFUND,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		}
		if errors.Is(err, fund.ErrDuplicateFundName) {
			return CreateFund400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{
					Code:    INVALIDFUND,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		}
		logError(ctx, "failed to create fund", err)
		return CreateFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to create fund",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	return CreateFund201JSONResponse(Fund{
		Id:         f.ID,
		Name:       f.Name,
		TotalUnits: f.TotalUnits,
		CreatedAt:  f.CreatedAt,
	}), nil
}

func (h *APIHandler) GetFund(ctx context.Context, request GetFundRequestObject) (GetFundResponseObject, error) {
	if h.fundService == nil {
		return GetFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "fund service not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	f, err := h.fundService.GetFund(ctx, request.FundId)
	if err != nil {
		if errors.Is(err, fund.ErrNotFound) {
			return GetFund404JSONResponse{
				FundNotFoundJSONResponse: FundNotFoundJSONResponse{
					Code:    FUNDNOTFOUND,
					Message: "fund not found",
					Details: errorDetails(ctx, map[string]interface{}{"fundId": request.FundId.String()}),
				},
			}, nil
		}
		logError(ctx, "failed to get fund", err, slog.String("fundId", request.FundId.String()))
		return GetFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to get fund",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	return GetFund200JSONResponse(Fund{
		Id:         f.ID,
		Name:       f.Name,
		TotalUnits: f.TotalUnits,
		CreatedAt:  f.CreatedAt,
	}), nil
}

func (h *APIHandler) GetCapTable(ctx context.Context, request GetCapTableRequestObject) (GetCapTableResponseObject, error) {
	if h.ownershipService == nil {
		return GetCapTable500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "ownership service not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	var fundTotalUnits int
	if h.fundService != nil {
		f, err := h.fundService.GetFund(ctx, request.FundId)
		if err != nil {
			if errors.Is(err, fund.ErrNotFound) {
				return GetCapTable404JSONResponse{
					FundNotFoundJSONResponse: FundNotFoundJSONResponse{
						Code:    FUNDNOTFOUND,
						Message: "fund not found",
						Details: errorDetails(ctx, map[string]interface{}{"fundId": request.FundId.String()}),
					},
				}, nil
			}
			logError(ctx, "failed to verify fund", err, slog.String("fundId", request.FundId.String()))
			return GetCapTable500JSONResponse{
				InternalErrorJSONResponse: InternalErrorJSONResponse{
					Code:    INTERNALERROR,
					Message: "failed to verify fund",
					Details: errorDetails(ctx, nil),
				},
			}, nil
		}
		fundTotalUnits = f.TotalUnits
	}

	params := ownership.ListParams{}
	if request.Params.Limit != nil {
		params.Limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		params.Offset = *request.Params.Offset
	}

	view, err := h.ownershipService.GetCapTable(ctx, request.FundId, params)
	if err != nil {
		logError(ctx, "failed to get cap table", err, slog.String("fundId", request.FundId.String()))
		return GetCapTable500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to get cap table",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	entries := make([]CapTableEntry, len(view.Entries))
	for i, e := range view.Entries {
		var percentage float64
		if fundTotalUnits > 0 {
			percentage = float64(e.Units) / float64(fundTotalUnits) * validation.PercentageMultiplier
		}
		entries[i] = CapTableEntry{
			OwnerName:  e.OwnerName,
			Units:      e.Units,
			AcquiredAt: e.AcquiredAt,
			Percentage: percentage,
		}
	}

	return GetCapTable200JSONResponse(CapTable{
		FundId:  request.FundId,
		Entries: entries,
		Total:   view.TotalCount,
		Limit:   view.Limit,
		Offset:  view.Offset,
	}), nil
}

func (h *APIHandler) ListTransfers(ctx context.Context, request ListTransfersRequestObject) (ListTransfersResponseObject, error) {
	if h.transferService == nil {
		return ListTransfers500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "transfer service not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	if h.fundService != nil {
		_, err := h.fundService.GetFund(ctx, request.FundId)
		if err != nil {
			if errors.Is(err, fund.ErrNotFound) {
				return ListTransfers404JSONResponse{
					FundNotFoundJSONResponse: FundNotFoundJSONResponse{
						Code:    FUNDNOTFOUND,
						Message: "fund not found",
						Details: errorDetails(ctx, map[string]interface{}{"fundId": request.FundId.String()}),
					},
				}, nil
			}
			logError(ctx, "failed to verify fund", err, slog.String("fundId", request.FundId.String()))
			return ListTransfers500JSONResponse{
				InternalErrorJSONResponse: InternalErrorJSONResponse{
					Code:    INTERNALERROR,
					Message: "failed to verify fund",
					Details: errorDetails(ctx, nil),
				},
			}, nil
		}
	}

	params := transfer.ListParams{}
	if request.Params.Limit != nil {
		params.Limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		params.Offset = *request.Params.Offset
	}

	list, err := h.transferService.ListTransfers(ctx, request.FundId, params)
	if err != nil {
		logError(ctx, "failed to list transfers", err, slog.String("fundId", request.FundId.String()))
		return ListTransfers500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to list transfers",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	transfers := make([]Transfer, len(list.Transfers))
	for i, t := range list.Transfers {
		transfers[i] = Transfer{
			Id:            t.ID,
			FundId:        t.FundID,
			FromOwner:     t.FromOwner,
			ToOwner:       t.ToOwner,
			Units:         t.Units,
			TransferredAt: t.TransferredAt,
		}
	}

	return ListTransfers200JSONResponse(TransferList{
		FundId:    request.FundId,
		Transfers: transfers,
		Total:     list.TotalCount,
		Limit:     list.Limit,
		Offset:    list.Offset,
	}), nil
}

func (h *APIHandler) CreateTransfer(ctx context.Context, request CreateTransferRequestObject) (CreateTransferResponseObject, error) {
	if h.transferService == nil {
		return CreateTransfer500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "transfer service not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	if request.Body == nil {
		return CreateTransfer400JSONResponse{
			TransferBadRequestJSONResponse: TransferBadRequestJSONResponse{
				Code:    INVALIDREQUEST,
				Message: "request body is required",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}

	if h.fundService != nil {
		_, err := h.fundService.GetFund(ctx, request.FundId)
		if err != nil {
			if errors.Is(err, fund.ErrNotFound) {
				return CreateTransfer404JSONResponse{
					TransferNotFoundJSONResponse: TransferNotFoundJSONResponse{
						Code:    FUNDNOTFOUND,
						Message: "fund not found",
						Details: errorDetails(ctx, map[string]interface{}{"fundId": request.FundId.String()}),
					},
				}, nil
			}
			logError(ctx, "failed to verify fund for transfer", err, slog.String("fundId", request.FundId.String()))
			return CreateTransfer500JSONResponse{
				InternalErrorJSONResponse: InternalErrorJSONResponse{
					Code:    INTERNALERROR,
					Message: "failed to verify fund",
					Details: errorDetails(ctx, nil),
				},
			}, nil
		}
	}

	req := transfer.Request{
		FundID:    request.FundId,
		FromOwner: request.Body.FromOwner,
		ToOwner:   request.Body.ToOwner,
		Units:     request.Body.Units,
	}
	if request.Body.IdempotencyKey != nil {
		key := uuid.UUID(*request.Body.IdempotencyKey)
		req.IdempotencyKey = &key
	}

	t, err := h.transferService.ExecuteTransfer(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, transfer.ErrInvalidOwner):
			return CreateTransfer400JSONResponse{
				TransferBadRequestJSONResponse: TransferBadRequestJSONResponse{
					Code:    INVALIDREQUEST,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		case errors.Is(err, transfer.ErrInvalidUnits):
			return CreateTransfer400JSONResponse{
				TransferBadRequestJSONResponse: TransferBadRequestJSONResponse{
					Code:    INVALIDREQUEST,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		case errors.Is(err, transfer.ErrSelfTransfer):
			return CreateTransfer400JSONResponse{
				TransferBadRequestJSONResponse: TransferBadRequestJSONResponse{
					Code:    INVALIDREQUEST,
					Message: fmt.Sprintf("Cannot transfer units to yourself. Both fromOwner and toOwner are '%s'.", request.Body.FromOwner),
					Details: errorDetails(ctx, map[string]interface{}{
						"fromOwner": request.Body.FromOwner,
						"toOwner":   request.Body.ToOwner,
					}),
				},
			}, nil
		case errors.Is(err, transfer.ErrOwnerNotFound):
			return CreateTransfer404JSONResponse{
				TransferNotFoundJSONResponse: TransferNotFoundJSONResponse{
					Code:    OWNERNOTFOUND,
					Message: fmt.Sprintf("Owner '%s' does not own any units in this fund. Check the cap table to see current owners.", request.Body.FromOwner),
					Details: errorDetails(ctx, map[string]interface{}{
						"ownerName": request.Body.FromOwner,
						"fundId":    request.FundId.String(),
						"hint":      "The fromOwner must be an existing owner in the fund's cap table",
					}),
				},
			}, nil
		case errors.Is(err, transfer.ErrInsufficientUnits):
			return CreateTransfer400JSONResponse{
				TransferBadRequestJSONResponse: TransferBadRequestJSONResponse{
					Code:    INSUFFICIENTUNITS,
					Message: fmt.Sprintf("Owner '%s' does not have enough units for this transfer. They need at least %d units.", request.Body.FromOwner, request.Body.Units),
					Details: errorDetails(ctx, map[string]interface{}{
						"ownerName":      request.Body.FromOwner,
						"requestedUnits": request.Body.Units,
						"hint":           "Check the cap table to see how many units this owner currently holds",
					}),
				},
			}, nil
		case errors.Is(err, transfer.ErrDuplicateIdempotencyKey):
			return CreateTransfer409JSONResponse{
				DuplicateTransferJSONResponse: DuplicateTransferJSONResponse{
					Code:    DUPLICATETRANSFER,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		default:
			logError(ctx, "failed to execute transfer", err,
				slog.String("fundId", request.FundId.String()),
				slog.String("fromOwner", request.Body.FromOwner),
				slog.String("toOwner", request.Body.ToOwner),
				slog.Int("units", request.Body.Units),
			)
			return CreateTransfer500JSONResponse{
				InternalErrorJSONResponse: InternalErrorJSONResponse{
					Code:    INTERNALERROR,
					Message: "failed to execute transfer",
					Details: errorDetails(ctx, nil),
				},
			}, nil
		}
	}

	return CreateTransfer201JSONResponse(Transfer{
		Id:            t.ID,
		FundId:        t.FundID,
		FromOwner:     t.FromOwner,
		ToOwner:       t.ToOwner,
		Units:         t.Units,
		TransferredAt: t.TransferredAt,
	}), nil
}

func (h *APIHandler) ResetDatabase(ctx context.Context, _ ResetDatabaseRequestObject) (ResetDatabaseResponseObject, error) {
	if h.pool == nil {
		return ResetDatabase500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "database pool not configured",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}


	var deletedTransfers, deletedOwnership, deletedFunds int

	result, err := h.pool.Exec(ctx, "DELETE FROM transfers")
	if err != nil {
		logError(ctx, "failed to delete transfers", err)
		return ResetDatabase500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to reset database",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}
	deletedTransfers = int(result.RowsAffected())

	result, err = h.pool.Exec(ctx, "DELETE FROM cap_table_entries")
	if err != nil {
		logError(ctx, "failed to delete ownership entries", err)
		return ResetDatabase500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to reset database",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}
	deletedOwnership = int(result.RowsAffected())

	result, err = h.pool.Exec(ctx, "DELETE FROM funds")
	if err != nil {
		logError(ctx, "failed to delete funds", err)
		return ResetDatabase500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to reset database",
				Details: errorDetails(ctx, nil),
			},
		}, nil
	}
	deletedFunds = int(result.RowsAffected())

	slog.InfoContext(ctx, "database reset completed",
		slog.Int("deletedFunds", deletedFunds),
		slog.Int("deletedTransfers", deletedTransfers),
		slog.Int("deletedOwnership", deletedOwnership),
	)

	return ResetDatabase200JSONResponse{
		Message:          ptr("Database reset successfully"),
		DeletedFunds:     ptr(deletedFunds),
		DeletedTransfers: ptr(deletedTransfers),
		DeletedOwnership: ptr(deletedOwnership),
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}

var _ StrictServerInterface = (*APIHandler)(nil)
