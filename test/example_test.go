package test

import (
	"fmt"
	"testing"

	ical "github.com/michael-gallo/simple-ical/parse"
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

func TestReadmeExample(t *testing.T) {
	calendar, err := ical.IcalString(testIcalString)
	if err != nil {
		t.Fatalf("Failed to parse iCal string: %v", err)
	}

	fmt.Println(calendar.Description)
}
