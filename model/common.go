// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"
)

// Organizer represents an ORGANIZER component in the iCalendar format, used in VEVENT, VTODO, and VJOURNAL
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

// BaseComponent represents common fields found in all top level calendar components.
type BaseComponent struct {

	// a DTSTAMP property defines the date and time that the instance of the calendar component was created.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.2
	// Note: This is technically mandatory in the spec, however I have seen examples in the wild where it is not present.
	// I will not be enforcing this requirement in the parser. I may at some point in the future add a strict mode.
	DTStamp time.Time

	// The unique identifier for the event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.7
	UID string
}

// Contact is used to represent contact information
// Can be specified in Events, Todos, Journals, and FreeBusy Components
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.2
type Contact = string

// Sequence is used to define the revision sequence number of the component
// Can be specified in Events, Todos, and Journals
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.4
type Sequence = int
