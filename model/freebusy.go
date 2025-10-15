// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"
)

// FreeBusyStatus represents the possible values for a VFREEBUSY's FREEBUSY property.
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.6
type FreeBusyStatus string

const (
	FreeBusyStatusFree            FreeBusyStatus = "FREE"
	FreeBusyStatusBusy            FreeBusyStatus = "BUSY"
	FreeBusyStatusBusyTentative   FreeBusyStatus = "BUSY-TENTATIVE"
	FreeBusyStatusBusyUnavailable FreeBusyStatus = "BUSY-UNAVAILABLE"
)

// FreeBusy represents a VFREEBUSY component in the iCalendar format.
// A VFREEBUSY is a grouping of component properties that describe either a request for free/busy time,
// describe a response to a request for free/busy time, or describe a published set of busy time.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.4
type FreeBusy struct {
	// REQUIRED, MUST NOT occur more than once
	// a DTSTAMP property defines the date and time that the instance of the calendar component was created.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.2
	DTStamp time.Time

	// REQUIRED, MUST NOT occur more than once
	// The unique identifier for the free/busy component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.7
	UID string

	// OPTIONAL, MUST NOT occur more than once
	// Specifies the contact information for the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.2
	Contact string

	// OPTIONAL, MUST NOT occur more than once
	// Specifies when the calendar component begins.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.4
	DTStart time.Time

	// OPTIONAL, MUST NOT occur more than once
	// Specifies when the calendar component ends.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.2
	DTEnd time.Time

	// OPTIONAL, MUST NOT occur more than once
	// The organizer of the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
	Organizer *Organizer

	// OPTIONAL, MUST NOT occur more than once
	// Specifies a URL associated with the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.6
	URL string

	// OPTIONAL, MAY occur more than once
	// Specifies the participants that are invited to the activity.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.1
	Attendees []url.URL

	// OPTIONAL, MAY occur more than once
	// Specifies non-processing information intended to provide a comment to the calendar user.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.4
	Comment []string

	// OPTIONAL, MAY occur more than once
	// Specifies one or more free or busy time intervals.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.6
	FreeBusy []FreeBusyTime

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
}

// FreeBusyTime represents a single free/busy time interval with its status.
type FreeBusyTime struct {
	// The start time of the free/busy interval.
	Start time.Time
	// The end time of the free/busy interval.
	End time.Time
	// The status of the time interval (FREE, BUSY, BUSY-TENTATIVE, BUSY-UNAVAILABLE).
	Status FreeBusyStatus
}
