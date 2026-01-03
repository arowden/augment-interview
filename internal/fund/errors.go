package fund

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ErrNotFound is a generic not found error that can be wrapped with context.
var ErrNotFound = errors.New("not found")

// ErrInvalidFund is returned when fund data fails validation.
var ErrInvalidFund = errors.New("invalid fund: name must be non-empty (max 255 chars) and totalUnits must be positive (max 2147483647)")

// ErrNilFund is returned when a nil fund is passed to a method that requires a valid fund.
var ErrNilFund = errors.New("fund: cannot operate on nil fund")

// ErrDuplicateFundName is returned when attempting to create a fund with a name that already exists.
var ErrDuplicateFundName = errors.New("fund name already exists")

// NotFoundError returns a wrapped not found error with the fund ID.
func NotFoundError(id uuid.UUID) error {
	return fmt.Errorf("fund %s: %w", id, ErrNotFound)
}
