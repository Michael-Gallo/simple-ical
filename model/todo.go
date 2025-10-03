// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import "time"

// Todo represents a VTODO component in the iCalendar format.
// A VTODO is a grouping of component properties that describe a to-do,
// appointment, or journal entry.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.2
type Todo struct {
	// TODO: Add fields for summary, description, due date, status, etc.
	// This struct will be expanded to include all VTODO properties
	// as defined in RFC 5545 section 3.6.2

	// Comment specifies non-processing information intended to provide a comment to the calendar user.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.4
	Comment []string

	// The unique identifier for the event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.7
	UID string

	// a DTSTAMP property defines the date and time that the instance of the calendar component was created.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.2
	// Note: This is technically mandatory in the spec, however I have seen examples in the wild where it is not present.
	// I will not be enforcing this requirement in the parser. I may at some point in the future add a strict mode.
	DTStamp time.Time

	Due      time.Time
	Duration time.Duration

	// TODO: RRULE?
}
