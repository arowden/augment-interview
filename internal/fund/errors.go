package fund

import (
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

// ErrNotFound is a generic not found error that can be wrapped with context.
var ErrNotFound = errors.New("not found")

// ErrInvalidFund is returned when fund data fails validation.
var ErrInvalidFund = fmt.Errorf("invalid fund: name must be non-empty (max %d chars) and totalUnits must be positive (max %d)", validation.MaxNameLength, validation.MaxUnits)

// ErrNilFund is returned when a nil fund is passed to a method that requires a valid fund.
var ErrNilFund = errors.New("fund: cannot operate on nil fund")

// ErrDuplicateFundName is returned when attempting to create a fund with a name that already exists.
var ErrDuplicateFundName = errors.New("fund name already exists")

// ErrPoolRequired is returned when a transactional operation is attempted without a pool.
var ErrPoolRequired = errors.New("fund: database pool is required for transactional operations")

// ErrOwnershipRepoRequired is returned when CreateFundWithInitialOwner is called without an ownership repository.
var ErrOwnershipRepoRequired = errors.New("fund: ownership repository is required for fund creation with initial owner")

// NotFoundError returns a wrapped not found error with the fund ID.
func NotFoundError(id uuid.UUID) error {
	return fmt.Errorf("fund %s: %w", id, ErrNotFound)
}
