package ownership

import (
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

var ErrOwnerNotFound = errors.New("owner not found")

var ErrInvalidOwner = fmt.Errorf("invalid owner: name must be non-empty (max %d chars)", validation.MaxNameLength)

var ErrInvalidUnits = fmt.Errorf("invalid units: must be between 0 and %d", validation.MaxUnits)

var ErrNilEntry = errors.New("ownership: cannot operate on nil entry")

func OwnerNotFoundError(fundID uuid.UUID, ownerName string) error {
	return fmt.Errorf("owner %q in fund %s: %w", ownerName, fundID, ErrOwnerNotFound)
}

func NotFoundError(id uuid.UUID) error {
	return fmt.Errorf("entry %s: %w", id, ErrNotFound)
}
