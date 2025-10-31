package icaldur

import (
	"errors"
	"strconv"
	"time"
)

var (
	ErrInvalidTimeFormat = errors.New("invalid iCal time format")
	ErrInvalidTimeValue  = errors.New("invalid time value")
)

// ParseIcalTime parses an iCal datetime string.
// Supports both UTC format (YYYYMMDDTHHMMSSZ) and floating time format (YYYYMMDDTHHMMSS).
// This manual implementation is faster than time.Parse for the fixed iCal format.
func ParseIcalTime(value string) (time.Time, error) {
	length := len(value)
	if length != 15 && length != 16 {
		return time.Time{}, ErrInvalidTimeFormat
	}

	if length == 16 {
		if value[15] != 'Z' {
			return time.Time{}, ErrInvalidTimeFormat
		}
	}

	// Parse year (positions 0-3)
	year, err := strconv.Atoi(value[0:4])
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}

	// Parse month (positions 4-5)
	month, err := strconv.Atoi(value[4:6])
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}
	if month < 1 || month > 12 {
		return time.Time{}, ErrInvalidTimeValue
	}

	// Parse day (positions 6-7)
	day, err := strconv.Atoi(value[6:8])
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}
	if day < 1 || day > 31 {
		return time.Time{}, ErrInvalidTimeValue
	}

	// Check for 'T' separator (position 8)
	if value[8] != 'T' {
		return time.Time{}, ErrInvalidTimeFormat
	}

	// Parse hour (positions 9-10)
	hour, err := strconv.Atoi(value[9:11])
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}
	if hour < 0 || hour > 23 {
		return time.Time{}, ErrInvalidTimeValue
	}

	// Parse minute (positions 11-12)
	minute, err := strconv.Atoi(value[11:13])
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}
	if minute < 0 || minute > 59 {
		return time.Time{}, ErrInvalidTimeValue
	}

	// Parse second (positions 13-14)
	second, err := strconv.Atoi(value[13:15])
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}
	if second < 0 || second > 59 {
		return time.Time{}, ErrInvalidTimeValue
	}

	// time.Date will validate the date (e.g., invalid day for month)
	// All times are returned in UTC (floating times are treated as UTC per iCal spec)
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
}
