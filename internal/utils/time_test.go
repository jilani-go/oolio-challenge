package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestToMidnightUTC tests the NormalizeToMidnightUTC function
func TestToMidnightUTC(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Mid-day time - Should normalize to midnight",
			input:    time.Date(2023, 1, 15, 12, 30, 45, 500000000, time.UTC),
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Already midnight - Should remain unchanged",
			input:    time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Late evening - Should normalize to midnight",
			input:    time.Date(2023, 1, 15, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Early morning - Should normalize to midnight",
			input:    time.Date(2023, 1, 15, 1, 1, 1, 1, time.UTC),
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Different timezone - Should convert to UTC midnight",
			input:    time.Date(2023, 1, 15, 10, 0, 0, 0, time.FixedZone("EST", -5*60*60)),
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result := NormalizeToMidnightUTC(tt.input)

			// Assert the result
			assert.True(t, tt.expected.Equal(result))
			assert.Equal(t, 0, result.Hour())
			assert.Equal(t, 0, result.Minute())
			assert.Equal(t, 0, result.Second())
			assert.Equal(t, 0, result.Nanosecond())
			assert.Equal(t, time.UTC, result.Location())
		})
	}
}
