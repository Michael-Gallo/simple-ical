// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parse

import "errors"

// Calendar-level errors.
var (
	errNoCalendarFound                   = errors.New("empty calendar sent")
	errInvalidCalendarFormatMissingBegin = errors.New("invalid calendar format: must start with BEGIN:VCALENDAR")
	errInvalidCalendarFormatMissingEnd   = errors.New("invalid calendar format: must end with END:VCALENDAR")
	errContentAfterEndBlock              = errors.New("content after END:VCALENDAR")
	errTemplateInvalidEndBlock           = errors.New("invalid end block")
	errTemplateInvalidStartBlock         = errors.New("invalid start block")

	// General parsing errors.
	errInvalidPropertyLine = errors.New("invalid property line in iCal data")
	errDuplicateProperty   = errors.New("duplicate property")

	// URI parsing errors.
	// errInvalidProtocol is one of the errors that could be returned when parsing a URI with the standard library.
	errInvalidProtocol = errors.New("parse \"://invalid\": missing protocol scheme")
)

// Event-specific errors.
var (
	errInvalidEventPropertyLineMissingColon = errors.New("invalid event property line missing colon in line")
	errInvalidEventProperty                 = errors.New("invalid event property")
	errInvalidEventPropertySequence         = errors.New("invalid event property in iCal Event: SEQUENCE must be an integer")
	errLineShouldStartWithOrganizer         = errors.New("line should start with ORGANIZER")

	// Event date/time property errors.
	errInvalidDatePropertyDtstart      = errors.New("invalid date property in iCal Event: DTSTART")
	errInvalidDatePropertyDtend        = errors.New("invalid date property in iCal Event: DTEND")
	errInvalidDatePropertyDTStamp      = errors.New("invalid date property in iCal Event: DTSTAMP")
	errInvalidDatePropertyLastModified = errors.New("invalid date property in iCal Event: LAST-MODIFIED")

	// Event duration property errors.
	errInvalidDurationProperty      = errors.New("invalid duration property in iCal Event: DURATION")
	errInvalidDurationPropertyDtend = errors.New("invalid duration property in iCal Event: DTEND and DURATION are mutually exclusive")

	// Event geographic property errors.
	errInvalidGeoProperty          = errors.New("invalid event property in iCal Event: GEO must be two floats separated by a semicolon")
	errInvalidGeoPropertyLatitude  = errors.New("invalid latitude in iCal Event: GEO must be a float")
	errInvalidGeoPropertyLongitude = errors.New("invalid longitude in iCal Event: GEO must be a float")
)
