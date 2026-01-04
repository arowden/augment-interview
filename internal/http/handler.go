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
)

// logError logs an error with structured context for debugging.
// Always includes the request ID for correlation.
func logError(ctx context.Context, msg string, err error, extraAttrs ...slog.Attr) {
	attrs := make([]slog.Attr, 0, len(extraAttrs)+2)
	attrs = append(attrs, slog.String("error", err.Error()))
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		attrs = append(attrs, slog.String("requestId", reqID))
	}
	attrs = append(attrs, extraAttrs...)
	slog.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

// errorDetails creates a details map with the request ID for traceability.
// Always includes requestId if available from context.
func errorDetails(ctx context.Context, extra map[string]interface{}) *map[string]interface{} {
	details := make(map[string]interface{})

	// Add request ID for traceability.
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		details["requestId"] = reqID
	}

	// Merge any extra details.
	for k, v := range extra {
		details[k] = v
	}

	if len(details) == 0 {
		return nil
	}
	return &details
}

// APIHandler implements the StrictServerInterface for the OpenAPI spec.
type APIHandler struct {
	fundService      *fund.Service
	ownershipService *ownership.Service
	transferService  *transfer.Service
}

// APIHandlerOption configures an APIHandler.
type APIHandlerOption func(*APIHandler)

// WithFundService sets the fund service.
func WithFundService(svc *fund.Service) APIHandlerOption {
	return func(h *APIHandler) {
		h.fundService = svc
	}
}

// WithOwnershipService sets the ownership service.
func WithOwnershipService(svc *ownership.Service) APIHandlerOption {
	return func(h *APIHandler) {
		h.ownershipService = svc
	}
}

// WithTransferService sets the transfer service.
func WithTransferService(svc *transfer.Service) APIHandlerOption {
	return func(h *APIHandler) {
		h.transferService = svc
	}
}

// NewAPIHandler creates a new APIHandler with the provided options.
// For production use, prefer NewAPIHandlerStrict which validates all required services.
func NewAPIHandler(opts ...APIHandlerOption) *APIHandler {
	h := &APIHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// NewAPIHandlerStrict creates a new APIHandler with validation that all required services are configured.
// Returns an error if any required service is missing. Use this in production to fail fast at startup.
func NewAPIHandlerStrict(opts ...APIHandlerOption) (*APIHandler, error) {
	h := &APIHandler{}
	for _, opt := range opts {
		opt(h)
	}

	// Validate all required services are configured.
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

// ListFunds lists all funds.
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

	// Build pagination params.
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

// CreateFund creates a new fund with initial ownership in a single atomic transaction.
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

	// CreateFundWithInitialOwner creates both fund and initial ownership atomically.
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

// GetFund gets a fund by ID.
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

// GetCapTable gets the cap table for a fund.
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

	// Verify fund exists and get total units for percentage calculation (single DB call).
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

	// Build pagination params.
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

// ListTransfers lists transfers for a fund.
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

	// Verify fund exists.
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

	// Build pagination params.
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

// CreateTransfer creates a new transfer.
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

	// Verify fund exists.
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

	// Build transfer request.
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
		// Map domain errors to HTTP responses.
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
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		case errors.Is(err, transfer.ErrOwnerNotFound):
			return CreateTransfer404JSONResponse{
				TransferNotFoundJSONResponse: TransferNotFoundJSONResponse{
					Code:    OWNERNOTFOUND,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
				},
			}, nil
		case errors.Is(err, transfer.ErrInsufficientUnits):
			return CreateTransfer400JSONResponse{
				TransferBadRequestJSONResponse: TransferBadRequestJSONResponse{
					Code:    INSUFFICIENTUNITS,
					Message: err.Error(),
					Details: errorDetails(ctx, nil),
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

// Ensure APIHandler implements StrictServerInterface.
var _ StrictServerInterface = (*APIHandler)(nil)
