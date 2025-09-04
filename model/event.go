package model

import (
	"time"
)

// EventStatus represents the possible values for a VEVENT's STATUS field, note VTODO's STATUS field accepts different values
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "CONFIRMED"
	EventStatusTentative EventStatus = "TENTATIVE"
	EventStatusCancelled EventStatus = "CANCELED"
)

// Event represents a VEVENT component in the iCalendar format.
// For more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
type Event struct {
	// A short, one-line summary about the activity or journal entry.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.12
	Summary string
	// Used to capture lengthy textual descriptions associated with the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.5
	Description string
	// dtstart in the ICAL format
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5
	Start time.Time
	// dtend in the ICAL format
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5
	End time.Time
	// The location where the event takes place.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.7
	Location string

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
	// defines the overall status or confirmation for the calendar component.
	Status EventStatus
	// The organizer of the event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
	Organizer *Organizer
}
