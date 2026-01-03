package fund

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	// DefaultListLimit is the default number of items returned per page.
	DefaultListLimit = 100
	// MaxListLimit is the maximum allowed items per page.
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

// ListResult contains paginated fund results.
type ListResult struct {
	Items  []*Fund
	Total  int
	Limit  int
	Offset int
}

// Repository defines the interface for fund persistence operations.
type Repository interface {
	// Create persists a new fund to the database.
	Create(ctx context.Context, fund *Fund) error

	// CreateTx persists a new fund within the provided transaction.
	// The caller is responsible for commit/rollback.
	CreateTx(ctx context.Context, tx pgx.Tx, fund *Fund) error

	// FindByID retrieves a fund by its UUID.
	// Returns ErrNotFound (wrapped) if the fund does not exist.
	FindByID(ctx context.Context, id uuid.UUID) (*Fund, error)

	// List retrieves funds with pagination, ordered by created_at descending.
	// Params are normalized: defaults applied, limits enforced.
	// Returns an empty Items slice (not nil) if no funds exist.
	List(ctx context.Context, params ListParams) (*ListResult, error)
}
