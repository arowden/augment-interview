package fund

import (
	"context"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ListParams = validation.ListParams

type ListResult struct {
	Items  []*Fund
	Total  int
	Limit  int
	Offset int
}

type Repository interface {
	Create(ctx context.Context, fund *Fund) error

	CreateTx(ctx context.Context, tx pgx.Tx, fund *Fund) error

	FindByID(ctx context.Context, id uuid.UUID) (*Fund, error)

	List(ctx context.Context, params ListParams) (*ListResult, error)
}
