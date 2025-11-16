// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"

	"github.com/michael-gallo/simpleical/rrule"
)

// EventStatus represents VEVENT STATUS values. Note VTODO STATUS values are different.
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11.
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "CONFIRMED"
	EventStatusTentative EventStatus = "TENTATIVE"
	EventStatusCancelled EventStatus = "CANCELED"
)

// EventTransp represents VEVENT TRANSP values. Note VTODO TRANSP values are different.
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.7.
type EventTransp string

const (
	EventTranspTransparent EventTransp = "TRANSPARENT"
	EventTranspOpaque      EventTransp = "OPAQUE"
)

// Event represents a VEVENT component in the iCalendar format.
// For more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1.
type Event struct {

	// DTStamp defines the date and time that the event was created.
	// Note: This is mandatory in RFC5545, but that is not enforced in this parser.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.2
	DTStamp time.Time

	// UID is the unique identifier for the event.
	// REQUIRED, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.7
	UID string

	// Start defines the date and time that the event begins. Refers to the DTSTART property.
	// REQUIRED if no METHOD property. MUST NOT occur more than once
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.4
	Start time.Time

	// Summary is a short, one-line summary about the event. Refers to the SUMMARY property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.12
	Summary string

	// Description is used to capture lengthy textual descriptions associated with the event. Refers to the DESCRIPTION property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.5
	Description string

	// Geo specifies the latitude and longitude of the activity specified by a calendar component.
	// Refers to the GEO property. Can be specified in Events and Todos.
	// Must be precise up to 6 decimal places.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.6.
	Geo []float64

	// LastModified specifies the date and time tthat the information associated with the calendar information was last revised.
	// Refers to the LAST-MODIFIED property. Can be specified in Events, Todos, Journals, and TimeZones.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.3.
	LastModified time.Time

	// Location is the location where the event takes place. Refers to the LOCATION property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.7.
	Location string

	// Organizer is the organizer of the event. Refers to the ORGANIZER property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3.
	Organizer *Organizer

	// Priority represents the priority of the event (0-9, where 0 is undefined, 1 is highest, 9 is lowest)
	// Refers to the PRIORITY property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.9.
	Priority int

	// Sequence is used to define the revision sequence number of the component
	// Refers to the SEQUENCE property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.4.
	Sequence int

	// Status defines the overall status or confirmation for the event. Refers to the STATUS property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11.
	Status EventStatus

	// Transp is the time transparency of the event. Refers to the TRANSP property.
	// Time transparency refers to whether the event is considered to consume time on the calendar.
	// ie: If an event is TRANSPARENT, participants are not to be considered busy during the event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.7.
	Transp EventTransp

	// URL specifies a URL associated with the event. Refers to the URL property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.6.
	URL string

	// RecurrenceID is the recurrence identifier for the event. Refers to the RECURRENCE-ID property.
	// OPTIONAL, MUST NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.4.
	RecurrenceID time.Time

	// RRule is the recurrence rule for the event. Refers to the RRULE property.
	// OPTIONAL, SHOULD NOT occur more than once.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.3
	RRule *rrule.RRule

	// Either dtend or duration (not both).
	// dtend in the ICAL format.
	// Can not be specified if a Duration is specified.
	// See the datetime specification for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5.
	End time.Time

	// The event's duration.
	// Can not be specified if an End time is specified.
	// See the Duration specficiation for more information: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.6.
	Duration time.Duration

	// OPTIONAL, MAY occur more than once.
	// Optional and can be defined multiple times.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.1.
	// TODO: define a more robust type for attachments.
	Attach []string

	// Attendee is used to represent an ATTENDEE component in the iCalendar format.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.1.
	Attendees []url.URL

	// Categories specifies the categories that the calendar component belongs to.
	// Can be specified in Events, Todos, and Journals.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.2.
	Categories []string

	// Comment specifies non-processing information intended to provide a comment to the calendar user.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.4.
	Comment []string

	// Contact is used to represent contact information.
	// Can be specified in Events, Todos, Journals, and FreeBusy Components.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.2.
	Contacts []string

	// Exception Date-Times, property name EXDATE.
	// This is optional and repeatable.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.1.
	ExceptionDates []time.Time

	// Property Name: REQUEST-STATUS Represented as RSTATUS.
	// The status code returned for a scheduling request.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.3.
	RequestStatus []string

	// Property Name: RELATED-TO.
	// Used to represent a relationship or reference between one calendar component and another.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.5.
	Related []string

	// Property Name: RESOURCES.
	// Defines equipment or resources anticipated for an event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.10.
	Resources []string

	// Recurrence Date-Times.
	// This is optional and repeatable.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.2.
	Rdate []time.Time

	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// This is optional and repeatable.
	// The keys of the map are expected to include the X-prefix.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2.
	XProp map[string]string

	// An IANA registered property name.
	// This is optional and repeatable.
	// As of right now this is implemented as a map of string to string with no validation of whether the property is a real IANA registered property.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1.
	IANAProp map[string]string

	// OPTIONAL, MAY occur more than once.
	// Sub-components: VALARM.
	Alarms []Alarm
}
