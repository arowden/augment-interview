package ownership

import (
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

// ErrNotFound is a generic not found error that can be wrapped with context.
var ErrNotFound = errors.New("not found")

// ErrOwnerNotFound is returned when an owner is not found in the cap table.
var ErrOwnerNotFound = errors.New("owner not found")

// ErrInvalidOwner is returned when owner name fails validation.
var ErrInvalidOwner = fmt.Errorf("invalid owner: name must be non-empty (max %d chars)", validation.MaxNameLength)

// ErrInvalidUnits is returned when units value fails validation.
var ErrInvalidUnits = fmt.Errorf("invalid units: must be between 0 and %d", validation.MaxUnits)

// ErrNilEntry is returned when a nil entry is passed to a method that requires a valid entry.
var ErrNilEntry = errors.New("ownership: cannot operate on nil entry")

// OwnerNotFoundError returns a wrapped owner not found error with fund and owner context.
func OwnerNotFoundError(fundID uuid.UUID, ownerName string) error {
	return fmt.Errorf("owner %q in fund %s: %w", ownerName, fundID, ErrOwnerNotFound)
}

// NotFoundError returns a wrapped not found error with the entry ID.
func NotFoundError(id uuid.UUID) error {
	return fmt.Errorf("entry %s: %w", id, ErrNotFound)
}
