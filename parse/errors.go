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
	errInvalidCalendarEmptyLine          = errors.New("invalid calendar format: must not contain empty lines")
	errContentAfterEndBlock              = errors.New("content after END:VCALENDAR")
	errTemplateInvalidEndBlock           = errors.New("invalid end block")
	errTemplateInvalidStartBlock         = errors.New("invalid start block")
	errMissingCalendarVersionProperty    = errors.New("calendar must have a VERSION property")
	errMissingCalendarProdIDProperty     = errors.New("calendar must have a PRODID property")

	// General parsing errors.
	errInvalidPropertyLine = errors.New("invalid property line in iCal data")
	errDuplicateProperty   = errors.New("duplicate property")
)

// Event-specific errors.
var (
	errInvalidEventProperty = errors.New("invalid event property")

	errMissingEventUIDProperty     = errors.New("event must have a UID property")
	errMissingEventDTStartProperty = errors.New("event must have a DTSTART property if no METHOD property is present for the top level calendar")

	// Event duration property errors.
	errInvalidDurationPropertyDtend = errors.New("invalid duration property in iCal Event: DTEND and DURATION are mutually exclusive")

	// Event geographic property errors.
	errInvalidGeoProperty          = errors.New("invalid event property in iCal Event: GEO must be two floats separated by a semicolon")
	errInvalidGeoPropertyLatitude  = errors.New("invalid latitude in iCal Event: GEO must be a float")
	errInvalidGeoPropertyLongitude = errors.New("invalid longitude in iCal Event: GEO must be a float")
)

// Todo-specific errors.
var (
	errInvalidTodoProperty = errors.New("invalid todo property")

	errMissingTodoUIDProperty = errors.New("todo must have a UID property")

	errMissingTodoDTStartProperty = errors.New("todo must have a DTSTART property")

	// Todo duration property errors.
	errInvalidDurationPropertyDue = errors.New("invalid duration property in iCal Todo: DUE and DURATION are mutually exclusive")
)

// Journal-specific errors.
var (
	errInvalidJournalProperty = errors.New("invalid journal property")

	errMissingJournalUIDProperty = errors.New("journal must have a UID property")

	errMissingJournalDTStartProperty = errors.New("journal must have a DTSTART property")
)

// FreeBusy-specific errors.
var (
	errInvalidFreeBusyProperty = errors.New("invalid freebusy property")

	errMissingFreeBusyUIDProperty = errors.New("freebusy must have a UID property")

	errInvalidFreeBusyFormat = errors.New("invalid FREEBUSY property format")

	errMissingFreeBusyDTStartProperty = errors.New("freebusy must have a DTSTART property")
)

// Timezone-specific errors.
var (
	errInvalidTimezoneProperty     = errors.New("invalid timezone property")
	errMissingTimezoneTZIDProperty = errors.New("timezone must have a TZID property")
)

// Alarm-specific errors.
var (
	errInvalidAlarmProperty = errors.New("invalid alarm property")

	errMissingAlarmActionProperty = errors.New("alarm must have an ACTION property")

	errMissingAlarmTriggerProperty = errors.New("alarm must have a TRIGGER property")

	errMissingAlarmDescriptionForDisplay = errors.New("DISPLAY alarm must have a DESCRIPTION property")

	errMissingAlarmDescriptionForEmail = errors.New("EMAIL alarm must have a DESCRIPTION property")

	errMissingAlarmSummaryForEmail = errors.New("EMAIL alarm must have a SUMMARY property")

	errMissingAlarmAttendeesForEmail = errors.New("EMAIL alarm must have at least one ATTENDEE property")
)

// Property Setter errors.

const errDuplicatePropertyInComponentFormat = "%w: %s set twice in component %s"

var (
	errDuplicatePropertyInComponent = errors.New("duplicate property error")
	errParseErrorInComponent        = errors.New("parse error in component")
)
