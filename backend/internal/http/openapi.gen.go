package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	DUPLICATETRANSFER ErrorCode = "DUPLICATE_TRANSFER"
	FUNDNOTFOUND      ErrorCode = "FUND_NOT_FOUND"
	INSUFFICIENTUNITS ErrorCode = "INSUFFICIENT_UNITS"
	INTERNALERROR     ErrorCode = "INTERNAL_ERROR"
	INVALIDFUND       ErrorCode = "INVALID_FUND"
	INVALIDREQUEST    ErrorCode = "INVALID_REQUEST"
	OWNERNOTFOUND     ErrorCode = "OWNER_NOT_FOUND"
	SELFTRANSFER      ErrorCode = "SELF_TRANSFER"
)

type CapTable struct {
	Entries []CapTableEntry `json:"entries"`

	FundId openapi_types.UUID `json:"fundId"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`

	Total int `json:"total"`
}

type CapTableEntry struct {
	AcquiredAt time.Time `json:"acquiredAt"`

	OwnerName string `json:"ownerName"`

	Percentage float64 `json:"percentage"`

	Units int `json:"units"`
}

type CreateFundRequest struct {
	InitialOwner string `json:"initialOwner"`

	Name string `json:"name"`

	TotalUnits int `json:"totalUnits"`
}

type CreateTransferRequest struct {
	FromOwner string `json:"fromOwner"`

	IdempotencyKey *openapi_types.UUID `json:"idempotencyKey,omitempty"`

	ToOwner string `json:"toOwner"`

	Units int `json:"units"`
}

type Error struct {
	Code ErrorCode `json:"code"`

	Details *map[string]interface{} `json:"details,omitempty"`

	Message string `json:"message"`
}

type ErrorCode string

type Fund struct {
	CreatedAt time.Time `json:"createdAt"`

	Id openapi_types.UUID `json:"id"`

	Name string `json:"name"`

	TotalUnits int `json:"totalUnits"`
}

type FundList struct {
	Funds []Fund `json:"funds"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`

	Total int `json:"total"`
}

type Transfer struct {
	FromOwner string `json:"fromOwner"`

	FundId openapi_types.UUID `json:"fundId"`

	Id openapi_types.UUID `json:"id"`

	ToOwner string `json:"toOwner"`

	TransferredAt time.Time `json:"transferredAt"`

	Units int `json:"units"`
}

type TransferList struct {
	FundId openapi_types.UUID `json:"fundId"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`

	Total int `json:"total"`

	Transfers []Transfer `json:"transfers"`
}

type FundId = openapi_types.UUID

type Limit = int

type Offset = int

type BadRequest = Error

type DuplicateTransfer = Error

type FundNotFound = Error

type InternalError = Error

type TransferBadRequest = Error

type TransferNotFound = Error

type ListFundsParams struct {
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type GetCapTableParams struct {
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type ListTransfersParams struct {
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type CreateFundJSONRequestBody = CreateFundRequest

type CreateTransferJSONRequestBody = CreateTransferRequest

type ServerInterface interface {
	ListFunds(w http.ResponseWriter, r *http.Request, params ListFundsParams)
	CreateFund(w http.ResponseWriter, r *http.Request)
	GetFund(w http.ResponseWriter, r *http.Request, fundId FundId)
	GetCapTable(w http.ResponseWriter, r *http.Request, fundId FundId, params GetCapTableParams)
	ListTransfers(w http.ResponseWriter, r *http.Request, fundId FundId, params ListTransfersParams)
	CreateTransfer(w http.ResponseWriter, r *http.Request, fundId FundId)
	ResetDatabase(w http.ResponseWriter, r *http.Request)
}


type Unimplemented struct{}

func (_ Unimplemented) ListFunds(w http.ResponseWriter, r *http.Request, params ListFundsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) CreateFund(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetFund(w http.ResponseWriter, r *http.Request, fundId FundId) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) GetCapTable(w http.ResponseWriter, r *http.Request, fundId FundId, params GetCapTableParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) ListTransfers(w http.ResponseWriter, r *http.Request, fundId FundId, params ListTransfersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) CreateTransfer(w http.ResponseWriter, r *http.Request, fundId FundId) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) ResetDatabase(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (siw *ServerInterfaceWrapper) ListFunds(w http.ResponseWriter, r *http.Request) {

	var err error

	var params ListFundsParams


	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}


	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListFunds(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) CreateFund(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateFund(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetFund(w http.ResponseWriter, r *http.Request) {

	var err error

	var fundId FundId

	err = runtime.BindStyledParameterWithOptions("simple", "fundId", chi.URLParam(r, "fundId"), &fundId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "fundId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetFund(w, r, fundId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetCapTable(w http.ResponseWriter, r *http.Request) {

	var err error

	var fundId FundId

	err = runtime.BindStyledParameterWithOptions("simple", "fundId", chi.URLParam(r, "fundId"), &fundId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "fundId", Err: err})
		return
	}

	var params GetCapTableParams


	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}


	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetCapTable(w, r, fundId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) ListTransfers(w http.ResponseWriter, r *http.Request) {

	var err error

	var fundId FundId

	err = runtime.BindStyledParameterWithOptions("simple", "fundId", chi.URLParam(r, "fundId"), &fundId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "fundId", Err: err})
		return
	}

	var params ListTransfersParams


	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}


	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListTransfers(w, r, fundId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) CreateTransfer(w http.ResponseWriter, r *http.Request) {

	var err error

	var fundId FundId

	err = runtime.BindStyledParameterWithOptions("simple", "fundId", chi.URLParam(r, "fundId"), &fundId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "fundId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTransfer(w, r, fundId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) ResetDatabase(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ResetDatabase(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/funds", wrapper.ListFunds)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/funds", wrapper.CreateFund)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/funds/{fundId}", wrapper.GetFund)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/funds/{fundId}/cap-table", wrapper.GetCapTable)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/funds/{fundId}/transfers", wrapper.ListTransfers)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/funds/{fundId}/transfers", wrapper.CreateTransfer)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/reset", wrapper.ResetDatabase)
	})

	return r
}

