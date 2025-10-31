// Package benchmarks provides comparative benchmarks against other Go iCalendar parsers
package benchmarks

import (
	"testing"

	"github.com/michael-gallo/simpleical/rrule"
	rrule_go "github.com/teambition/rrule-go"
)

func BenchmarkParseRRule(b *testing.B) {
	rruleTests := []struct {
		name  string
		input string
	}{
		{
			name:  "Simple rule with count",
			input: "FREQ=DAILY;INTERVAL=1;COUNT=10",
		},
		{
			name:  "Simple rule with until",
			input: "FREQ=DAILY;INTERVAL=1;UNTIL=20250928T183000Z",
		},
		{
			name:  "String from teambition's rrule.go example",
			input: "FREQ=DAILY;DTSTART=20060101T150405Z;COUNT=5",
		},
		{
			name:  "Every 20th Monday of the year, forever",
			input: "FREQ=YEARLY;BYDAY=20MO",
		},
	}
	for _, test := range rruleTests {
		b.Run(test.name, func(b *testing.B) {
			benchmarkRrule(b, test.input)
		})
	}
}

func benchmarkRrule(b *testing.B, rruleString string) {
	b.Run("SimpleIcal", func(b *testing.B) {
		for b.Loop() {
			_, err := rrule.ParseRRule(rruleString)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("RRuleGo", func(b *testing.B) {
		for b.Loop() {
			_, err := rrule_go.StrToRRule(rruleString)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
