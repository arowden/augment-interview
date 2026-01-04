// Package fund provides the Fund domain entity and related operations.
package fund

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

// Fund represents an investment fund with ownership units.
type Fund struct {
	ID         uuid.UUID
	Name       string
	TotalUnits int
	CreatedAt  time.Time
}

// NewFund creates a new Fund with validation.
// Returns ErrInvalidFund if:
//   - name is empty/whitespace or exceeds validation.MaxNameLength
//   - totalUnits is not positive or exceeds validation.MaxUnits
//
// Note: CreatedAt is set to time.Now() at call time. For tests requiring
// deterministic timestamps or ordering by CreatedAt, callers should introduce
// delays between calls or use a clock abstraction for time-sensitive testing.
func NewFund(name string, totalUnits int) (*Fund, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" || utf8.RuneCountInString(trimmedName) > validation.MaxNameLength {
		return nil, ErrInvalidFund
	}
	if totalUnits <= 0 || totalUnits > validation.MaxUnits {
		return nil, ErrInvalidFund
	}
	return &Fund{
		ID:         uuid.New(),
		Name:       trimmedName,
		TotalUnits: totalUnits,
		CreatedAt:  time.Now(),
	}, nil
}
