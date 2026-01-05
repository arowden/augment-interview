package fund

import (
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

var ErrInvalidFund = fmt.Errorf("invalid fund: name must be non-empty (max %d chars) and totalUnits must be positive (max %d)", validation.MaxNameLength, validation.MaxUnits)

var ErrNilFund = errors.New("fund: cannot operate on nil fund")

var ErrDuplicateFundName = errors.New("fund name already exists")

var ErrPoolRequired = errors.New("fund: database pool is required for transactional operations")

var ErrOwnershipRepoRequired = errors.New("fund: ownership repository is required for fund creation with initial owner")

func NotFoundError(id uuid.UUID) error {
	return fmt.Errorf("fund %s: %w", id, ErrNotFound)
}
