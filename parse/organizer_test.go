package parse

import (
	"net/url"
	"testing"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/stretchr/testify/assert"
)

func BenchmarkParseOrganizer(b *testing.B) {
	params := map[string]string{"CN": "My Org"}
	value := "MAILTO:dc@example.com"
	for b.Loop() {
		_, _ = parseOrganizer(value, params)
	}
}

func TestParseOrganizer(t *testing.T) {
	testCases := []struct {
		name              string
		value             string
		params            map[string]string
		expectedOrganizer *model.Organizer
		expectedError     error
	}{
		{
			name:              "Valid organizer line",
			value:             "MAILTO:dc@example.com",
			params:            map[string]string{"CN": "My Org"},
			expectedOrganizer: &model.Organizer{CommonName: "My Org", CalAddress: &url.URL{Scheme: "mailto", Opaque: "dc@example.com"}},
			expectedError:     nil,
		},
		{
			name:              "Valid organizer line with no common name",
			value:             "MAILTO:dc@example.com",
			expectedOrganizer: &model.Organizer{CalAddress: &url.URL{Scheme: "mailto", Opaque: "dc@example.com"}},
			expectedError:     nil,
		},
		{
			name:   "Mailto has a port",
			value:  "MAILTO:dc@example.com:8080",
			params: map[string]string{"CN": "My Org"},
			expectedOrganizer: &model.Organizer{
				CommonName: "My Org",
				CalAddress: &url.URL{Scheme: "mailto", Opaque: "dc@example.com:8080"},
			},
			expectedError: nil,
		},
		{
			name:   "Valid organizer line with non MAILTO URI",
			value:  "http://www.ietf.org/rfc/rfc2396.txt",
			params: map[string]string{"CN": "My Org"},
			expectedOrganizer: &model.Organizer{
				CommonName: "My Org",
				CalAddress: &url.URL{Scheme: "http", Host: "www.ietf.org", Path: "/rfc/rfc2396.txt"},
			},
			expectedError: nil,
		},
		{
			name:  "Valid organizer line with quoted string",
			value: "mailto:jsmith@example.com",
			params: map[string]string{
				"MISCFIELD":  "TEST",
				"MISCFIELD2": "TEST2",
				"CN":         "JohnSmith",
				"DIR":        "ldap://example.com:6666/o=DC%20Associates,c=US???(cn=John%20Smith)",
			},
			expectedOrganizer: &model.Organizer{
				CommonName: "JohnSmith",
				CalAddress: &url.URL{Scheme: "mailto", Opaque: "jsmith@example.com"},
				Directory:  &url.URL{Scheme: "ldap", Host: "example.com:6666", Path: "/o=DC Associates,c=US", RawQuery: "??(cn=John%20Smith)"},
				OtherParams: map[string]string{
					"MISCFIELD":  "TEST",
					"MISCFIELD2": "TEST2",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			organizer, err := parseOrganizer(testCase.value, testCase.params)
			if testCase.expectedError != nil {
				assert.ErrorIs(t, err, testCase.expectedError)
				assert.Nil(t, organizer)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedOrganizer, organizer)
		})
	}
}
