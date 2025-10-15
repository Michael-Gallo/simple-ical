package icaldur

import "time"

// iCalDateTimeFormat represents the standard iCal datetime format
// Format: YYYYMMDDTHHMMSSZ (e.g., 20250928T183000Z).
const iCalDateTimeFormat = "20060102T150405Z"

func ParseIcalTime(value string) (time.Time, error) {
	return time.Parse(iCalDateTimeFormat, value)
}
