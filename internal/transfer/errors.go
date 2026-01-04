// Package transfer provides transfer domain entities and operations.
package transfer

import (
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/validation"
)

// ErrInsufficientUnits is returned when the from_owner doesn't have enough units.
var ErrInsufficientUnits = fmt.Errorf("insufficient units for transfer")

// ErrSelfTransfer is returned when from_owner equals to_owner.
var ErrSelfTransfer = fmt.Errorf("cannot transfer to self")

// ErrInvalidUnits is returned when units is zero, negative, or exceeds max.
var ErrInvalidUnits = fmt.Errorf("units must be between %d and %d", validation.MinUnits, validation.MaxUnits)

// ErrInvalidOwner is returned when owner name is empty, whitespace, or exceeds max length.
var ErrInvalidOwner = fmt.Errorf("owner name must be non-empty (max %d chars)", validation.MaxNameLength)

// ErrOwnerNotFound is returned when the from_owner doesn't exist in the cap table.
var ErrOwnerNotFound = errors.New("owner not found")

// ErrNilTransfer is returned when a nil transfer is passed to a method.
var ErrNilTransfer = errors.New("transfer: cannot operate on nil transfer")

// ErrDuplicateIdempotencyKey is returned when an idempotency key is reused with different transfer data.
var ErrDuplicateIdempotencyKey = errors.New("idempotency key already used with different transfer data")
