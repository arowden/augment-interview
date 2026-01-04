package ownership

import (
	"context"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListParams is an alias for the shared pagination type.
// Deprecated: Use validation.ListParams directly for new code.
type ListParams = validation.ListParams

// Repository defines the interface for cap table entry persistence operations.
type Repository interface {
	// Create persists a new cap table entry to the database.
	Create(ctx context.Context, entry *Entry) error

	// CreateTx persists a new cap table entry within the provided transaction.
	// The caller is responsible for commit/rollback.
	CreateTx(ctx context.Context, tx pgx.Tx, entry *Entry) error

	// FindByFundID retrieves all cap table entries for a fund with pagination.
	// Entries are ordered by units descending.
	// Returns an empty slice (not nil) if no entries exist.
	FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error)

	// FindByFundAndOwner retrieves a single cap table entry by fund and owner name.
	// Returns ErrOwnerNotFound (wrapped) if the entry does not exist.
	FindByFundAndOwner(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error)

	// FindByFundAndOwnerForUpdateTx retrieves and locks a cap table entry within a transaction.
	// Uses SELECT FOR UPDATE to prevent concurrent modifications.
	// Returns ErrOwnerNotFound (wrapped) if the entry does not exist.
	FindByFundAndOwnerForUpdateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string) (*Entry, error)

	// DecrementUnitsTx decreases an owner's units within a transaction.
	// The caller must ensure sufficient units exist before calling.
	DecrementUnitsTx(ctx context.Context, tx pgx.Tx, entryID uuid.UUID, units int) error

	// IncrementOrCreateTx adds units to an existing owner or creates a new entry.
	// Uses upsert to handle both cases atomically.
	IncrementOrCreateTx(ctx context.Context, tx pgx.Tx, fundID uuid.UUID, ownerName string, units int) error

	// Upsert creates or updates a cap table entry.
	// If the entry exists (by fund_id + owner_name), updates units and updated_at.
	// Preserves the original acquired_at timestamp on update.
	Upsert(ctx context.Context, entry *Entry) error

	// UpsertTx creates or updates a cap table entry within the provided transaction.
	// The caller is responsible for commit/rollback.
	// Preserves the original acquired_at timestamp on update.
	UpsertTx(ctx context.Context, tx pgx.Tx, entry *Entry) error
}
