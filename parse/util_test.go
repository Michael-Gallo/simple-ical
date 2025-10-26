package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIcalLine(t *testing.T) {
	testCases := []struct {
		name                 string
		line                 string
		expectedPropertyName string
		expectedParams       map[string]string
		expectedValue        string
		expectedError        error
	}{
		{
			name:                 "Valid line",
			line:                 "DTSTART:20250928T183000Z",
			expectedPropertyName: "DTSTART",
			expectedParams:       nil,
			expectedValue:        "20250928T183000Z",
			expectedError:        nil,
		},
		{
			name:                 "Valid line with params",
			line:                 "DTSTART;VALUE=DATE:20250928",
			expectedPropertyName: "DTSTART",
			expectedParams:       map[string]string{"VALUE": "DATE"},
			expectedValue:        "20250928",
			expectedError:        nil,
		},
		{
			name:                 "Valid line with quote string",
			line:                 "ATTENDEE;ROLE=REQ-PARTICIPANT;DELEGATED-FROM=\"mailto:bob@example.com\";PARTSTAT=ACCEPTED;CN=Jane Doe:mailto:jdoe@example.com",
			expectedPropertyName: "ATTENDEE",
			expectedParams:       map[string]string{"ROLE": "REQ-PARTICIPANT", "DELEGATED-FROM": "mailto:bob@example.com", "PARTSTAT": "ACCEPTED", "CN": "Jane Doe"},
			expectedValue:        "mailto:jdoe@example.com",
			expectedError:        nil,
		},
		{
			name:                 "valid Geo Line",
			line:                 "GEO:37.386013;-122.082932",
			expectedPropertyName: "GEO",
			expectedParams:       nil,
			expectedValue:        "37.386013;-122.082932",
			expectedError:        nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			propertyName, params, value, err := parseIcalLineWithReusableMap(testCase.line, make(map[string]string))
			assert.Equal(t, testCase.expectedPropertyName, propertyName)
			assert.Equal(t, testCase.expectedParams, params)
			assert.Equal(t, testCase.expectedValue, value)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestFindUnquotedColonIndex(t *testing.T) {
	testCases := []struct {
		name          string
		line          string
		expectedIndex int
	}{
		{name: "Valid line",
			line:          "DTSTART:20250928T183000Z",
			expectedIndex: 7,
		},
		{name: "Valid line with quote string",
			line:          "\":\":",
			expectedIndex: 3,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			index := findUnquotedColonIndex(testCase.line)
			assert.Equal(t, testCase.expectedIndex, index)
		})
	}
}

func TestSplitParameters(t *testing.T) {
	testCases := []struct {
		name        string
		paramString string
		want        map[string]string
	}{
		{
			name:        "Standard organize line with a common name",
			paramString: "CN=ORG",
			want:        map[string]string{"CN": "ORG"},
		},
		{
			name:        "Organizer line with all parameters set",
			paramString: "CN=ORG;DIR=http://example.com;LANGUAGE=en;SENT-BY=mailto:org@example.com",
			want:        map[string]string{"CN": "ORG", "DIR": "http://example.com", "LANGUAGE": "en", "SENT-BY": "mailto:org@example.com"},
		},
		{
			name:        "Organizer line with quoted string",
			paramString: "CN=ORG;DIR=\"http://example.com\";LANGUAGE=en;SENT-BY=mailto:org@example.com",
			want:        map[string]string{"CN": "ORG", "DIR": "http://example.com", "LANGUAGE": "en", "SENT-BY": "mailto:org@example.com"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := splitParameters(testCase.paramString)
			assert.Equal(t, testCase.want, got)
		})
	}
}
