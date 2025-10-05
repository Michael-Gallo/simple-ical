// Package benchmarks provides comparative benchmarks against other Go iCalendar parsers
package benchmarks

import (
	"testing"

	"github.com/michael-gallo/simple-ical/rrule"
	rrule_go "github.com/teambition/rrule-go"
)

func BenchmarkParseRRule(b *testing.B) {
	rruleStrings := []string{
		// Simple rule with count
		"FREQ=DAILY;INTERVAL=1;COUNT=10",
		// simple rule with until
		"FREQ=DAILY;INTERVAL=1;UNTIL=20250928T183000Z",
		// String from teambition's rrule.go example
		"FREQ=DAILY;DTSTART=20060101T150405Z;COUNT=5",
	}

	for _, rruleString := range rruleStrings {
		benchmarkRrule(b, rruleString)
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
