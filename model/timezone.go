// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

// TimeZone represents a VTIMEZONE component in the iCalendar format.
// A grouping of component properties that defines a time zone.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.5
type TimeZone struct {
	// Represented by TZID
	// The time zone identifier for the time zone used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.1
	TimeZoneID string

	// The time zone offset from UTC when daylight saving time is in effect.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.3
	TimeZoneOffsetFrom string

	// The time zone offset from UTC when standard time is in effect.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.4
	TimeZoneOffsetTo string
}
