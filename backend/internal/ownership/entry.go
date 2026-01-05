package ownership

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

type Entry struct {
	ID         uuid.UUID
	FundID     uuid.UUID
	OwnerName  string
	Units      int
	AcquiredAt time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

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
