# Simple-ical

A very much not ready ICAL parser for Golang intended to follow the official [ICAL 2.0 spec](https://datatracker.ietf.org/doc/html/rfc5545) as closely as is reasonable.

Focused on ease of use and good documentation, with frequent links to the spec.

## Deviations from spec

1. The VCALENDAR spec does not address whitespace or blank lines. We assume in this parser that they are to be ignored and trim all whitespace from every line.
2. The `DTSTAMP` property is [mandatory](https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1), however, I have seen real life examples where it is not filled out. Ergo I will not be enforcing it here. If I do enforce it in the future, it will be in an opt-in strict mode.

## License

This project is licensed under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.


## Installation


```sh
go get github.com/michael-gallo/simple-ical
```


## Usage

```go
package main
import (
    ical "github.com/michael-gallo/simple-ical/parse"
    "fmt"
    )


const testIcalString string =  `BEGIN:VCALENDAR
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

func main(){
    calendar,err := ical.ParseIcalString(testIcalString)
    if err != nil {
        panic("Broken calendar string")
    }
    fmt.Println(calendar.Description)

}

```
