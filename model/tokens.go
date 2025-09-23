// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

// SectionToken represents the names of the top level components in a VCALENDAR
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6
type SectionToken string

const (
	SectionTokenVCalendar SectionToken = "VCALENDAR"
	SectionTokenVEvent    SectionToken = "VEVENT"
	SectionTokenVTodo     SectionToken = "VTODO"
	SectionTokenVJournal  SectionToken = "VJOURNAL"
	SectionTokenVTimezone SectionToken = "VTIMEZONE"
	SectionTokenVFreebusy SectionToken = "VFREEBUSY"
	SectionTokenVAlarm    SectionToken = "VALARM"
	SectionTokenVStandard SectionToken = "STANDARD"
)

// EventToken represents the names of the properties in a VEVENT
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.1
type EventToken string

const (
	EventTokenSummary      EventToken = "SUMMARY"
	EventTokenDescription  EventToken = "DESCRIPTION"
	EventTokenLocation     EventToken = "LOCATION"
	EventTokenOrganizer    EventToken = "ORGANIZER"
	EventTokenStatus       EventToken = "STATUS"
	EventTokenSequence     EventToken = "SEQUENCE"
	EventTokenTransp       EventToken = "TRANSP"
	EventTokenDtstart      EventToken = "DTSTART"
	EventTokenDtend        EventToken = "DTEND"
	EventTokenUID          EventToken = "UID"
	EventTokenDTStamp      EventToken = "DTSTAMP"
	EventTokenContact      EventToken = "CONTACT"
	EventTokenLastModified EventToken = "LAST-MODIFIED"
	EventTokenComment      EventToken = "COMMENT"
	EventTokenCategories   EventToken = "CATEGORIES"
	EventTokenGeo          EventToken = "GEO"
)
