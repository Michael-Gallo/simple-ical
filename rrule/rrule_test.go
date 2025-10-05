package rrule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: replace with calls to New once go 1.26 is released
func getPointer[T any](v T) *T {
	return &v
}

func TestParseRRule(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *RRule
		expectError error
	}{

		{
			name:  "Valid daily rule",
			input: "FREQ=DAILY;INTERVAL=1;COUNT=10",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  1,
				Count:     getPointer(10),
				Until:     nil,
			},
			expectError: nil,
		},
		{
			name:        "Invalid rule: missing frequency",
			input:       "INTERVAL=1;COUNT=10",
			want:        nil,
			expectError: ErrFrequencyRequired,
		},
		{
			name:        "Invalid rule: count and until cannot both be set",
			input:       "FREQ=DAILY;COUNT=10;UNTIL=19730429T070000Z",
			want:        nil,
			expectError: ErrCountAndUntilBothSet,
		},
		{
			name:        "Invalid rule: interval must be a positive integer",
			input:       "FREQ=DAILY;INTERVAL=0;COUNT=10",
			want:        nil,
			expectError: ErrInvalidInterval,
		},
		{
			name:        "Invalid rule: malformed rrule string",
			input:       "FREQ=DAILY;INVALID",
			want:        nil,
			expectError: ErrInvalidRRuleString,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule, err := ParseRRule(test.input)
			if test.expectError != nil {
				assert.ErrorIs(t, err, test.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.want, rule)
		})

	}
}
