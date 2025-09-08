package parse

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/stretchr/testify/assert"
)

//go:embed test_data/test_event.ical
var testIcalInput string

//go:embed test_data/test_event_invalid_organizer.ical
var testIcalInvalidOrganizerInput string

//go:embed test_data/test_event_invalid_start.ical
var testIcalInvalidStartInput string

//go:embed test_data/test_event_invalid_end.ical
var testIcalInvalidEndInput string

//go:embed test_data/test_event_content_after_end_block.ical
var testIcalContentAfterEndBlockInput string

func TestParse(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
		expectedError    error
	}{
		{
			name:  "Valid iCal event",
			input: testIcalInput,
			expectedCalendar: &model.Calendar{
				Events: []model.Event{
					{
						BaseComponent: model.BaseComponent{
							DTStamp: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
							UID:     "13235@example.com",
						},
						Start:       time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC),
						End:         time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC),
						Summary:     "Event Summary",
						Description: "Event Description",
						Location:    "555 Fake Street",
						Organizer: &model.Organizer{
							CommonName: "Org",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "hello@world"},
						},
						Status:   model.EventStatusConfirmed,
						Sequence: 0,
						Transp:   model.EventTranspOpaque,
					},
				},
				TimeZones: []model.TimeZone{
					{
						TimeZoneID: "America/Detroit",
						Standard: []model.TimeZoneProperty{
							{
								TimeZoneOffsetFrom: "+0000",
								TimeZoneOffsetTo:   "+0000",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:             "Empty input",
			input:            "",
			expectedCalendar: nil,
			expectedError:    errNoCalendarFound,
		},
		{
			name:  "No VEVENT block",
			input: "BEGIN:VCALENDAR\nVERSION:2.0\nEND:VCALENDAR",
			expectedCalendar: &model.Calendar{
				Events: []model.Event{},
			},
			expectedError: nil,
		},
		{
			name:             "Invalid organizer",
			input:            testIcalInvalidOrganizerInput,
			expectedCalendar: nil,
			expectedError:    errInvalidProtocol,
		},
		{
			name:             "Calendar with no BEGIN:VCALENDAR",
			input:            "VERSION:2.0\nPRODID:-//Event//Event Calendar//EN\nCALSCALE:GREGORIAN\nMETHOD:REQUEST\nMETHOD:REQUEST",
			expectedCalendar: nil,
			expectedError:    errInvalidCalendarFormatMissingBegin,
		},
		{
			name:             "Calendar with no END:VCALENDAR",
			input:            "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Event//Event Calendar//EN\nCALSCALE:GREGORIAN\nMETHOD:REQUEST\n",
			expectedCalendar: nil,
			expectedError:    errInvalidCalendarFormatMissingEnd,
		},
		{
			name:             "Invalid start date",
			input:            testIcalInvalidStartInput,
			expectedCalendar: nil,
			expectedError:    errInvalidDatePropertyDtstart,
		},
		{
			name:             "Invalid end date",
			input:            testIcalInvalidEndInput,
			expectedCalendar: nil,
			expectedError:    errInvalidDatePropertyDtend,
		},
		{
			name:             "Content after END:VCALENDAR",
			input:            testIcalContentAfterEndBlockInput,
			expectedCalendar: nil,
			expectedError:    errContentAfterEndBlock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := IcalString(tc.input)

			if tc.expectedError != nil {
				assert.ErrorContains(t, err, tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)

			if tc.expectedCalendar == nil {
				assert.Nil(t, calendar)
				return
			}

			assert.NotNil(t, calendar)
			assert.Equal(t, len(tc.expectedCalendar.Events), len(calendar.Events))
			assert.Equal(t, len(tc.expectedCalendar.Todos), len(calendar.Todos))
			assert.Equal(t, len(tc.expectedCalendar.TimeZones), len(calendar.TimeZones))

			// Compare events if they exist
			if len(tc.expectedCalendar.Events) > 0 {
				assert.Equal(t, tc.expectedCalendar.Events, calendar.Events)
			}
			if len(tc.expectedCalendar.TimeZones) > 0 {
				assert.Equal(t, tc.expectedCalendar.TimeZones, calendar.TimeZones)
			}
		})
	}
}

func TestParseOrganizer(t *testing.T) {
	testCases := []struct {
		name               string
		line               string
		expectedURIScheme  string
		expectedCommonName string
		expectedError      error
	}{
		{
			name:               "Valid organizer line",
			line:               "ORGANIZER;CN=My Org:MAILTO:dc@example.com",
			expectedCommonName: "My Org",
			expectedURIScheme:  "mailto",
			expectedError:      nil,
		},
		{
			name:               "Valid organizer line with no common name",
			line:               "ORGANIZER:MAILTO:dc@example.com",
			expectedCommonName: "",
			expectedURIScheme:  "mailto",
			expectedError:      nil,
		},
		{
			name:               "Invalid Organizer line",
			line:               "Not a valid line",
			expectedCommonName: "",
			expectedURIScheme:  "",
			expectedError:      errLineShouldStartWithOrganizer,
		},
		{
			name:               "Mailto has a port",
			line:               "ORGANIZER;CN=My Org:MAILTO:dc@example.com:8080",
			expectedCommonName: "My Org",
			expectedURIScheme:  "mailto",
			expectedError:      nil,
		},
		{
			name:               "Valid organizer line with non MAILTO URI",
			line:               "ORGANIZER;CN=My Org:http://www.ietf.org/rfc/rfc2396.txt",
			expectedCommonName: "My Org",
			expectedURIScheme:  "http",
			expectedError:      nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			organizer, err := parseOrganizer(testCase.line)
			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
				assert.Nil(t, organizer)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedURIScheme, organizer.CalAddress.Scheme)
			assert.Equal(t, testCase.expectedCommonName, organizer.CommonName)
		})
	}
}

func BenchmarkIcalString(b *testing.B) {
	for b.Loop() {
		_, _ = IcalString(testIcalInput)
	}
}

func BenchmarkParseOrganizer(b *testing.B) {
	line := "ORGANIZER;CN=My Org:MAILTO:dc@example.com"
	for b.Loop() {
		_, _ = parseOrganizer(line)
	}
}
