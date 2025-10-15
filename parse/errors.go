// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parse

import "errors"

// Calendar-level errors.
var (
	ErrNoCalendarFound                   = errors.New("empty calendar sent")
	ErrInvalidCalendarFormatMissingBegin = errors.New("invalid calendar format: must start with BEGIN:VCALENDAR")
	ErrInvalidCalendarFormatMissingEnd   = errors.New("invalid calendar format: must end with END:VCALENDAR")
	ErrInvalidCalendarEmptyLine          = errors.New("invalid calendar format: must not contain empty lines")
	ErrContentAfterEndBlock              = errors.New("content after END:VCALENDAR")
	ErrTemplateInvalidEndBlock           = errors.New("invalid end block")
	ErrTemplateInvalidStartBlock         = errors.New("invalid start block")
	ErrMissingCalendarVersionProperty    = errors.New("calendar must have a VERSION property")
	ErrMissingCalendarProdIDProperty     = errors.New("calendar must have a PRODID property")

	// General parsing errors.
	ErrInvalidPropertyLine = errors.New("invalid property line in iCal data")
	ErrDuplicateProperty   = errors.New("duplicate property")

	// URI parsing errors.
	// ErrInvalidProtocol is one of the errors that could be returned when parsing a URI with the standard library.
	ErrInvalidProtocol = errors.New("parse \"://invalid\": missing protocol scheme")
)

// Event-specific errors.
var (
	ErrInvalidEventProperty = errors.New("invalid event property")

	ErrMissingEventUIDProperty     = errors.New("event must have a UID property")
	ErrMissingEventDTStartProperty = errors.New("event must have a DTSTART property if no METHOD property is present for the top level calendar")

	// Event duration property errors.
	ErrInvalidDurationPropertyDtend = errors.New("invalid duration property in iCal Event: DTEND and DURATION are mutually exclusive")

	// Event geographic property errors.
	ErrInvalidGeoProperty          = errors.New("invalid event property in iCal Event: GEO must be two floats separated by a semicolon")
	ErrInvalidGeoPropertyLatitude  = errors.New("invalid latitude in iCal Event: GEO must be a float")
	ErrInvalidGeoPropertyLongitude = errors.New("invalid longitude in iCal Event: GEO must be a float")
)

// Todo-specific errors.
var (
	ErrInvalidTodoProperty = errors.New("invalid todo property")

	ErrMissingTodoUIDProperty = errors.New("todo must have a UID property")

	ErrMissingTodoDTStartProperty = errors.New("todo must have a DTSTART property")

	// Todo duration property errors.
	ErrInvalidDurationPropertyDue = errors.New("invalid duration property in iCal Todo: DUE and DURATION are mutually exclusive")
)

// Journal-specific errors.
var (
	ErrInvalidJournalProperty = errors.New("invalid journal property")

	ErrMissingJournalUIDProperty = errors.New("journal must have a UID property")

	ErrMissingJournalDTStartProperty = errors.New("journal must have a DTSTART property")
)

// FreeBusy-specific errors.
var (
	ErrInvalidFreeBusyProperty = errors.New("invalid freebusy property")

	ErrMissingFreeBusyUIDProperty = errors.New("freebusy must have a UID property")

	ErrInvalidFreeBusyFormat = errors.New("invalid FREEBUSY property format")

	ErrMissingFreeBusyDTStartProperty = errors.New("freebusy must have a DTSTART property")
)

// Timezone-specific errors.
var (
	ErrInvalidTimezoneProperty       = errors.New("invalid timezone property")
	ErrMissingTimezoneTZIDProperty   = errors.New("timezone must have a TZID property")
	ErrInvalidTimezoneDatetimeFormat = errors.New("invalid timezone datetime format")
)

// Alarm-specific errors.
var (
	ErrInvalidAlarmProperty = errors.New("invalid alarm property")

	ErrMissingAlarmActionProperty = errors.New("alarm must have an ACTION property")

	ErrMissingAlarmTriggerProperty = errors.New("alarm must have a TRIGGER property")

	ErrMissingAlarmDescriptionForDisplay = errors.New("DISPLAY alarm must have a DESCRIPTION property")

	ErrMissingAlarmDescriptionForEmail = errors.New("EMAIL alarm must have a DESCRIPTION property")

	ErrMissingAlarmSummaryForEmail = errors.New("EMAIL alarm must have a SUMMARY property")

	ErrMissingAlarmAttendeesForEmail = errors.New("EMAIL alarm must have at least one ATTENDEE property")
)

// Property Setter errors.

const ErrDuplicatePropertyInComponentFormat = "%w: %s set twice in component %s"

var (
	ErrDuplicatePropertyInComponent = errors.New("duplicate property error")
	ErrParseErrorInComponent        = errors.New("parse error in component")
)
