package test

import (
	_ "embed"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/michael-gallo/simple-ical/parse"
	"github.com/stretchr/testify/assert"
)

// TODO: more heavily test multiple alarms, alarms in VEVENT and VTODOs

var (

	//go:embed test_data/events/test_event_invalid_organizer.ical
	testIcalInvalidOrganizerInput string
	//go:embed test_data/events/test_event_full_organizer.ical
	testIcalFullOrganizerInput string
	//go:embed test_data/events/test_event_invalid_start.ical
	testIcalInvalidStartInput string
	//go:embed test_data/events/test_event_invalid_end.ical
	testIcalInvalidEndInput string
	//go:embed test_data/events/test_event_content_after_end_block.ical
	testIcalContentAfterEndBlockInput string
	//go:embed test_data/events/test_event_duplicate_uid.ical
	testIcalDuplicateUIDInput string
	//go:embed test_data/events/test_event_duplicate_sequence.ical
	testIcalDuplicateSequenceInput string
	//go:embed test_data/events/test_event_both_duration_and_end.ical
	testIcalBothDurationAndEndInput string
	//go:embed test_data/events/test_event_both_duration_and_end_duration_first.ical
	testIcalBothDurationAndEndDurationFirstInput string
	//go:embed test_data/events/test_event_missing_colon.ical
	testIcalMissingColonInput string
	//go:embed test_data/events/test_event_missing_uid.ical
	testIcalMissingUIDInput string
	//go:embed test_data/events/test_event_missing_dtstart.ical
	testIcalMissingDTStartInput string

	// VTODO test files
	//go:embed test_data/todos/test_todo.ical
	testTodoInput string
	//go:embed test_data/todos/test_todo_missing_uid.ical
	testTodoMissingUIDInput string
	//go:embed test_data/todos/test_todo_both_due_and_duration.ical
	testTodoBothDueAndDurationInput string
	//go:embed test_data/todos/test_todo_duplicate_uid.ical
	testTodoDuplicateUIDInput string
	//go:embed test_data/todos/test_todo_invalid_geo.ical
	testTodoInvalidGeoInput string

	// VJOURNAL test files
	//go:embed test_data/journals/test_journal.ical
	testJournalInput string
	//go:embed test_data/journals/test_journal_missing_uid.ical
	testJournalMissingUIDInput string
	//go:embed test_data/journals/test_journal_duplicate_uid.ical
	testJournalDuplicateUIDInput string
	//go:embed test_data/journals/test_journal_multiple_exdates.ical
	testJournalMultipleExdatesInput string

	// VFREEBUSY test files
	//go:embed test_data/freebusy/test_freebusy.ical
	testFreeBusyInput string
	//go:embed test_data/freebusy/test_freebusy_missing_uid.ical
	testFreeBusyMissingUIDInput string
	//go:embed test_data/freebusy/test_freebusy_duplicate_uid.ical
	testFreeBusyDuplicateUIDInput string
	//go:embed test_data/freebusy/test_freebusy_invalid_freebusy.ical
	testFreeBusyInvalidFreeBusyInput string

	// VTIMEZONE test files
	//go:embed test_data/timezones/test_timezone.ical
	testTimezoneInput string
	//go:embed test_data/timezones/test_timezone_missing_tzid.ical
	testTimezoneMissingTZIDInput string
	//go:embed test_data/timezones/test_timezone_duplicate_tzid.ical
	testTimezoneDuplicateTZIDInput string
	//go:embed test_data/timezones/test_timezone_invalid_dtstart.ical
	testTimezoneInvalidDTStartInput string

	// VALARM test files (within events)
	//go:embed test_data/events/test_event_with_alarm.ical
	testEventWithAlarmInput string
	//go:embed test_data/events/test_event_alarm_missing_action.ical
	testEventAlarmMissingActionInput string
	//go:embed test_data/events/test_event_alarm_missing_description_display.ical
	testEventAlarmMissingDescriptionDisplayInput string
	//go:embed test_data/events/test_event_alarm_missing_attendee_email.ical
	testEventAlarmMissingAttendeeEmailInput string
)

