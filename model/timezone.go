// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"

	"github.com/michael-gallo/simpleical/rrule"
)

// TimeZone represents a VTIMEZONE component in the iCalendar format.
// A grouping of component properties that defines a time zone.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.5
type TimeZone struct {
	// REQUIRED, MUST NOT occur more than once
	// Represented by TZID
	// The time zone identifier for the time zone used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.1
	TimeZoneID string

	// OPTIONAL, MUST NOT occur more than once
	// The last modification time of the time zone.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.7.3
	LastMod time.Time

	// OPTIONAL, MUST NOT occur more than once
	// Time Zone URL, represented as tzurl in the spec, can be any valid URI
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.5
	TimeZoneURL *url.URL

	// OPTIONAL, MAY occur more than once
	// Either Standard or Daylight must be present.
	// STANDARD sub-components define the time zone standard time rules
	Standard []TimeZoneProperty

	// OPTIONAL, MAY occur more than once
	// Either Standard or Daylight must be present.
	// DAYLIGHT sub-components define the time zone daylight saving time rules
	Daylight []TimeZoneProperty

	// OPTIONAL, MAY occur more than once
	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2
	XProp map[string]string

	// OPTIONAL, MAY occur more than once
	// An IANA registered property name.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1
	IANAProp map[string]string
}

// TimeZoneProperty is defined in the spec as tzprop and describes the fields that are used to represent either a standard or daylight sub-component in a timezone.
type TimeZoneProperty struct {
	// REQUIRED, MUST NOT occur more than once
	// The time zone offset from UTC when daylight saving time is in effect.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.3
	TimeZoneOffsetFrom string

	// REQUIRED, MUST NOT occur more than once
	// The time zone offset from UTC when standard time is in effect.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.4
	TimeZoneOffsetTo string

	// REQUIRED, MUST NOT occur more than once
	// Date-Time Start, used to specify when the calendar event starts
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.4
	DTStart time.Time

	// OPTIONAL, MAY occur more than once
	// A comment to describe the Time Zone Property
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.4
	Comment []string

	// OPTIONAL, MAY occur more than once
	// Recurrence Date-Times
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.2
	Rdate []time.Time

	// OPTIONAL, MAY occur more than once
	// Represented by tzname
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.2
	TimeZoneName []string

	// OPTIONAL, SHOULD NOT occur more than once
	RRule *rrule.RRule

	// OPTIONAL, MAY occur more than once
	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2
	XProp map[string]string

	// OPTIONAL, MAY occur more than once
	// An IANA registered property name.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1
	IANAProp map[string]string
}
