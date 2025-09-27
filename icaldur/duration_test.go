package icaldur

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseICalDuration(t *testing.T) {
	tests := []struct {
		input       string
		want        time.Duration
		expectError error
	}{
		{input: "PT1H", want: time.Hour},
		{input: "PT1M", want: time.Minute},
		{input: "PT1S", want: time.Second},
		{input: "PT1H30M", want: time.Hour + time.Minute*30},
		{input: "PT1H30M1S", want: time.Hour + time.Minute*30 + time.Second},
		{input: "P15DT5H0M20S", want: time.Hour*24*15 + time.Hour*5 + time.Minute*0 + time.Second*20},
		{input: "+P15DT5H0M20S", want: time.Hour*24*15 + time.Hour*5 + time.Minute*0 + time.Second*20},
		{input: "-P15DT5H0M20S", want: -(time.Hour*24*15 + time.Hour*5 + time.Minute*0 + time.Second*20)},
		{input: "", want: 0, expectError: ErrEmpty},
		{input: "+Q15DT5H0M20S", expectError: ErrBadPrefix},
		{input: "+P15DT5H0M20G", expectError: ErrUnexpectedChar},
		{input: "+P15DT5H0M20", expectError: ErrMissingUnit},
		{input: "+P15DT5H0M20S20S", expectError: ErrDuplicateUnit},
	}
	for _, test := range tests {
		got, err := ParseICalDuration(test.input)
		if test.expectError != nil {
			assert.ErrorIs(t, err, test.expectError)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, test.want, got)
	}
}

func BenchmarkParseICalDuration(b *testing.B) {
	for b.Loop() {
		_, err := ParseICalDuration("P15DT5H0M20S")
		if err != nil {
			b.Fatal(err)
		}
	}
}
