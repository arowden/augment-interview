// Package ownership provides the cap table domain entities and operations.
package ownership

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

// Entry represents a single ownership record in the cap table.
type Entry struct {
	ID         uuid.UUID
	FundID     uuid.UUID
	OwnerName  string
	Units      int
	AcquiredAt time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

// NewCapTableEntry creates a new cap table entry with validation.
// Returns an error if:
//   - ownerName is empty/whitespace or exceeds validation.MaxNameLength
//   - units is negative or exceeds validation.MaxUnits (zero is valid for sold-out positions)
//
// The ownerName is trimmed of leading/trailing whitespace.
// AcquiredAt and UpdatedAt are set to time.Now() at call time.
func NewCapTableEntry(fundID uuid.UUID, ownerName string, units int) (*Entry, error) {
	trimmedName := strings.TrimSpace(ownerName)
	if trimmedName == "" || utf8.RuneCountInString(trimmedName) > validation.MaxNameLength {
		return nil, ErrInvalidOwner
	}
	if units < 0 || units > validation.MaxUnits {
		return nil, ErrInvalidUnits
	}

	now := time.Now()
	return &Entry{
		ID:         uuid.New(),
		FundID:     fundID,
		OwnerName:  trimmedName,
		Units:      units,
		AcquiredAt: now,
		UpdatedAt:  now,
		DeletedAt:  nil,
	}, nil
}
