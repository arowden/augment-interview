package http

import (
	"context"
	"errors"

	"github.com/arowden/augment-fund/internal/fund"
	"github.com/arowden/augment-fund/internal/ownership"
)

// APIHandler implements the StrictServerInterface for the OpenAPI spec.
type APIHandler struct {
	fundService      *fund.Service
	ownershipService *ownership.Service
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

// NewAPIHandler creates a new APIHandler with the provided options.
func NewAPIHandler(opts ...APIHandlerOption) *APIHandler {
	h := &APIHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// ListFunds lists all funds.
func (h *APIHandler) ListFunds(ctx context.Context, request ListFundsRequestObject) (ListFundsResponseObject, error) {
	if h.fundService == nil {
		return ListFunds500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "fund service not configured",
			},
		}, nil
	}

	result, err := h.fundService.ListFunds(ctx, fund.ListParams{})
	if err != nil {
		return ListFunds500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to list funds",
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

	return ListFunds200JSONResponse(funds), nil
}

// CreateFund creates a new fund.
func (h *APIHandler) CreateFund(ctx context.Context, request CreateFundRequestObject) (CreateFundResponseObject, error) {
	if h.fundService == nil {
		return CreateFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "fund service not configured",
			},
		}, nil
	}

	if request.Body == nil {
		return CreateFund400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{
				Code:    INVALIDREQUEST,
				Message: "request body is required",
			},
		}, nil
	}

	f, err := h.fundService.CreateFund(ctx, request.Body.Name, request.Body.TotalUnits)
	if err != nil {
		if errors.Is(err, fund.ErrInvalidFund) {
			return CreateFund400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{
					Code:    INVALIDFUND,
					Message: err.Error(),
				},
			}, nil
		}
		if errors.Is(err, fund.ErrDuplicateFundName) {
			return CreateFund400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{
					Code:    INVALIDFUND,
					Message: err.Error(),
				},
			}, nil
		}
		return CreateFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to create fund",
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
				},
			}, nil
		}
		return GetFund500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to get fund",
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
					},
				}, nil
			}
			return GetCapTable500JSONResponse{
				InternalErrorJSONResponse: InternalErrorJSONResponse{
					Code:    INTERNALERROR,
					Message: "failed to verify fund",
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
		return GetCapTable500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Code:    INTERNALERROR,
				Message: "failed to get cap table",
			},
		}, nil
	}

	entries := make([]CapTableEntry, len(view.Entries))
	for i, e := range view.Entries {
		var percentage float64
		if fundTotalUnits > 0 {
			percentage = float64(e.Units) / float64(fundTotalUnits) * 100
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
	// TODO: Implement when transfer domain is added.
	return ListTransfers200JSONResponse([]Transfer{}), nil
}

// CreateTransfer creates a new transfer.
func (h *APIHandler) CreateTransfer(ctx context.Context, request CreateTransferRequestObject) (CreateTransferResponseObject, error) {
	// TODO: Implement when transfer domain is added.
	return CreateTransfer500JSONResponse{
		InternalErrorJSONResponse: InternalErrorJSONResponse{
			Code:    INTERNALERROR,
			Message: "transfers not yet implemented",
		},
	}, nil
}

// Ensure APIHandler implements StrictServerInterface.
var _ StrictServerInterface = (*APIHandler)(nil)