type BadRequestJSONResponse Error

type DuplicateTransferJSONResponse Error

type FundNotFoundJSONResponse Error

type InternalErrorJSONResponse Error

type TransferBadRequestJSONResponse Error

type TransferNotFoundJSONResponse Error

type ListFundsRequestObject struct {
	Params ListFundsParams
}

type ListFundsResponseObject interface {
	VisitListFundsResponse(w http.ResponseWriter) error
}

type ListFunds200JSONResponse FundList

func (response ListFunds200JSONResponse) VisitListFundsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListFunds500JSONResponse struct{ InternalErrorJSONResponse }

func (response ListFunds500JSONResponse) VisitListFundsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateFundRequestObject struct {
	Body *CreateFundJSONRequestBody
}

type CreateFundResponseObject interface {
	VisitCreateFundResponse(w http.ResponseWriter) error
}

type CreateFund201JSONResponse Fund

func (response CreateFund201JSONResponse) VisitCreateFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateFund400JSONResponse struct{ BadRequestJSONResponse }

func (response CreateFund400JSONResponse) VisitCreateFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateFund500JSONResponse struct{ InternalErrorJSONResponse }

func (response CreateFund500JSONResponse) VisitCreateFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetFundRequestObject struct {
	FundId FundId `json:"fundId"`
}

type GetFundResponseObject interface {
	VisitGetFundResponse(w http.ResponseWriter) error
}

type GetFund200JSONResponse Fund

func (response GetFund200JSONResponse) VisitGetFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetFund400JSONResponse struct{ BadRequestJSONResponse }

func (response GetFund400JSONResponse) VisitGetFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetFund404JSONResponse struct{ FundNotFoundJSONResponse }

func (response GetFund404JSONResponse) VisitGetFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetFund500JSONResponse struct{ InternalErrorJSONResponse }

func (response GetFund500JSONResponse) VisitGetFundResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetCapTableRequestObject struct {
	FundId FundId `json:"fundId"`
	Params GetCapTableParams
}

type GetCapTableResponseObject interface {
	VisitGetCapTableResponse(w http.ResponseWriter) error
}

type GetCapTable200JSONResponse CapTable

func (response GetCapTable200JSONResponse) VisitGetCapTableResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetCapTable400JSONResponse struct{ BadRequestJSONResponse }

func (response GetCapTable400JSONResponse) VisitGetCapTableResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetCapTable404JSONResponse struct{ FundNotFoundJSONResponse }

func (response GetCapTable404JSONResponse) VisitGetCapTableResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetCapTable500JSONResponse struct{ InternalErrorJSONResponse }

func (response GetCapTable500JSONResponse) VisitGetCapTableResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type ListTransfersRequestObject struct {
	FundId FundId `json:"fundId"`
	Params ListTransfersParams
}

type ListTransfersResponseObject interface {
	VisitListTransfersResponse(w http.ResponseWriter) error
}

type ListTransfers200JSONResponse TransferList

func (response ListTransfers200JSONResponse) VisitListTransfersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListTransfers400JSONResponse struct{ BadRequestJSONResponse }

func (response ListTransfers400JSONResponse) VisitListTransfersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type ListTransfers404JSONResponse struct{ FundNotFoundJSONResponse }

func (response ListTransfers404JSONResponse) VisitListTransfersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type ListTransfers500JSONResponse struct{ InternalErrorJSONResponse }

