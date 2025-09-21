// Package benchmarks provides comparative benchmarks against other Go iCalendar parsers
package benchmarks

import (
	"bytes"
	_ "embed"
	"os"
	"testing"

	"github.com/apognu/gocal"
	golangical "github.com/arran4/golang-ical"
	"github.com/michael-gallo/simple-ical/parse"
)

const commonName = "Org"

func BenchmarkAll(b *testing.B) {
	// Pre-load file content to avoid I/O overhead in benchmarks
	fileContent, err := os.ReadFile("./test_event.ical")
	if err != nil {
		panic("Invalid File")
	}
	var reader bytes.Reader

	b.Run("Gocal", func(b *testing.B) {
		for b.Loop() {
			reader.Reset(fileContent)
			c := gocal.NewParser(&reader)
			err := c.Parse()
			if err != nil {
				panic(err)
			}
			if c.Events[0].Organizer.Cn != commonName {
				panic("Invalid Organizer")
			}
		}
	})

	b.Run("SimpleIcal", func(b *testing.B) {
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

	b.Run("GolangIcal", func(b *testing.B) {
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
