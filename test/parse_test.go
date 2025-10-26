package test

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/michael-gallo/simple-ical/parse"
	"github.com/stretchr/testify/assert"
)

// TODO: more heavily test multiple alarms, alarms in VEVENT and VTODOs

var (

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
