package transfer

import (
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/validation"
)

var ErrInsufficientUnits = fmt.Errorf("insufficient units for transfer")

var ErrSelfTransfer = fmt.Errorf("cannot transfer to self")

var ErrInvalidUnits = fmt.Errorf("units must be between %d and %d", validation.MinUnits, validation.MaxUnits)

var ErrInvalidOwner = fmt.Errorf("owner name must be non-empty (max %d chars)", validation.MaxNameLength)

var ErrOwnerNotFound = errors.New("owner not found")

var ErrNilTransfer = errors.New("transfer: cannot operate on nil transfer")

var ErrDuplicateIdempotencyKey = errors.New("idempotency key already used with different transfer data")
