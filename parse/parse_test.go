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

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedEvent *model.Event
		expectedError error
	}{
		{
			name:  "Valid iCal event",
			input: testIcalInput,
			expectedEvent: &model.Event{
				Start:       time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC),
				End:         time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC),
				Summary:     "Event Summary",
				Description: "Event Description",
				Location:    "555 Fake Street",
				Organizer: &model.Organizer{
					CommonName: "Org",
					CalAddress: &url.URL{Scheme: "mailto", Opaque: "hello@world"},
				},
				Status: model.EventStatusConfirmed,
			},
			expectedError: nil,
		},
		{
			name:          "Empty input",
			input:         "",
			expectedEvent: nil,
			expectedError: errNoCalendarFound,
		},
		{
			name:          "No VEVENT block",
			input:         "BEGIN:VCALENDAR\nVERSION:2.0\nEND:VCALENDAR",
			expectedEvent: &model.Event{},
			expectedError: nil,
		},
		{
			name:          "Invalid organizer",
			input:         testIcalInvalidOrganizerInput,
			expectedEvent: nil,
			expectedError: errInvalidProtocol,
		},
		{
			name:          "Calendar with no BEGIN:VCALENDAR",
			input:         "VERSION:2.0\nPRODID:-//Event//Event Calendar//EN\nCALSCALE:GREGORIAN\nMETHOD:REQUEST\nMETHOD:REQUEST",
			expectedEvent: nil,
			expectedError: errInvalidCalendarFormatMissingBegin,
		},
		{
			name:          "Calendar with no END:VCALENDAR",
			input:         "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Event//Event Calendar//EN\nCALSCALE:GREGORIAN\nMETHOD:REQUEST\n",
			expectedEvent: nil,
			expectedError: errInvalidCalendarFormatMissingEnd,
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

			if tc.expectedEvent == nil {
				assert.Nil(t, calendar)
				return
			}
			assert.NotNil(t, calendar)
			assert.Equal(t, tc.expectedEvent.Start, calendar.Events[0].Start)
			assert.Equal(t, tc.expectedEvent.End, calendar.Events[0].End)
			assert.Equal(t, tc.expectedEvent.Summary, calendar.Events[0].Summary)
			assert.Equal(t, tc.expectedEvent.Description, calendar.Events[0].Description)
			assert.Equal(t, tc.expectedEvent.Location, calendar.Events[0].Location)
			assert.Equal(t, tc.expectedEvent.Status, calendar.Events[0].Status)

			if tc.expectedEvent.Organizer != nil {
				assert.Equal(t, tc.expectedEvent.Organizer.CommonName, calendar.Events[0].Organizer.CommonName)
				assert.Equal(t, tc.expectedEvent.Organizer.CalAddress.String(), calendar.Events[0].Organizer.CalAddress.String())
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
