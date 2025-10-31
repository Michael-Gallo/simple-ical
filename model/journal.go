// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"

	"github.com/michael-gallo/simpleical/rrule"
)

// JournalStatus represents the possible values for a VJOURNAL's STATUS field.
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
type JournalStatus string

const (
	JournalStatusDraft     JournalStatus = "DRAFT"
	JournalStatusFinal     JournalStatus = "FINAL"
	JournalStatusCancelled JournalStatus = "CANCELLED"
)

// JournalClass represents the possible values for a VJOURNAL's CLASS field.
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.3
type JournalClass string

const (
	JournalClassPublic       JournalClass = "PUBLIC"
	JournalClassPrivate      JournalClass = "PRIVATE"
	JournalClassConfidential JournalClass = "CONFIDENTIAL"
)

// Journal represents a VJOURNAL component in the iCalendar format.
// A VJOURNAL is a grouping of component properties that describe a journal entry.
// Does not take up time on a calendar.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.3
type Journal struct {
	// REQUIRED, MUST NOT occur more than once
	// a DTSTAMP property defines the date and time that the instance of the calendar component was created.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.2
	DTStamp time.Time

	// REQUIRED, MUST NOT occur more than once
	// The unique identifier for the journal entry.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.7
	UID string

	// OPTIONAL, MUST NOT occur more than once
	// Access Classification for the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.3
	Class JournalClass

	// OPTIONAL, MUST NOT occur more than once
	// Specifies the date and time that the calendar information was created.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.1
	Created time.Time

	// OPTIONAL, MUST NOT occur more than once
	// Specifies when the calendar component begins.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.4
	DTStart time.Time

	// OPTIONAL, MUST NOT occur more than once
	// Specifies the date and time that the information associated with the calendar component was last revised.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.3
	LastModified time.Time

	// OPTIONAL, MUST NOT occur more than once
	// The organizer of the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
	Organizer *Organizer

	// OPTIONAL, MUST NOT occur more than once
	// Specifies the revision sequence number of the calendar component within a sequence of revisions.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.4
	RecurrenceID time.Time

	// OPTIONAL, MUST NOT occur more than once
	// Specifies the revision sequence number of the calendar component within a sequence of revisions.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.4
	Sequence int

	// OPTIONAL, MUST NOT occur more than once
	// Defines the overall status or confirmation for the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.11
	Status JournalStatus

	// OPTIONAL, MUST NOT occur more than once
	// A short, one-line summary about the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.12
	Summary string

	// OPTIONAL, MUST NOT occur more than once
	// Specifies a URL associated with the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.6
	URL string

	// OPTIONAL, SHOULD NOT occur more than once
	// TODO: RRULE - define once per journal
	RRule *rrule.RRule

	// OPTIONAL, MAY occur more than once
	// Provides the capability to associate a document object with a calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.1
	Attach []string

	// OPTIONAL, MAY occur more than once
	// Specifies the participants that are invited to the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.1
	Attendees []url.URL

	// OPTIONAL, MAY occur more than once
	// Specifies the categories that the calendar component belongs to.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.2
	Categories []string

	// OPTIONAL, MAY occur more than once
	// Specifies non-processing information intended to provide a comment to the calendar user.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.4
	Comment []string

	// OPTIONAL, MAY occur more than once
	// Specifies the contact information for the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.2
	Contacts []string

	// OPTIONAL, MAY occur more than once
	// Used to capture lengthy textual descriptions associated with the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.5
	Description []string

	// OPTIONAL, MAY occur more than once
	// Specifies the list of date/time exceptions for a recurring calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.1
	ExceptionDates []time.Time

	// OPTIONAL, MAY occur more than once
	// Specifies a relationship or reference between one calendar component and another.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.5
	Related []string

	// OPTIONAL, MAY occur more than once
	// Specifies the list of date/time values for recurring activities.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.2
	Rdate []time.Time

	// OPTIONAL, MAY occur more than once
	// Specifies the status code returned for a scheduling request.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.3
	RequestStatus []string

	// OPTIONAL, MAY occur more than once
	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2
	XProp map[string]string

	// OPTIONAL, MAY occur more than once
	// An IANA registered property name.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1
	IANAProp map[string]string

	// OPTIONAL, MAY occur more than once
	// Sub-components: VALARM
	Alarms []Alarm
}
