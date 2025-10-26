package benchmarks

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/apognu/gocal"
	golangical "github.com/arran4/golang-ical"
	"github.com/michael-gallo/simple-ical/parse"
)

const commonName = "Org"

const (
	singleEventFileName    = "./test_event.ical"
	multipleEventsFileName = "./test_multiple_events.ical"

	// An extremely minimal ical file, with a single event with only required properties
	simpleFileName   = "./test_simple.ical"
	singleFileName   = "./test_event.ical"
	multipleFileName = "./test_multiple_events.ical"
	complexFileName  = "./test_complex.ical"
)

func BenchmarkAllScenarios(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		benchmarkFile(b, simpleFileName, "Simple Event")
	})

	b.Run("Single", func(b *testing.B) {
		benchmarkFile(b, singleFileName, "Single Event")
	})

	b.Run("Multiple", func(b *testing.B) {
		benchmarkFile(b, multipleFileName, "Multiple Events")
	})

	b.Run("Complex", func(b *testing.B) {
		benchmarkFile(b, complexFileName, "Complex Calendar")
	})
}

func benchmarkFile(b *testing.B, fileName string, description string) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		b.Fatalf("Failed to read file %s: %v", fileName, err)
	}

	var reader bytes.Reader
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Reset(fileContent)
		cal, err := parse.IcalReader(&reader)
		if err != nil {
			b.Fatalf("Failed to parse %s: %v", description, err)
		}

		// Basic validation to ensure parsing worked
		if cal == nil {
			b.Fatal("Calendar is nil")
		}

		// Prevent optimization
		_ = cal
	}
}

// Runs benchmarks to test simmple-ical against other parsers
func BenchmarkComparativeAll(b *testing.B) {
	// Pre-load file content to avoid I/O overhead in benchmarks
	benchmarkFileComparison(b, singleEventFileName, "Single Event")
	benchmarkFileComparison(b, multipleEventsFileName, "Multiple Events")
}

func benchmarkFileComparison(b *testing.B, fileName string, testName string) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		panic("Invalid File")
	}
	var reader bytes.Reader

	b.Run(fmt.Sprintf("%s - Gocal", testName), func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			c := gocal.NewParser(&reader)
			c.SkipBounds = true // Parse all events regardless of date
			err := c.Parse()
			if err != nil {
				panic(err)
			}
			if c.Events[0].Organizer.Cn != commonName {
				panic("Invalid Organizer")
			}
		}
	})

	b.Run(fmt.Sprintf("%s - SimpleIcal", testName), func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			cal, err := parse.IcalReader(&reader)
			if err != nil {
				panic(err)
			}
			if cal.Events[0].Organizer.CommonName != commonName {
				panic("Invalid Organizer")
			}
		}
	})

	b.Run(fmt.Sprintf("%s - GolangIcal", testName), func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			cal, err := golangical.ParseCalendar(&reader)
			if err != nil {
				panic(err)
			}
			organizerProp := cal.Events()[0].GetProperty(golangical.ComponentPropertyOrganizer)
			organizerValue := organizerProp.ICalParameters["CN"][0]
			if organizerValue != commonName {
				panic("Invalid organizer value")
			}
		}
	})
}
