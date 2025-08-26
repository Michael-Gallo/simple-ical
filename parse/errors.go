package parse

import "errors"

var (
	ErrNoVEventFound     = errors.New("no VEVENT found in iCal data")
	ErrInvalidDateFormat = errors.New("invalid date format in iCal data")
	ErrMissingDTSTART    = errors.New("missing DTSTART in VEVENT")
	ErrMissingDTEND      = errors.New("missing DTEND in VEVENT")
	ErrInvalidTimeFormat = errors.New("invalid time format in iCal data")

	ErrLineShouldStartWithOrganizer = errors.New("line should start with ORGANIZER")
)
