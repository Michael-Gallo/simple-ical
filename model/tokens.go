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
	SectionTokenVDaylight SectionToken = "DAYLIGHT"
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
	EventTokenDuration     EventToken = "DURATION"
	EventTokenGeo          EventToken = "GEO"
)

// TodoToken represents the names of the properties in a VTODO
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.2
type TodoToken string

const (
	TodoTokenClass           TodoToken = "CLASS"
	TodoTokenCompleted       TodoToken = "COMPLETED"
	TodoTokenCreated         TodoToken = "CREATED"
	TodoTokenDescription     TodoToken = "DESCRIPTION"
	TodoTokenDTStart         TodoToken = "DTSTART"
	TodoTokenDue             TodoToken = "DUE"
	TodoTokenDuration        TodoToken = "DURATION"
	TodoTokenGeo             TodoToken = "GEO"
	TodoTokenLastModified    TodoToken = "LAST-MODIFIED"
	TodoTokenLocation        TodoToken = "LOCATION"
	TodoTokenOrganizer       TodoToken = "ORGANIZER"
	TodoTokenPercentComplete TodoToken = "PERCENT-COMPLETE"
	TodoTokenPriority        TodoToken = "PRIORITY"
	TodoTokenRecurrenceID    TodoToken = "RECURRENCE-ID"
	TodoTokenSequence        TodoToken = "SEQUENCE"
	TodoTokenStatus          TodoToken = "STATUS"
	TodoTokenSummary         TodoToken = "SUMMARY"
	TodoTokenTransp          TodoToken = "TRANSP"
	TodoTokenURL             TodoToken = "URL"
	TodoTokenUID             TodoToken = "UID"
	TodoTokenDTStamp         TodoToken = "DTSTAMP"
	TodoTokenAttach          TodoToken = "ATTACH"
	TodoTokenAttendee        TodoToken = "ATTENDEE"
	TodoTokenCategories      TodoToken = "CATEGORIES"
	TodoTokenComment         TodoToken = "COMMENT"
	TodoTokenContact         TodoToken = "CONTACT"
	TodoTokenExceptionDates  TodoToken = "EXDATE"
	TodoTokenRequestStatus   TodoToken = "RSTATUS"
	TodoTokenRelated         TodoToken = "RELATED"
	TodoTokenResources       TodoToken = "RESOURCES"
	TodoTokenRdate           TodoToken = "RDATE"
)

// JournalToken represents the names of the properties in a VJOURNAL
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.3
type JournalToken string

const (
	JournalTokenClass          JournalToken = "CLASS"
	JournalTokenCreated        JournalToken = "CREATED"
	JournalTokenDTStart        JournalToken = "DTSTART"
	JournalTokenLastModified   JournalToken = "LAST-MODIFIED"
	JournalTokenOrganizer      JournalToken = "ORGANIZER"
	JournalTokenRecurrenceID   JournalToken = "RECURRENCE-ID"
	JournalTokenSequence       JournalToken = "SEQUENCE"
	JournalTokenStatus         JournalToken = "STATUS"
	JournalTokenSummary        JournalToken = "SUMMARY"
	JournalTokenURL            JournalToken = "URL"
	JournalTokenUID            JournalToken = "UID"
	JournalTokenDTStamp        JournalToken = "DTSTAMP"
	JournalTokenAttach         JournalToken = "ATTACH"
	JournalTokenAttendee       JournalToken = "ATTENDEE"
	JournalTokenCategories     JournalToken = "CATEGORIES"
	JournalTokenComment        JournalToken = "COMMENT"
	JournalTokenContact        JournalToken = "CONTACT"
	JournalTokenDescription    JournalToken = "DESCRIPTION"
	JournalTokenExceptionDates JournalToken = "EXDATE"
	JournalTokenRelated        JournalToken = "RELATED"
	JournalTokenRdate          JournalToken = "RDATE"
	JournalTokenRequestStatus  JournalToken = "RSTATUS"
)

// FreeBusyToken represents the names of the properties in a VFREEBUSY
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.4
type FreeBusyToken string

const (
	FreeBusyTokenContact       FreeBusyToken = "CONTACT"
	FreeBusyTokenDTStart       FreeBusyToken = "DTSTART"
	FreeBusyTokenDTEnd         FreeBusyToken = "DTEND"
	FreeBusyTokenOrganizer     FreeBusyToken = "ORGANIZER"
	FreeBusyTokenURL           FreeBusyToken = "URL"
	FreeBusyTokenUID           FreeBusyToken = "UID"
	FreeBusyTokenDTStamp       FreeBusyToken = "DTSTAMP"
	FreeBusyTokenAttendee      FreeBusyToken = "ATTENDEE"
	FreeBusyTokenComment       FreeBusyToken = "COMMENT"
	FreeBusyTokenFreeBusy      FreeBusyToken = "FREEBUSY"
	FreeBusyTokenRequestStatus FreeBusyToken = "RSTATUS"
)

// TimezoneToken represents the names of the properties in a VTIMEZONE
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.5
type TimezoneToken string

const (
	TimezoneTokenTimeZoneID         TimezoneToken = "TZID"
	TimezoneTokenLastMod            TimezoneToken = "LAST-MODIFIED"
	TimezoneTokenTimeZoneURL        TimezoneToken = "TZURL"
	TimezoneTokenTimeZoneOffsetFrom TimezoneToken = "TZOFFSETFROM"
	TimezoneTokenTimeZoneOffsetTo   TimezoneToken = "TZOFFSETTO"
	TimezoneTokenDTStart            TimezoneToken = "DTSTART"
	TimezoneTokenComment            TimezoneToken = "COMMENT"
	TimezoneTokenRdate              TimezoneToken = "RDATE"
	TimezoneTokenTimeZoneName       TimezoneToken = "TZNAME"
)

// AlarmToken represents the names of the properties in a VALARM
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.6
type AlarmToken string

const (
	AlarmTokenAction      AlarmToken = "ACTION"
	AlarmTokenTrigger     AlarmToken = "TRIGGER"
	AlarmTokenAttach      AlarmToken = "ATTACH"
	AlarmTokenDuration    AlarmToken = "DURATION"
	AlarmTokenDescription AlarmToken = "DESCRIPTION"
	AlarmTokenRepeat      AlarmToken = "REPEAT"
	AlarmTokenSummary     AlarmToken = "SUMMARY"
	AlarmTokenAttendee    AlarmToken = "ATTENDEE"
)
