package model

type SectionToken string

const (
	SectionTokenVCalendar SectionToken = "VCALENDAR"
	SectionTokenVEvent    SectionToken = "VEVENT"
	SectionTokenVTodo     SectionToken = "VTODO"
	SectionTokenVJournal  SectionToken = "VJOURNAL"
	SectionTokenVTimezone SectionToken = "VTIMEZONE"
	SectionTokenVFreebusy SectionToken = "VFREEBUSY"
	SectionTokenVAlarm    SectionToken = "VALARM"
)
