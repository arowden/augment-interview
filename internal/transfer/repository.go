package transfer

import (
	"context"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListParams is an alias for the shared pagination type.
// Deprecated: Use validation.ListParams directly for new code.
type ListParams = validation.ListParams

// TransferList holds a paginated list of transfers.
type TransferList struct {
	Transfers  []*Transfer
	TotalCount int
	Limit      int
	Offset     int
}

// Repository defines the interface for transfer persistence operations.
type Repository interface {
	// Create persists a new transfer record.
	Create(ctx context.Context, transfer *Transfer) error

	// CreateTx persists a new transfer record within the provided transaction.
	CreateTx(ctx context.Context, tx pgx.Tx, transfer *Transfer) error

	// FindByFundID retrieves transfers for a fund with pagination.
	// Transfers are ordered by transferred_at ascending.
	FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*TransferList, error)

	// FindByIdempotencyKey looks up a transfer by its idempotency key within a transaction.
	// Returns nil, nil if not found.
	FindByIdempotencyKey(ctx context.Context, tx pgx.Tx, key uuid.UUID) (*Transfer, error)
}
