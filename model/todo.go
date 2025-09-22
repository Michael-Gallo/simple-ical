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

	BaseComponent

	Due      time.Time
	Duration time.Duration

	// TODO: RRULE?
}
