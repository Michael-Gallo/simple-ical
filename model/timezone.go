// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"
)

// TimeZone represents a VTIMEZONE component in the iCalendar format.
// A grouping of component properties that defines a time zone.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.5
type TimeZone struct {
	// Represented by TZID
	// The time zone identifier for the time zone used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.1
	TimeZoneID string

	// The last modification time of the time zone.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.5
	LastMod time.Time

	// Time Zone URL, represented as tzurl in the spec, can be any valid URI
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.5
	TimeZoneURL *url.URL

	Standard []TimeZoneProperty
	Daylight []TimeZoneProperty
}

// TimeZoneProperty is defined in the spec as tzprop and describes the fields that are used to represent either a standard or daylight sub-component in a timezone.
type TimeZoneProperty struct {
	// The time zone offset from UTC when daylight saving time is in effect.
	// This property is required
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.3
	TimeZoneOffsetFrom string

	// The time zone offset from UTC when standard time is in effect.
	// This property is required
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.4
	TimeZoneOffsetTo string

	// Date-Time Start, used to specify when the calendar event starts
	// This property is required
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.2.4
	DTStart time.Time

	// A comment to describe the Time Zone Property
	// This property is optional.
	Comment []string

	// Recurrence Date-Times
	// This property is optional.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.2
	Rdate []string

	// Represented by tzname
	// This property is optional.
	TimeZoneName []string

	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// This property is optional.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2
	XProp []string

	// An IANA registered property name.
	// This property is optional and repeatable.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1
	IanaProp []string

	// TODO: Implement RRule
	// RRule *rrule.RRule
}
