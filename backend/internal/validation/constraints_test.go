package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListParams_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		input    ListParams
		expected ListParams
	}{
		{
			name:     "zero values get defaults",
			input:    ListParams{},
			expected: ListParams{Limit: DefaultLimit, Offset: 0},
		},
		{
			name:     "negative limit gets default",
			input:    ListParams{Limit: -1, Offset: 0},
			expected: ListParams{Limit: DefaultLimit, Offset: 0},
		},
		{
			name:     "negative offset gets zero",
			input:    ListParams{Limit: 10, Offset: -5},
			expected: ListParams{Limit: 10, Offset: 0},
		},
		{
			name:     "limit exceeding max gets capped",
			input:    ListParams{Limit: 9999, Offset: 0},
			expected: ListParams{Limit: MaxLimit, Offset: 0},
		},
		{
			name:     "valid params unchanged",
			input:    ListParams{Limit: 50, Offset: 100},
			expected: ListParams{Limit: 50, Offset: 100},
		},
		{
			name:     "limit at max stays at max",
			input:    ListParams{Limit: MaxLimit, Offset: 0},
			expected: ListParams{Limit: MaxLimit, Offset: 0},
		},
		{
			name:     "limit of 1 is valid",
			input:    ListParams{Limit: 1, Offset: 0},
			expected: ListParams{Limit: 1, Offset: 0},
		},
		{
			name:     "both negative values normalized",
			input:    ListParams{Limit: -100, Offset: -50},
			expected: ListParams{Limit: DefaultLimit, Offset: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Normalize()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestListParams_Normalize_DoesNotMutateOriginal(t *testing.T) {
	original := ListParams{Limit: -1, Offset: -5}
	_ = original.Normalize()

	assert.Equal(t, -1, original.Limit)
	assert.Equal(t, -5, original.Offset)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, 1, MinNameLength)
	assert.Equal(t, 255, MaxNameLength)
	assert.Equal(t, 1, MinUnits)
	assert.Equal(t, 2_147_483_647, MaxUnits)
	assert.Equal(t, 100, DefaultLimit)
	assert.Equal(t, 1000, MaxLimit)
	assert.Equal(t, 100.0, PercentageMultiplier)
}
