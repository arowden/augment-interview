package transfer

import (
	"strings"
	"unicode/utf8"

	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/arowden/augment-fund/internal/validation"
)

// Validator validates transfer requests.
type Validator struct{}

// NewValidator creates a new Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateBasic performs format validation without ownership lookup.
// Checks: owner names valid (non-empty, max length), units in valid range, from != to.
func (v *Validator) ValidateBasic(req Request) error {
	fromOwner := strings.TrimSpace(req.FromOwner)
	toOwner := strings.TrimSpace(req.ToOwner)

	if fromOwner == "" || utf8.RuneCountInString(fromOwner) > validation.MaxNameLength {
		return ErrInvalidOwner
	}
	if toOwner == "" || utf8.RuneCountInString(toOwner) > validation.MaxNameLength {
		return ErrInvalidOwner
	}
	if req.Units <= 0 || req.Units > validation.MaxUnits {
		return ErrInvalidUnits
	}
	if fromOwner == toOwner {
		return ErrSelfTransfer
	}
	return nil
}

// Validate performs full validation including ownership checks.
// The fromEntry must be provided (looked up by caller with lock).
func (v *Validator) Validate(req Request, fromEntry *ownership.Entry) error {
	if err := v.ValidateBasic(req); err != nil {
		return err
	}
	if fromEntry == nil {
		return ErrOwnerNotFound
	}
	if fromEntry.Units < req.Units {
		return ErrInsufficientUnits
	}
	return nil
}
