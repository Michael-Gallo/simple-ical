package model

// Calendar represents a VCALENDAR component in the iCalendar format.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.4
type Calendar struct {
	// Specifies the identifier corresponding to the
	// highest version number or the minimum and maximum range of the
	// iCalendar specification that is required in order to interpret the
	// iCalendar object.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.4
	Version string
	// Product Identifier.
	// This property specifies the identifier for the product that
	// created the iCalendar object.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.3
	ProdId string
	// CalScale specifies the calendar scale used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.1
	CalScale string
	// Method specifies the method used by the calendar component.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.2
	Method    string
	TimeZones []TimeZone

	// A grouping of component properties that describe an event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
	Events []Event

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.2
	Todos []Todo
}
