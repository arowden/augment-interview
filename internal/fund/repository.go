package fund

import (
	"context"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListParams is an alias for the shared pagination type.
// Deprecated: Use validation.ListParams directly for new code.
type ListParams = validation.ListParams

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