func TestParseSuccess(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
		{
			name:  "Valid organizer with all parameters set",
			input: testIcalFullOrganizerInput,
			expectedCalendar: &model.Calendar{
				ProdID:   "-//Event//Event Calendar//EN",
				Version:  "2.0",
				Method:   "REQUEST",
				CalScale: "GREGORIAN",
				Events: []model.Event{
					{
						DTStamp:     time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
						UID:         "13235@example.com",
						Start:       time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC),
						End:         time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC),
						Summary:     "Event Summary",
						Description: "Event Description",
						Location:    "555 Fake Street",
						Organizer: &model.Organizer{
							CommonName: "JohnSmith",
							Directory:  &url.URL{Scheme: "ldap", Host: "example.com:6666", Path: "/o=DC Associates,c=US", RawQuery: "??(cn=John%20Smith)"},
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "jsmith@example.com"},
							Language:   "en-us",
							SentBy:     &url.URL{Scheme: "mailto", Opaque: "mailtojsmith@example.com"},
							OtherParams: map[string]string{
								"MISCFIELD":  "TEST",
								"MISCFIELD2": "TEST2",
							},
						},
						Status:       model.EventStatusConfirmed,
						Sequence:     1,
						Comment:      []string{"I Am", "A Comment"},
						Categories:   []string{"first", "second", "third"},
						Geo:          []float64{37.386013, -122.082932},
						Transp:       model.EventTranspOpaque,
						Contacts:     []string{"Jim Dolittle, ABC Industries, +1-919-555-1234"},
						LastModified: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				TimeZones: []model.TimeZone{
					{
						TimeZoneID: "America/Detroit",
						Standard: []model.TimeZoneProperty{
							{
								TimeZoneOffsetFrom: "+0000",
								TimeZoneOffsetTo:   "+0000",
								DTStart:            time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
		{
			name:  "Valid VTODO",
			input: testTodoInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Todo Calendar//EN",
				Version: "2.0",
				Todos: []model.Todo{
					{
						UID:             "todo123@example.com",
						DTStamp:         time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						Summary:         "Complete project documentation",
						Description:     []string{"Write comprehensive documentation for the new API", "Include examples and usage patterns"},
						Location:        "Office",
						Class:           model.TodoClassConfidential,
						Status:          model.TodoStatusInProcess,
						Priority:        1,
						PercentComplete: 75,
						Created:         time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						LastModified:    time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC),
						DTStart:         time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						Due:             time.Date(2024, time.January, 30, 17, 0, 0, 0, time.UTC),
						Organizer: &model.Organizer{
							CommonName: "Project Manager",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "pm@example.com"},
						},
						Attendees:  []url.URL{{Scheme: "mailto", Opaque: "dev1@example.com"}, {Scheme: "mailto", Opaque: "dev2@example.com"}},
						Contacts:   []string{"John Doe, Engineering Team, +1-555-0123"},
						Categories: []string{"work", "urgent", "project"},
						Comment:    []string{"This is a critical task for the Q1 release"},
						Resources:  []string{"laptop", "meeting-room"},
						Geo:        []float64{37.7749, -122.4194},
						URL:        "https://project.example.com/todo/123",
					},
				},
			},
		},
		{
			name:  "Valid VJOURNAL",
			input: testJournalInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Journal Calendar//EN",
				Version: "2.0",
				Journals: []model.Journal{
					{
						UID:          "journal123@example.com",
						DTStamp:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						Summary:      "Project status update",
						Description:  []string{"Completed the initial research phase", "Identified key stakeholders and requirements"},
						Class:        model.JournalClassConfidential,
						Status:       model.JournalStatusFinal,
						Created:      time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						LastModified: time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC),
						DTStart:      time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						Organizer: &model.Organizer{
							CommonName: "Project Lead",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "lead@example.com"},
						},
						Attendees:  []url.URL{{Scheme: "mailto", Opaque: "stakeholder1@example.com"}, {Scheme: "mailto", Opaque: "stakeholder2@example.com"}},
						Contacts:   []string{"Jane Doe, Project Manager, +1-555-0456"},
						Categories: []string{"work", "project", "status"},
						Comment:    []string{"This journal entry documents the completion of Phase 1"},
						URL:        "https://project.example.com/journal/123",
					},
				},
			},
		},
		{
			name:  "Valid VJOURNAL with Multiple Exception Dates",
			input: testJournalMultipleExdatesInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Journal Calendar//EN",
				Version: "2.0",
				Journals: []model.Journal{
					{
						UID:         "journal123@example.com",
						DTStamp:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						DTStart:     time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						Summary:     "Journal with Multiple Exception Dates",
						Description: []string{"This journal has multiple exception dates to test the append functionality"},
						Class:       model.JournalClassConfidential,
						Status:      model.JournalStatusFinal,
						ExceptionDates: []time.Time{
							time.Date(2024, time.January, 15, 9, 0, 0, 0, time.UTC),
							time.Date(2024, time.January, 22, 9, 0, 0, 0, time.UTC),
							time.Date(2024, time.January, 29, 9, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		{
			name:  "Valid VFREEBUSY",
			input: testFreeBusyInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//FreeBusy Calendar//EN",
				Version: "2.0",
				FreeBusys: []model.FreeBusy{
					{
						UID:     "freebusy123@example.com",
						DTStamp: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						Contact: "John Doe, Scheduling Assistant, +1-555-0123",
						DTStart: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						DTEnd:   time.Date(2024, time.January, 31, 23, 59, 59, 0, time.UTC),
						Organizer: &model.Organizer{
							CommonName: "Calendar Owner",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "owner@example.com"},
						},
						Attendees: []url.URL{{Scheme: "mailto", Opaque: "user1@example.com"}, {Scheme: "mailto", Opaque: "user2@example.com"}},
						Comment:   []string{"Available for meetings during business hours"},
						FreeBusy: []model.FreeBusyTime{
							{
								Start:  time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
								End:    time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
								Status: model.FreeBusyStatusBusy,
							},
							{
								Start:  time.Date(2024, time.January, 1, 13, 0, 0, 0, time.UTC),
								End:    time.Date(2024, time.January, 1, 17, 0, 0, 0, time.UTC),
								Status: model.FreeBusyStatusBusy,
							},
							{
								Start:  time.Date(2024, time.January, 2, 10, 0, 0, 0, time.UTC),
								End:    time.Date(2024, time.January, 2, 11, 0, 0, 0, time.UTC),
								Status: model.FreeBusyStatusBusyTentative,
							},
						},
						URL: "https://calendar.example.com/freebusy/123",
					},
				},
			},
		},
		{
			name:  "Valid VTIMEZONE",
			input: testTimezoneInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Timezone Calendar//EN",
				Version: "2.0",
				TimeZones: []model.TimeZone{
					{
						TimeZoneID:  "America/New_York",
						LastMod:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						TimeZoneURL: &url.URL{Scheme: "http", Host: "tzurl.org", Path: "/zoneinfo-outlook/America/New_York"},
						Standard: []model.TimeZoneProperty{
							{
								TimeZoneOffsetFrom: "-0400",
								TimeZoneOffsetTo:   "-0500",
								DTStart:            time.Date(2024, time.January, 1, 2, 0, 0, 0, time.UTC),
								TimeZoneName:       []string{"EST"},
								Comment:            []string{"Eastern Standard Time"},
								Rdate:              []time.Time{time.Date(2024, time.January, 1, 2, 0, 0, 0, time.UTC)},
							},
						},
						Daylight: []model.TimeZoneProperty{
							{
								TimeZoneOffsetFrom: "-0500",
								TimeZoneOffsetTo:   "-0400",
								DTStart:            time.Date(2024, time.March, 1, 2, 0, 0, 0, time.UTC),
								TimeZoneName:       []string{"EDT"},
								Comment:            []string{"Eastern Daylight Time"},
								Rdate:              []time.Time{time.Date(2024, time.March, 1, 2, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
		},
		{
			name:  "Valid VEVENT with VALARM",
			input: testEventWithAlarmInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Event//Event Calendar//EN",
				Version: "2.0",
				Events: []model.Event{
					{
						UID:         "13235@example.com",
						DTStamp:     time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
						Start:       time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC),
						End:         time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC),
						Summary:     "Event with Alarm",
						Description: "Event Description",
						Alarms: []model.Alarm{
							{
								Action:      model.AlarmActionDisplay,
								Trigger:     "-PT15M",
								Description: []string{"Reminder: Event starting in 15 minutes"},
								Repeat:      2,
								Duration:    5 * time.Minute,
							},
							{
								Action:      model.AlarmActionEmail,
								Trigger:     "-PT1H",
								Description: []string{"Email reminder for upcoming event"},
								Summary:     "Event Reminder",
								Attendees:   []url.URL{{Scheme: "mailto", Opaque: "user@example.com"}},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, *tc.expectedCalendar, *calendar)
		})
	}
}

func TestParseError(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "Empty input",
			input:         "",
			expectedError: parse.ErrNoCalendarFound,
		},
		{
			name:          "Invalid organizer",
			input:         testIcalInvalidOrganizerInput,
			expectedError: parse.ErrInvalidProtocol,
		},
		{
			name:          "Invalid start date",
			input:         testIcalInvalidStartInput,
			expectedError: parse.ErrParseErrorInComponent,
		},
		{
			name:          "Invalid end date",
			input:         testIcalInvalidEndInput,
			expectedError: parse.ErrParseErrorInComponent,
		},
		{
			name:          "Content after END:VCALENDAR",
			input:         testIcalContentAfterEndBlockInput,
			expectedError: parse.ErrContentAfterEndBlock,
		},
		{
			name:          "Duplicate UID",
			input:         testIcalDuplicateUIDInput,
			expectedError: parse.ErrDuplicateProperty,
		},
		{
			name:          "Duplicate sequence",
			input:         testIcalDuplicateSequenceInput,
			expectedError: fmt.Errorf(parse.ErrDuplicatePropertyInComponentFormat, parse.ErrDuplicatePropertyInComponent, model.EventTokenSequence, "Event"),
		},
		{
			name:          "Both duration and end date are specified, DTEND first",
			input:         testIcalBothDurationAndEndInput,
			expectedError: parse.ErrInvalidDurationPropertyDtend,
		},
		{
			name:          "Both duration and end date are specified, DURATION first",
			input:         testIcalBothDurationAndEndDurationFirstInput,
			expectedError: parse.ErrInvalidDurationPropertyDtend,
		},
		{
			name:          "Missing colon in event property line",
			input:         testIcalMissingColonInput,
			expectedError: fmt.Errorf("%w: %s", parse.ErrInvalidPropertyLine, "STATUSCONFIRMED"),
		},
		{
			name:          "Missing UID",
			input:         testIcalMissingUIDInput,
			expectedError: parse.ErrMissingEventUIDProperty,
		},
		{
			name:          "Missing DTSTART",
			input:         testIcalMissingDTStartInput,
			expectedError: parse.ErrMissingEventDTStartProperty,
		},
		{
			name:          "VTODO missing UID",
			input:         testTodoMissingUIDInput,
			expectedError: parse.ErrMissingTodoUIDProperty,
		},
		{
			name:          "VTODO both DUE and DURATION",
			input:         testTodoBothDueAndDurationInput,
			expectedError: parse.ErrInvalidDurationPropertyDue,
		},
		{
			name:          "VTODO invalid GEO",
			input:         testTodoInvalidGeoInput,
			expectedError: parse.ErrInvalidGeoProperty,
		},
		{
			name:          "VJOURNAL missing UID",
			input:         testJournalMissingUIDInput,
			expectedError: parse.ErrMissingJournalUIDProperty,
		},
		{
			name:          "VFREEBUSY missing UID",
			input:         testFreeBusyMissingUIDInput,
			expectedError: parse.ErrMissingFreeBusyUIDProperty,
		},
		{
			name:          "VFREEBUSY invalid FREEBUSY format",
			input:         testFreeBusyInvalidFreeBusyInput,
			expectedError: parse.ErrInvalidFreeBusyFormat,
		},
		{
			name:          "VTIMEZONE missing TZID",
			input:         testTimezoneMissingTZIDInput,
			expectedError: parse.ErrMissingTimezoneTZIDProperty,
		},
		{
			name:          "VTIMEZONE invalid DTSTART",
			input:         testTimezoneInvalidDTStartInput,
			expectedError: parse.ErrInvalidTimezoneProperty,
		},
		{
			name:          "VALARM missing ACTION",
			input:         testEventAlarmMissingActionInput,
			expectedError: parse.ErrMissingAlarmActionProperty,
		},
		{
			name:          "VALARM DISPLAY missing DESCRIPTION",
			input:         testEventAlarmMissingDescriptionDisplayInput,
			expectedError: parse.ErrMissingAlarmDescriptionForDisplay,
		},
		{
			name:          "VALARM EMAIL missing ATTENDEE",
			input:         testEventAlarmMissingAttendeeEmailInput,
			expectedError: parse.ErrMissingAlarmAttendeesForEmail,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.ErrorContains(t, err, tc.expectedError.Error())
			assert.Nil(t, calendar)
		})
	}
}

func BenchmarkIcalString(b *testing.B) {
	for b.Loop() {
		_, _ = parse.IcalString(testIcalWithEventAndTimezoneInput)
	}
}
