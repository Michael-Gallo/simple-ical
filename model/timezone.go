package model

// a grouping of component properties that defines a time zone.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.5
type TimeZone struct {
	// Represented by TZID
	// The time zone identifier for the time zone used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.1
	TimeZoneId string

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.3
	TimeZoneOffsetFrom string

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.3.4
	TimeZoneOffsetTo string
}
