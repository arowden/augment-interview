package ownership

import (
	"context"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ListParams = validation.ListParams

type Repository interface {
	Create(ctx context.Context, entry *Entry) error

	CreateTx(ctx context.Context, tx pgx.Tx, entry *Entry) error

	FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error)

	FindByFundAndOwner(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error)

	FindByFundAndOwnerForUpdateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string) (*Entry, error)

	DecrementUnitsTx(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, units int) error

	IncrementOrCreateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string, units int) error

	Upsert(ctx context.Context, entry *Entry) error

	UpsertTx(ctx context.Context, tx pgx.Tx, entry *Entry) error
}
