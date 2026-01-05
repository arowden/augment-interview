package transfer

import (
	"testing"

	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateBasic(t *testing.T) {
	v := NewValidator()
	fundID := uuid.New()

	t.Run("valid request passes validation", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.NoError(t, err)
	})

	t.Run("empty from_owner returns ErrInvalidOwner", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "",
			ToOwner:   "Bob",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("whitespace-only from_owner returns ErrInvalidOwner", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "   ",
			ToOwner:   "Bob",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("empty to_owner returns ErrInvalidOwner", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("whitespace-only to_owner returns ErrInvalidOwner", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "   ",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("zero units returns ErrInvalidUnits", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     0,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrInvalidUnits)
	})

	t.Run("negative units returns ErrInvalidUnits", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     -50,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrInvalidUnits)
	})

	t.Run("same from and to owner returns ErrSelfTransfer", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Alice",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrSelfTransfer)
	})

	t.Run("different case owners are distinct", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "ALICE",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.NoError(t, err)
	})

	t.Run("self-transfer check trims whitespace", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "  Alice  ",
			ToOwner:   "Alice",
			Units:     100,
		}
		err := v.ValidateBasic(req)
		assert.ErrorIs(t, err, ErrSelfTransfer)
	})
}

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()
	fundID := uuid.New()

	t.Run("valid request with sufficient units passes", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}
		fromEntry := &ownership.Entry{
			ID:        uuid.New(),
			FundID:    fundID,
			OwnerName: "Alice",
			Units:     500,
		}
		err := v.Validate(req, fromEntry)
		assert.NoError(t, err)
	})

	t.Run("nil fromEntry returns ErrOwnerNotFound", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}
		err := v.Validate(req, nil)
		assert.ErrorIs(t, err, ErrOwnerNotFound)
	})

	t.Run("insufficient units returns ErrInsufficientUnits", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     500,
		}
		fromEntry := &ownership.Entry{
			ID:        uuid.New(),
			FundID:    fundID,
			OwnerName: "Alice",
			Units:     100,
		}
		err := v.Validate(req, fromEntry)
		assert.ErrorIs(t, err, ErrInsufficientUnits)
	})

	t.Run("exact units allowed", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "Alice",
			ToOwner:   "Bob",
			Units:     100,
		}
		fromEntry := &ownership.Entry{
			ID:        uuid.New(),
			FundID:    fundID,
			OwnerName: "Alice",
			Units:     100,
		}
		err := v.Validate(req, fromEntry)
		assert.NoError(t, err)
	})

	t.Run("basic validation errors still surface", func(t *testing.T) {
		req := Request{
			FundID:    fundID,
			FromOwner: "",
			ToOwner:   "Bob",
			Units:     100,
		}
		fromEntry := &ownership.Entry{
			ID:        uuid.New(),
			FundID:    fundID,
			OwnerName: "Alice",
			Units:     500,
		}
		err := v.Validate(req, fromEntry)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})
}
