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
			expectedError: ErrNoCalendarFound,
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
			expectedError: ErrInvalidProtocol,
		},
		{
			name:          "Calendar with no BEGIN:VCALENDAR",
			input:         "VERSION:2.0\nPRODID:-//Event//Event Calendar//EN\nCALSCALE:GREGORIAN\nMETHOD:REQUEST\nMETHOD:REQUEST",
			expectedEvent: nil,
			expectedError: ErrInvalidCalendarFormatMissingBegin,
		},
		{
			name:          "Calendar with no END:VCALENDAR",
			input:         "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Event//Event Calendar//EN\nCALSCALE:GREGORIAN\nMETHOD:REQUEST\n",
			expectedEvent: nil,
			expectedError: ErrInvalidCalendarFormatMissingEnd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, err := IcalString(tc.input)

			if tc.expectedError != nil {
				assert.ErrorContains(t, err, tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, event)

			if tc.expectedEvent != nil {
				assert.Equal(t, tc.expectedEvent.Start, event.Start)
				assert.Equal(t, tc.expectedEvent.End, event.End)
				assert.Equal(t, tc.expectedEvent.Summary, event.Summary)
				assert.Equal(t, tc.expectedEvent.Description, event.Description)
				assert.Equal(t, tc.expectedEvent.Location, event.Location)
				assert.Equal(t, tc.expectedEvent.Status, event.Status)

				if tc.expectedEvent.Organizer != nil {
					assert.Equal(t, tc.expectedEvent.Organizer.CommonName, event.Organizer.CommonName)
					assert.Equal(t, tc.expectedEvent.Organizer.CalAddress.String(), event.Organizer.CalAddress.String())
				}
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
			expectedError:      ErrLineShouldStartWithOrganizer,
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
