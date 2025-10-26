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
)

func BenchmarkSimpleIcalSingleEvent(b *testing.B) {
	fileContent, err := os.ReadFile(singleEventFileName)
	if err != nil {
		panic("Invalid File")
	}

	var reader bytes.Reader
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
