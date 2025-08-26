package parse

import (
	_ "embed"
	"simple-ical/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:embed test_event.ical
var testIcalInput string

func TestParse(t *testing.T) {
	event, err := ParseIcalString(testIcalInput)
	assert.NoError(t, err)
	assert.NotNil(t, event)

	// Expected start time: September 28, 2025 at 18:30:00 UTC
	expectedStart := time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC)
	assert.Equal(t, expectedStart, event.Start)

	// Expected end time: September 28, 2025 at 20:30:00 UTC
	expectedEnd := time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC)
	assert.Equal(t, expectedEnd, event.End)

	assert.Equal(t, "Event Summary", event.Summary)
	assert.Equal(t, "Event Description", event.Description)
	assert.Equal(t, "555 Fake Street", event.Location)
	assert.Equal(t, "Org", event.Organizer.CommonName)
	assert.Equal(t, "hello@world", event.Organizer.CalAddress.URI)
}

func TestParseOrganizer(t *testing.T) {
	testCases := []struct {
		name              string
		line              string
		expectedOrganizer *model.Organizer
		expectedError     error
	}{
		{
			name: "Valid organizer line",
			line: "ORGANIZER;CN=My Org:MAILTO:dc@example.com",
			expectedOrganizer: &model.Organizer{
				CommonName: "My Org",
				CalAddress: model.CalendarAddress{
					URI:      "dc@example.com",
					IsMailTo: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "Valid organizer line with no common name",
			line: "ORGANIZER:MAILTO:dc@example.com",
			expectedOrganizer: &model.Organizer{
				CommonName: "",
				CalAddress: model.CalendarAddress{
					URI:      "dc@example.com",
					IsMailTo: true,
				},
			},
			expectedError: nil,
		},
		{
			name:              "Invalid Organizer line",
			line:              "Not a valid line",
			expectedOrganizer: nil,
			expectedError:     ErrLineShouldStartWithOrganizerError,
		},
		{
			name: "Mailto has a port",
			line: "ORGANIZER;CN=My Org:MAILTO:dc@example.com:8080",
			expectedOrganizer: &model.Organizer{
				CommonName: "My Org",
				CalAddress: model.CalendarAddress{
					URI:      "dc@example.com:8080",
					IsMailTo: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "Valid organizer line with non MAILTO URI",
			line: "ORGANIZER;CN=My Org:http://www.ietf.org/rfc/rfc2396.txt",
			expectedOrganizer: &model.Organizer{
				CommonName: "My Org",
				CalAddress: model.CalendarAddress{
					URI:      "http://www.ietf.org/rfc/rfc2396.txt",
					IsMailTo: false,
				},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			organizer, err := parseOrganizer(testCase.line)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedOrganizer, organizer)
		})
	}
}
