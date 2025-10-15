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

var (
	//go:embed test_data/calendar/valid_calendar_with_event_and_timezone.ical
	testIcalWithEventAndTimezoneInput string
	//go:embed test_data/calendar/valid_calendar.ical
	testValidCalendarInput string
	//go:embed test_data/calendar/valid_empty_calendar.ical
	testEmptyCalendarInput string
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
	}{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Nil(t, calendar)
		})
	}
}
