// Package model contains structs used throughout the project
package model

import (
	"net/url"
	"time"
)

// The possible values for a VEVENT's STATUS field, note VTODO's STATUS field accepts different values
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "CONFIRMED"
	EventStatusTentative EventStatus = "TENTATIVE"
	EventStatusCancelled EventStatus = "CANCELLED"
)

// An Event in the iCalendar format
// for more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
type Event struct {
	// a short, one-line summary about the activity or journal entry.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.12
	Summary string
	// Used tocapture lengthy textual descriptions associated with the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.5
	Description string
	// dtstart in the ICAL format
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5
	Start time.Time
	// dtend in the ICAL format
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5
	End time.Time
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.7
	Location string

	// Represented by TZID in the spec
	// The time zone identifier for the time zone used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.1
	TimeZoneId string

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.3
	TimeZoneOffsetFrom string

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.4
	TimeZoneOffsetTo string

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
	// defines the overall status or confirmation for the calendar component.
	Status    EventStatus
	Organizer *Organizer
}

// An Organizer in the iCalendar format
// for more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
type Organizer struct {
	// denoted by CN= in the spec
	CommonName string
	// Note: Any Valid URI
	// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.3
	CalAddress *url.URL
	// denoted by DIR= in the spec
	Directory string
}
