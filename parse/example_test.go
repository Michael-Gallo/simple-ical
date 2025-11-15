package parse_test

import (
	"fmt"
	"strings"

	"github.com/michael-gallo/simpleical/parse"
)

const testIcalString string = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Event//Event Calendar//EN
CALSCALE:GREGORIAN
METHOD:REQUEST
BEGIN:VTIMEZONE
TZID:America/Detroit
BEGIN:STANDARD
DTSTART:19700101T000000
TZOFFSETFROM:+0000
TZOFFSETTO:+0000
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
UID:13235@example.com
DTSTART:20250928T183000Z
DTEND:20250928T203000Z
SUMMARY:Event Summary
DESCRIPTION:Event Description
LOCATION:555 Fake Street
ORGANIZER;CN=Org:MAILTO:hello@world
STATUS:CONFIRMED
SEQUENCE:0
TRANSP:OPAQUE
END:VEVENT
END:VCALENDAR
`

func ExampleIcalString() {
	calendar, err := parse.IcalString(testIcalString)
	if err != nil {
		panic(err)
	}

	fmt.Println(calendar.ProdID)
	fmt.Println(calendar.TimeZones[0].TimeZoneID)
	fmt.Println(calendar.Events[0].Summary)
	// Output:
	// -//Event//Event Calendar//EN
	// America/Detroit
	// Event Summary
}

func ExampleIcalFromFileName() {
	calendar, err := parse.IcalFromFileName("../test/test_data/calendar/valid_calendar.ical")
	if err != nil {
		panic(err)
	}

	fmt.Println(calendar.ProdID)
	fmt.Println(calendar.CalScale)
	// Output:
	// -//Event//Event Calendar//EN
	// GREGORIAN
}

func ExampleIcalReader() {
	reader := strings.NewReader(testIcalString)
	calendar, err := parse.IcalReader(reader)
	if err != nil {
		panic(err)
	}

	fmt.Println(calendar.ProdID)
	fmt.Println(calendar.TimeZones[0].TimeZoneID)
	fmt.Println(calendar.Events[0].Summary)
	// Output:
	// -//Event//Event Calendar//EN
	// America/Detroit
	// Event Summary
}
