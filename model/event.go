package model

import "time"

// An Event in the iCalendar format
// for more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
type Event struct {
	Summary     string
	Description string
	// dtstart
	Start     time.Time
	End       time.Time
	Location  string
	Organizer *Organizer
}

// An Organizer in the iCalendar format
// for more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
type Organizer struct {
	// denoted by CN= in the spec
	CommonName string
	// Note: Any Valid URI
	// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.3
	CalAddress CalendarAddress
	// denoted by DIR= in the spec
	Directory string
}

type CalendarAddress struct {
	URI      string
	IsMailTo bool
}
