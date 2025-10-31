package icaldur

import (
	"strings"
	"time"
)

// iCalDateTimeFormat represents the standard iCal datetime format
// Format: YYYYMMDDTHHMMSSZ (e.g., 20250928T183000Z).
const iCalDateTimeFormat = "20060102T150405Z"
const iCalDateTimeFormatNoZ = "20060102T150405"

func ParseIcalTime(value string) (time.Time, error) {
	if strings.HasSuffix(value, "Z") {
		return time.Parse(iCalDateTimeFormat, value)
	}

	return time.Parse(iCalDateTimeFormatNoZ, value)
}
