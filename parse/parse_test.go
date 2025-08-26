package parse

import (
	_ "embed"
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
	assert.Equal(t, "mailto:hello@world", event.Organizer.CalAddress.String())
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
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedURIScheme, organizer.CalAddress.Scheme)
			assert.Equal(t, testCase.expectedCommonName, organizer.CommonName)
		})
	}
}
