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
		expectedParams       []string
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
			expectedParams:       []string{"VALUE=DATE"},
			expectedValue:        "20250928",
			expectedError:        nil,
		},
		{
			name:                 "Valid line with quote string",
			line:                 "ATTENDEE;ROLE=REQ-PARTICIPANT;DELEGATED-FROM=\"mailto:bob@example.com\";PARTSTAT=ACCEPTED;CN=Jane Doe:mailto:jdoe@example.com",
			expectedPropertyName: "ATTENDEE",
			expectedParams:       []string{"ROLE=REQ-PARTICIPANT", "DELEGATED-FROM=\"mailto:bob@example.com\"", "PARTSTAT=ACCEPTED", "CN=Jane Doe"},
			expectedValue:        "mailto:jdoe@example.com",
			expectedError:        nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			propertyName, params, value, err := parseIcalLine(testCase.line)
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
