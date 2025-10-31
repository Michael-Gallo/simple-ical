package icaldur

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseIcalTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        time.Time
		expectError bool
	}{
		{
			name:        "Valid time with Z",
			input:       "20250928T183000Z",
			want:        time.Date(2025, 9, 28, 18, 30, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Valid time without Z",
			input:       "20240101T000000",
			want:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			input:       "20231231T235959Z",
			want:        time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			expectError: false,
		},
		{
			input:       "20000101T120000Z",
			want:        time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Invalid time with Z",
			input:       "20250928T1830Z",
			expectError: true,
		},
		{
			name:        "Invalid Format",
			input:       "2025-09-28T18:30:00Z",
			expectError: true,
		},
		{
			name:        "Empty input",
			input:       "",
			expectError: true,
		},
		{
			name:        "Invalid input",
			input:       "invalid",
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseIcalTime(test.input)
			if test.expectError {
				assert.Error(t, err, "expected error for input: %s", test.input)
				return
			}
			assert.NoError(t, err, "unexpected error for input: %s", test.input)
			assert.Equal(t, test.want, got, "mismatch for input: %s", test.input)
		})
	}
}

func BenchmarkParseIcalTime(b *testing.B) {
	times := []string{
		"20250928T183000Z",
		"20240101T000000Z",
		"20231231T235959Z",
		"20000101T120000",
	}
	for b.Loop() {
		for _, t := range times {
			_, err := ParseIcalTime(t)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
