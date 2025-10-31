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
	testCases := []struct {
		fileName string
		testName string
	}{
		{simpleFileName, "Simple Event"},
		{singleFileName, "Single Event"},
		{multipleFileName, "Multiple Events"},
		{complexFileName, "Complex Calendar"},
	}
	for _, testCase := range testCases {
		b.Run(testCase.testName, func(b *testing.B) {
			benchmarkFile(b, testCase.fileName, testCase.testName)
		})
	}
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

// Runs benchmarks to test simple-ical against other parsers
func BenchmarkComparativeAll(b *testing.B) {
	testCases := []struct {
		fileName string
		testName string
	}{
		{singleEventFileName, "Single Event"},
		{multipleEventsFileName, "Multiple Events"},
		{simpleFileName, "Simple Event"},
		{complexFileName, "Complex Calendar"},
	}

	for _, testCase := range testCases {
		benchmarkFileComparison(b, testCase.fileName, testCase.testName)
	}
}

func benchmarkFileComparison(b *testing.B, fileName string, testName string) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		panic("Invalid File")
	}
	var reader bytes.Reader

	b.Run(fmt.Sprintf("%s - SimpleIcal", testName), func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			_, err := parse.IcalReader(&reader)
			if err != nil {
				panic(err)
			}
		}
	})

	b.Run(fmt.Sprintf("%s - Gocal", testName), func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			c := gocal.NewParser(&reader)
			c.SkipBounds = true // Parse all events regardless of date
			err := c.Parse()
			if err != nil {
				panic(err)
			}
		}
	})

	b.Run(fmt.Sprintf("%s - GolangIcal", testName), func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			_, err := golangical.ParseCalendar(&reader)
			if err != nil {
				panic(err)
			}
		}
	})
}
