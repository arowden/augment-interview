// Package fund provides the Fund domain entity and related operations.
package fund

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	// MaxNameLength is the maximum allowed length for a fund name.
	MaxNameLength = 255
	// MaxTotalUnits is the maximum allowed value for total units (PostgreSQL INTEGER max).
	MaxTotalUnits = 2_147_483_647
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
//   - name is empty/whitespace or exceeds MaxNameLength
//   - totalUnits is not positive or exceeds MaxTotalUnits
func NewFund(name string, totalUnits int) (*Fund, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" || utf8.RuneCountInString(trimmedName) > MaxNameLength {
		return nil, ErrInvalidFund
	}
	if totalUnits <= 0 || totalUnits > MaxTotalUnits {
		return nil, ErrInvalidFund
	}
	return &Fund{
		ID:         uuid.New(),
		Name:       trimmedName,
		TotalUnits: totalUnits,
		CreatedAt:  time.Now(),
	}, nil
}
