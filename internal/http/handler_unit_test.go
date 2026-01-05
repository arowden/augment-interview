package http

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIHandler(t *testing.T) {
	t.Run("creates handler with no options", func(t *testing.T) {
		h := NewAPIHandler()
		assert.NotNil(t, h)
		assert.Nil(t, h.fundService)
		assert.Nil(t, h.ownershipService)
		assert.Nil(t, h.transferService)
	})
}

func TestNewAPIHandlerStrict(t *testing.T) {
	t.Run("returns error when fund service missing", func(t *testing.T) {
		h, err := NewAPIHandlerStrict()
		assert.Nil(t, h)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fund service")
	})

	t.Run("returns error when ownership service missing", func(t *testing.T) {
		h, err := NewAPIHandlerStrict()
		assert.Nil(t, h)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ownership service")
	})

	t.Run("returns error when transfer service missing", func(t *testing.T) {
		h, err := NewAPIHandlerStrict()
		assert.Nil(t, h)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transfer service")
	})
}

func TestErrorDetails(t *testing.T) {
	t.Run("returns nil for empty context and no extras", func(t *testing.T) {
		details := errorDetails(context.Background(), nil)
		assert.Nil(t, details)
	})

	t.Run("includes extra details", func(t *testing.T) {
		extras := map[string]interface{}{"fundId": "123"}
		details := errorDetails(context.Background(), extras)
		require.NotNil(t, details)
		assert.Equal(t, "123", (*details)["fundId"])
	})

	t.Run("merges multiple extra details", func(t *testing.T) {
		extras := map[string]interface{}{
			"fundId":  "123",
			"ownerId": "456",
		}
		details := errorDetails(context.Background(), extras)
		require.NotNil(t, details)
		assert.Equal(t, "123", (*details)["fundId"])
		assert.Equal(t, "456", (*details)["ownerId"])
	})
}

func TestListFunds_NilService(t *testing.T) {
	h := NewAPIHandler()

	resp, err := h.ListFunds(context.Background(), ListFundsRequestObject{})
	require.NoError(t, err)

	errResp, ok := resp.(ListFunds500JSONResponse)
	require.True(t, ok)
	assert.Equal(t, INTERNALERROR, errResp.Code)
	assert.Contains(t, errResp.Message, "fund service not configured")
}

func TestCreateFund_NilService(t *testing.T) {
	h := NewAPIHandler()

	resp, err := h.CreateFund(context.Background(), CreateFundRequestObject{
		Body: &CreateFundJSONRequestBody{
			Name:         "Test Fund",
			TotalUnits:   1000,
			InitialOwner: "Owner",
		},
	})
	require.NoError(t, err)

	errResp, ok := resp.(CreateFund500JSONResponse)
	require.True(t, ok)
	assert.Equal(t, INTERNALERROR, errResp.Code)
	assert.Contains(t, errResp.Message, "fund service not configured")
}


func TestGetFund_NilService(t *testing.T) {
	h := NewAPIHandler()

	resp, err := h.GetFund(context.Background(), GetFundRequestObject{})
	require.NoError(t, err)

	errResp, ok := resp.(GetFund500JSONResponse)
	require.True(t, ok)
	assert.Equal(t, INTERNALERROR, errResp.Code)
	assert.Contains(t, errResp.Message, "fund service not configured")
}

func TestGetCapTable_NilService(t *testing.T) {
	h := NewAPIHandler()

	resp, err := h.GetCapTable(context.Background(), GetCapTableRequestObject{})
	require.NoError(t, err)

	errResp, ok := resp.(GetCapTable500JSONResponse)
	require.True(t, ok)
	assert.Equal(t, INTERNALERROR, errResp.Code)
	assert.Contains(t, errResp.Message, "ownership service not configured")
}

func TestListTransfers_NilService(t *testing.T) {
	h := NewAPIHandler()

	resp, err := h.ListTransfers(context.Background(), ListTransfersRequestObject{})
	require.NoError(t, err)

	errResp, ok := resp.(ListTransfers500JSONResponse)
	require.True(t, ok)
	assert.Equal(t, INTERNALERROR, errResp.Code)
	assert.Contains(t, errResp.Message, "transfer service not configured")
}

func TestCreateTransfer_NilService(t *testing.T) {
	h := NewAPIHandler()

	resp, err := h.CreateTransfer(context.Background(), CreateTransferRequestObject{
		Body: &CreateTransferJSONRequestBody{
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		},
	})
	require.NoError(t, err)

	errResp, ok := resp.(CreateTransfer500JSONResponse)
	require.True(t, ok)
	assert.Equal(t, INTERNALERROR, errResp.Code)
	assert.Contains(t, errResp.Message, "transfer service not configured")
}


func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(logger)

	t.Run("logs error without request ID", func(t *testing.T) {
		buf.Reset()
		ctx := context.Background()
		testErr := errors.New("test error")

		logError(ctx, "test message", testErr)

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "test error")
		assert.NotContains(t, output, "requestId")
	})

	t.Run("logs error with request ID", func(t *testing.T) {
		buf.Reset()
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "req-123")
		testErr := errors.New("another error")

		logError(ctx, "request failed", testErr)

		output := buf.String()
		assert.Contains(t, output, "request failed")
		assert.Contains(t, output, "another error")
		assert.Contains(t, output, "requestId")
		assert.Contains(t, output, "req-123")
	})

	t.Run("logs error with extra attributes", func(t *testing.T) {
		buf.Reset()
		ctx := context.Background()
		testErr := errors.New("detail error")

		logError(ctx, "operation failed", testErr,
			slog.String("fundId", "fund-456"),
			slog.Int("units", 100),
		)

		output := buf.String()
		assert.Contains(t, output, "operation failed")
		assert.Contains(t, output, "detail error")
		assert.Contains(t, output, "fundId")
		assert.Contains(t, output, "fund-456")
		assert.Contains(t, output, "units")
		assert.Contains(t, output, "100")
	})
}

func TestErrorDetails_WithRequestID(t *testing.T) {
	t.Run("includes request ID from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "req-789")
		details := errorDetails(ctx, nil)
		require.NotNil(t, details)
		assert.Equal(t, "req-789", (*details)["requestId"])
	})

	t.Run("merges request ID with extras", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "req-abc")
		extras := map[string]interface{}{"fundId": "fund-123"}
		details := errorDetails(ctx, extras)
		require.NotNil(t, details)
		assert.Equal(t, "req-abc", (*details)["requestId"])
		assert.Equal(t, "fund-123", (*details)["fundId"])
	})
}

func TestWithTransferService(t *testing.T) {
	t.Run("sets transfer service on handler", func(t *testing.T) {
		h := NewAPIHandler(WithTransferService(nil))
		assert.Nil(t, h.transferService)
	})
}

func TestNewAPIHandlerStrict_AllServices(t *testing.T) {
	t.Run("returns error listing all missing services", func(t *testing.T) {
		h, err := NewAPIHandlerStrict()
		assert.Nil(t, h)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fund service")
		assert.Contains(t, err.Error(), "ownership service")
		assert.Contains(t, err.Error(), "transfer service")
	})
}
