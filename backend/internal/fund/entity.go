package fund

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
)

type Fund struct {
	ID         uuid.UUID
	Name       string
	TotalUnits int
	CreatedAt  time.Time
}

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