func (response ListTransfers500JSONResponse) VisitListTransfersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateTransferRequestObject struct {
	FundId FundId `json:"fundId"`
	Body   *CreateTransferJSONRequestBody
}

type CreateTransferResponseObject interface {
	VisitCreateTransferResponse(w http.ResponseWriter) error
}

type CreateTransfer200JSONResponse Transfer

func (response CreateTransfer200JSONResponse) VisitCreateTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type CreateTransfer201JSONResponse Transfer

func (response CreateTransfer201JSONResponse) VisitCreateTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateTransfer400JSONResponse struct{ TransferBadRequestJSONResponse }

func (response CreateTransfer400JSONResponse) VisitCreateTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateTransfer404JSONResponse struct{ TransferNotFoundJSONResponse }

func (response CreateTransfer404JSONResponse) VisitCreateTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type CreateTransfer409JSONResponse struct{ DuplicateTransferJSONResponse }

func (response CreateTransfer409JSONResponse) VisitCreateTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type CreateTransfer500JSONResponse struct{ InternalErrorJSONResponse }

func (response CreateTransfer500JSONResponse) VisitCreateTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type ResetDatabaseRequestObject struct {
}

type ResetDatabaseResponseObject interface {
	VisitResetDatabaseResponse(w http.ResponseWriter) error
}

type ResetDatabase200JSONResponse struct {
	DeletedFunds     *int    `json:"deletedFunds,omitempty"`
	DeletedOwnership *int    `json:"deletedOwnership,omitempty"`
	DeletedTransfers *int    `json:"deletedTransfers,omitempty"`
	Message          *string `json:"message,omitempty"`
}

func (response ResetDatabase200JSONResponse) VisitResetDatabaseResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ResetDatabase500JSONResponse struct{ InternalErrorJSONResponse }

func (response ResetDatabase500JSONResponse) VisitResetDatabaseResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type StrictServerInterface interface {
	ListFunds(ctx context.Context, request ListFundsRequestObject) (ListFundsResponseObject, error)
	CreateFund(ctx context.Context, request CreateFundRequestObject) (CreateFundResponseObject, error)
	GetFund(ctx context.Context, request GetFundRequestObject) (GetFundResponseObject, error)
	GetCapTable(ctx context.Context, request GetCapTableRequestObject) (GetCapTableResponseObject, error)
	ListTransfers(ctx context.Context, request ListTransfersRequestObject) (ListTransfersResponseObject, error)
	CreateTransfer(ctx context.Context, request CreateTransferRequestObject) (CreateTransferResponseObject, error)
	ResetDatabase(ctx context.Context, request ResetDatabaseRequestObject) (ResetDatabaseResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

func (sh *strictHandler) ListFunds(w http.ResponseWriter, r *http.Request, params ListFundsParams) {
	var request ListFundsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListFunds(ctx, request.(ListFundsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListFunds")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListFundsResponseObject); ok {
		if err := validResponse.VisitListFundsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) CreateFund(w http.ResponseWriter, r *http.Request) {
	var request CreateFundRequestObject

	var body CreateFundJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateFund(ctx, request.(CreateFundRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateFund")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateFundResponseObject); ok {
		if err := validResponse.VisitCreateFundResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetFund(w http.ResponseWriter, r *http.Request, fundId FundId) {
	var request GetFundRequestObject

	request.FundId = fundId

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetFund(ctx, request.(GetFundRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetFund")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetFundResponseObject); ok {
		if err := validResponse.VisitGetFundResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) GetCapTable(w http.ResponseWriter, r *http.Request, fundId FundId, params GetCapTableParams) {
	var request GetCapTableRequestObject

	request.FundId = fundId
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetCapTable(ctx, request.(GetCapTableRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetCapTable")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetCapTableResponseObject); ok {
		if err := validResponse.VisitGetCapTableResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) ListTransfers(w http.ResponseWriter, r *http.Request, fundId FundId, params ListTransfersParams) {
	var request ListTransfersRequestObject

	request.FundId = fundId
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListTransfers(ctx, request.(ListTransfersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListTransfers")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListTransfersResponseObject); ok {
		if err := validResponse.VisitListTransfersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) CreateTransfer(w http.ResponseWriter, r *http.Request, fundId FundId) {
	var request CreateTransferRequestObject

	request.FundId = fundId

	var body CreateTransferJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateTransfer(ctx, request.(CreateTransferRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateTransfer")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateTransferResponseObject); ok {
		if err := validResponse.VisitCreateTransferResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

func (sh *strictHandler) ResetDatabase(w http.ResponseWriter, r *http.Request) {
	var request ResetDatabaseRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ResetDatabase(ctx, request.(ResetDatabaseRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ResetDatabase")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ResetDatabaseResponseObject); ok {
		if err := validResponse.VisitResetDatabaseResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
