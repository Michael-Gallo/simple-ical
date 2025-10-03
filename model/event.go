// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
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

// EventTransp represents the possible values for a VEVENT's TRANSP field, note VTODO's TRANSP field accepts different values
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.7
type EventTransp string

const (
	EventTranspTransparent EventTransp = "TRANSPARENT"
	EventTranspOpaque      EventTransp = "OPAQUE"
)

// Event represents a VEVENT component in the iCalendar format.
// For more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
type Event struct {

	// a DTSTAMP property defines the date and time that the instance of the calendar component was created.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.2
	// Note: This is technically mandatory in the spec, however I have seen examples in the wild where it is not present.
	// I will not be enforcing this requirement in the parser. I may at some point in the future add a strict mode.
	DTStamp time.Time

	// The unique identifier for the event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.7
	// REQUIRED, MUST NOT occur more than once
	UID string

	// REQUIRED if no METHOD property, MUST NOT occur more than once
	// dtstart in the ICAL format
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5
	Start time.Time

	// OPTIONAL, MUST NOT occur more than once
	// A short, one-line summary about the activity or journal entry.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.12
	Summary string

	// Used to capture lengthy textual descriptions associated with the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.5
	Description string

	// Geo specifies the latitude and longitude of the activity specified by a calendar component
	// Can be specified in Events and Todos
	// Must be precise up to 6 decimal places
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.6
	Geo []float64

	// LastModified specifies the date and time tthat the information associated with the calendar information was last revised
	// Can be specified in Events, Todos, Journals, and TimeZones
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.3
	LastModified time.Time

	// The location where the event takes place.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.7
	Location string

	// The organizer of the event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
	Organizer *Organizer

	// Priority represents the priority of the event (0-9, where 0 is undefined, 1 is highest, 9 is lowest)
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.9
	Priority int

	// Sequence is used to define the revision sequence number of the component
	// Can be specified in Events, Todos, and Journals
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.4
	Sequence int

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
	// defines the overall status or confirmation for the calendar component.
	Status EventStatus

	// The Time Transparency of the event.
	// This refers to whether the event is considered to consume time on the calendar
	// ie: if an event is TRANSPARENT, that means that participants are not to be considered busy during the event
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.7
	Transp EventTransp

	// URL specifies a URL associated with the event
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.6
	URL string

	// RecurrenceID specifies the recurrence identifier for the event
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.4
	RecurrenceID time.Time

	// OPTIONAL, SHOULD NOT occur more than once

	// TODO: RRULE , define once per event, doc: https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
	// RRule *RecurrenceRule

	// Either dtend or duration (not both)
	// dtend in the ICAL format
	// Can not be specified if a Duration is specified
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5
	End time.Time

	// The event's duration
	// Can not be specified if an End time is specified
	// See the Duration specficiation for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.6
	Duration time.Duration

	// OPTIONAL, MAY occur more than once
	// Optional and can be defined multiple times
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.1
	// TODO: define a more robust type for attachments
	Attach []string

	// Attendee is used to represent an ATTENDEE component in the iCalendar format
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.1
	Attendees []url.URL

	// Categories specifies the categories that the calendar component belongs to
	// Can be specified in Events, Todos, and Journals
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.2
	Categories []string

	// Comment specifies non-processing information intended to provide a comment to the calendar user.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.4
	Comment []string

	// Contact is used to represent contact information
	// Can be specified in Events, Todos, Journals, and FreeBusy Components
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.2
	Contacts []string

	// Exception Date-Times, property name EXDATE
	// This is optional and repeatable.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.1
	ExceptionDates []time.Time

	// Property Name: REQUEST-STATUS Represented as RSTATUS
	// The status code returned for a scheduling request
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.3
	RequestStatus []string

	// Property Name: RELATED-TO
	// Used to represent a relationship or reference between one calendar component and another
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.5
	Related []string

	// Property Name: RESOURCES
	// Defines equipment or resources anticipated for an event
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.10
	Resources []string

	// Recurrence Date-Times
	// This is optional and repeatable.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.2
	Rdate []time.Time

	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// This is optional and repeatable.
	// The keys of the map are expected to include the X-prefix
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2
	XProp map[string]string

	// An IANA registered property name.
	// This is optional and repeatable.
	// As of right now this is implemented as a map of string to string with no validation of whether the property is a real IANA registered property.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1
	IANAProp map[string]string
}
