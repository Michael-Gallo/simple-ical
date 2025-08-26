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
	assert.Equal(t, "hello@world", event.Organizer.Mailto)
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
				Mailto:     "dc@example.com",
			},
			expectedError: nil,
		},
		{
			name: "Valid organizer line with no common name",
			line: "ORGANIZER:MAILTO:dc@example.com",
			expectedOrganizer: &model.Organizer{
				CommonName: "",
				Mailto:     "dc@example.com",
			},
			expectedError: nil,
		},
		{
			name:              "Invalid Organizer line",
			line:              "Not a valid line",
			expectedOrganizer: nil,
			expectedError:     ErrLineShouldStartWithOrganizerError,
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
