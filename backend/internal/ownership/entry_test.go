package ownership

import (
	"strings"
	"testing"

	"github.com/arowden/augment-fund/internal/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCapTableEntry(t *testing.T) {
	fundID := uuid.New()

	t.Run("valid inputs", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "John Doe", 1000)
		require.NoError(t, err)
		assert.NotEmpty(t, entry.ID)
		assert.Equal(t, fundID, entry.FundID)
		assert.Equal(t, "John Doe", entry.OwnerName)
		assert.Equal(t, 1000, entry.Units)
		assert.False(t, entry.AcquiredAt.IsZero())
		assert.False(t, entry.UpdatedAt.IsZero())
		assert.Nil(t, entry.DeletedAt)
	})

	t.Run("trims whitespace from owner name", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "  Jane Doe  ", 500)
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", entry.OwnerName)
	})

	t.Run("empty owner name returns error", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "", 1000)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("whitespace-only owner name returns error", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "   ", 1000)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("zero units is valid", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "Empty Holder", 0)
		require.NoError(t, err)
		assert.Equal(t, 0, entry.Units)
	})

	t.Run("negative units returns error", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "Test Owner", -100)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidUnits)
	})

	t.Run("owner name exceeding max length returns error", func(t *testing.T) {
		longName := strings.Repeat("A", validation.MaxNameLength+1)
		entry, err := NewCapTableEntry(fundID, longName, 1000)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("owner name at max length succeeds", func(t *testing.T) {
		maxName := strings.Repeat("A", validation.MaxNameLength)
		entry, err := NewCapTableEntry(fundID, maxName, 1000)
		require.NoError(t, err)
		assert.Equal(t, maxName, entry.OwnerName)
	})

	t.Run("unicode owner name counts runes not bytes", func(t *testing.T) {
		unicodeName := strings.Repeat("基", validation.MaxNameLength)
		entry, err := NewCapTableEntry(fundID, unicodeName, 1000)
		require.NoError(t, err)
		assert.Equal(t, unicodeName, entry.OwnerName)
	})

	t.Run("unicode owner name exceeding max runes returns error", func(t *testing.T) {
		unicodeName := strings.Repeat("基", validation.MaxNameLength+1)
		entry, err := NewCapTableEntry(fundID, unicodeName, 1000)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, ErrInvalidOwner)
	})

	t.Run("acquiredAt equals updatedAt on creation", func(t *testing.T) {
		entry, err := NewCapTableEntry(fundID, "New Owner", 500)
		require.NoError(t, err)
		assert.Equal(t, entry.AcquiredAt, entry.UpdatedAt)
	})
}
