package test

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simpleical/model"
	"github.com/michael-gallo/simpleical/parse"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed test_data/calendar/valid_calendar_with_event_and_timezone.ical
	testIcalWithEventAndTimezoneInput string
	//go:embed test_data/calendar/valid_calendar.ical
	testValidCalendarInput string
	//go:embed test_data/calendar/valid_empty_calendar.ical
	testEmptyCalendarInput string
	//go:embed test_data/calendar/valid_calendar_trailing_whitespace.ical
	testTrailingWithSpaceInput string
	//go:embed test_data/calendar/no_begin_calendar.ical
	testInvalidBeginCalendarInput string
	//go:embed test_data/calendar/no_end_calendar.ical
	testInvalidEndCalendarInput string
	//go:embed test_data/calendar/empty_line_calendar.ical
	testInvalidEmptyLineCalendarInput string
	//go:embed test_data/calendar/calendar_missing_version.ical
	testCalendarMissingVersionInput string
	//go:embed test_data/calendar/calendar_missing_prodid.ical
	testCalendarMissingProdIDInput string
)

func TestParseCalendarSuccess(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
		{
			name:  "Valid iCal event",
			input: testIcalWithEventAndTimezoneInput,
			expectedCalendar: &model.Calendar{
				ProdID:   "-//Event//Event Calendar//EN",
				Version:  "2.0",
				Method:   "REQUEST",
				CalScale: "GREGORIAN",
				Events: []model.Event{
					{
						DTStamp:     time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
						UID:         "13235@example.com",
						Comment:     []string{"I Am", "A Comment"},
						Start:       time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC),
						End:         time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC),
						Summary:     "Event Summary",
						Description: "Event Description",
						Location:    "555 Fake Street",
						Organizer: &model.Organizer{
							CommonName: "Org",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "hello@world"},
						},
						Status:       model.EventStatusConfirmed,
						Sequence:     1,
						Transp:       model.EventTranspOpaque,
						Contacts:     []string{"Jim Dolittle, ABC Industries, +1-919-555-1234"},
						LastModified: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
						Categories:   []string{"first", "second", "third"},
						Geo:          []float64{37.386013, -122.082932},
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
			name:  "Valid calendar",
			input: testValidCalendarInput,
			expectedCalendar: &model.Calendar{
				ProdID:   "-//Event//Event Calendar//EN",
				Version:  "2.0",
				Method:   "REQUEST",
				CalScale: "GREGORIAN",
			},
		},
		{
			name:  "No VEVENT block",
			input: testEmptyCalendarInput,
			expectedCalendar: &model.Calendar{
				Version: "2.0",
				ProdID:  "Id",
				Events:  nil,
			},
		},
		{
			name:  "Calendar with trailing space",
			input: testTrailingWithSpaceInput,
			expectedCalendar: &model.Calendar{
				ProdID:   "-//Event//Event Calendar//EN",
				Version:  "2.0",
				Method:   "REQUEST",
				CalScale: "GREGORIAN",
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

func TestParseCalendarError(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "Calendar with no BEGIN:VCALENDAR",
			input:         testInvalidBeginCalendarInput,
			expectedError: parse.ErrInvalidCalendarFormatMissingBegin,
		},
		{
			name:          "Calendar with no END:VCALENDAR",
			input:         testInvalidEndCalendarInput,
			expectedError: parse.ErrInvalidCalendarFormatMissingEnd,
		},
		{
			name:          "Empty line in calendar",
			input:         testInvalidEmptyLineCalendarInput,
			expectedError: parse.ErrInvalidCalendarEmptyLine,
		},
		{
			name:          "Calendar missing VERSION property",
			input:         testCalendarMissingVersionInput,
			expectedError: parse.ErrMissingCalendarVersionProperty,
		},
		{
			name:          "Calendar missing PRODID property",
			input:         testCalendarMissingProdIDInput,
			expectedError: parse.ErrMissingCalendarProdIDProperty,
		},
		{
			name:          "Empty input",
			input:         "",
			expectedError: parse.ErrNoCalendarFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expectedError.Error())
			assert.Nil(t, calendar)
		})
	}
}
