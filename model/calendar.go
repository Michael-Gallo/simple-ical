// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package model contains structs used throughout the project
package model

// Calendar represents a VCALENDAR component in the iCalendar format.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.4
// Documentation on the properties can be found here:
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7
type Calendar struct {
	// Specifies the identifier corresponding to the
	// highest version number or the minimum and maximum range of the
	// iCalendar specification that is required in order to interpret the
	// iCalendar object. This property is required.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.4
	Version string
	// Product Identifier.
	// This property specifies the identifier for the product that
	// created the iCalendar object.
	// This property is required.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.3
	ProdID string
	// CalScale specifies the calendar scale used by the calendar component.
	// This property is optional
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.1
	CalScale string
	// Method specifies the method used by the calendar component.
	// This property is optional.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.7.2
	Method string

	TimeZones []TimeZone

	// A grouping of component properties that describe an event.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
	Events []Event

	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.2
	Todos []Todo
}
