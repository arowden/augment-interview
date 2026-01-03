package ownership

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	// DefaultListLimit is the default number of entries returned per page.
	DefaultListLimit = 100
	// MaxListLimit is the maximum allowed entries per page.
	MaxListLimit = 1000
)

// ListParams configures pagination for list operations.
type ListParams struct {
	Limit  int
	Offset int
}

// Normalize applies defaults and constraints to ListParams.
// Returns a new ListParams with normalized values; the original is unchanged.
// Use: params = params.Normalize()
func (p ListParams) Normalize() ListParams {
	if p.Limit <= 0 {
		p.Limit = DefaultListLimit
	}
	if p.Limit > MaxListLimit {
		p.Limit = MaxListLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}

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

	// Upsert creates or updates a cap table entry.
	// If the entry exists (by fund_id + owner_name), updates units and updated_at.
	// Preserves the original acquired_at timestamp on update.
	Upsert(ctx context.Context, entry *Entry) error

	// UpsertTx creates or updates a cap table entry within the provided transaction.
	// The caller is responsible for commit/rollback.
	// Preserves the original acquired_at timestamp on update.
	UpsertTx(ctx context.Context, tx pgx.Tx, entry *Entry) error
}
