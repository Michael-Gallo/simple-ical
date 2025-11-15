# Simple-ical

A very much not ready ICAL parser for Golang intended to follow the official [ICAL 2.0 spec](https://datatracker.ietf.org/doc/html/rfc5545) as closely as is reasonable.

Focused on ease of use and good documentation, with frequent links to the spec.

[![Go Reference](https://pkg.go.dev/badge/github.com/michael-gallo/simpleical.svg)](https://pkg.go.dev/github.com/michael-gallo/simpleical)

## Documentation

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/michael-gallo/simpleical).

## Deviations from spec

1. The VCALENDAR spec does not address whitespace at the end of lines. We assume in this parser it is to be ignored and right trim all whitespace.
2. The `DTSTAMP` property is [mandatory](https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1), however, I have seen real life examples where it is not filled out. Ergo I will not be enforcing it here. If I do enforce it in the future, it will be in an opt-in strict mode.

## License

This project is licensed under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.


## Installation


```sh
go get github.com/michael-gallo/simpleical
```


## Usage

```go
package main
import (
    ical "github.com/michael-gallo/simpleical/parse"
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

## Performance
Performance tests were ran against [golang-ical v0.3.2](https://github.com/arran4/golang-ical/releases/tag/v0.3.2) and [gocal v0.9.1](https://github.com/apognu/gocal/releases/tag/v0.9.1)

### Specs
All tests were ran on a 5700X3D Processor with 32GB of RAM.

### Single Event Calendar File

|         | Gocal       | SimpleIcal  | GolangIcal |
|---------|-------------|-------------| -----------|
| sec/op  | 11.32µ ± 1% | 4.760µ ± 1% | 26.83µ ± 0%|
| B/op    | 11.50Ki ± 0%| 6.320Ki ± 0%|17.89Ki ± 0%|
|allocs/op| 198.0 ± 0%  | 53.00 ± 0%  |439.0 ± 0%  |


### Multiple Event Calendar File

|         | Gocal       | SimpleIcal  | GolangIcal  |
|---------|-------------|-------------| ------------|
| sec/op  | 18.06µ ± 1% | 7.585µ ± 1% | 42.74µ ± 0% |
| B/op    | 16.79Ki ± 0%| 8.391Ki ± 0%| 27.09Ki ± 0%|
|allocs/op| 314.0 ± 0%  | 87.00 ± 0%  | 692.0 ± 0%  |


### Complex Calendar File

|         | Gocal       | SimpleIcal  | GolangIcal  |
|---------|-------------|-------------| ------------|
| sec/op  | 21.24µ ± 1% | 9.832µ ± 1% | 58.01µ ± 1% |
| B/op    | 18.78Ki ± 0%| 10.27Ki ± 0%| 32.35Ki ± 0%|
|allocs/op|412.0 ± 0%   | 114.0 ± 0%  |959.0 ± 0%   |
