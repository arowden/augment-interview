package transfer

import (
	"context"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ListParams = validation.ListParams

type TransferList struct {
	Transfers  []*Transfer
	TotalCount int
	Limit      int
	Offset     int
}

type Repository interface {
	Create(ctx context.Context, transfer *Transfer) error

	CreateTx(ctx context.Context, tx pgx.Tx, transfer *Transfer) error

	FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*TransferList, error)

	FindByIdempotencyKey(ctx context.Context, tx pgx.Tx, key uuid.UUID) (*Transfer, error)
}
