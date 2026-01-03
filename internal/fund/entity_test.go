package fund

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFund(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		fund, err := NewFund("Growth Fund", 1000)
		require.NoError(t, err)
		assert.NotEmpty(t, fund.ID)
		assert.Equal(t, "Growth Fund", fund.Name)
		assert.Equal(t, 1000, fund.TotalUnits)
		assert.False(t, fund.CreatedAt.IsZero())
	})

	t.Run("trims whitespace from name", func(t *testing.T) {
		fund, err := NewFund("  Trimmed Fund  ", 500)
		require.NoError(t, err)
		assert.Equal(t, "Trimmed Fund", fund.Name)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		fund, err := NewFund("", 1000)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("whitespace-only name returns error", func(t *testing.T) {
		fund, err := NewFund("   ", 1000)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("zero totalUnits returns error", func(t *testing.T) {
		fund, err := NewFund("Test Fund", 0)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("negative totalUnits returns error", func(t *testing.T) {
		fund, err := NewFund("Test Fund", -100)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("name exceeding max length returns error", func(t *testing.T) {
		longName := strings.Repeat("A", MaxNameLength+1)
		fund, err := NewFund(longName, 1000)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("name at max length succeeds", func(t *testing.T) {
		maxName := strings.Repeat("A", MaxNameLength)
		fund, err := NewFund(maxName, 1000)
		require.NoError(t, err)
		assert.Equal(t, maxName, fund.Name)
	})

	t.Run("totalUnits exceeding max returns error", func(t *testing.T) {
		fund, err := NewFund("Test Fund", MaxTotalUnits+1)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})

	t.Run("totalUnits at max succeeds", func(t *testing.T) {
		fund, err := NewFund("Test Fund", MaxTotalUnits)
		require.NoError(t, err)
		assert.Equal(t, MaxTotalUnits, fund.TotalUnits)
	})

	t.Run("unicode name counts runes not bytes", func(t *testing.T) {
		// 255 CJK characters (each is 3 bytes in UTF-8, so 765 bytes total)
		unicodeName := strings.Repeat("基", MaxNameLength)
		fund, err := NewFund(unicodeName, 1000)
		require.NoError(t, err)
		assert.Equal(t, unicodeName, fund.Name)
	})

	t.Run("unicode name exceeding max runes returns error", func(t *testing.T) {
		// 256 CJK characters should fail
		unicodeName := strings.Repeat("基", MaxNameLength+1)
		fund, err := NewFund(unicodeName, 1000)
		assert.Nil(t, fund)
		assert.ErrorIs(t, err, ErrInvalidFund)
	})
}
