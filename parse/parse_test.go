package parse

import (
	_ "embed"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/stretchr/testify/assert"
)

var (

	//go:embed test_data/events/test_event.ical
	testIcalInput string
	//go:embed test_data/events/test_event_invalid_organizer.ical
	testIcalInvalidOrganizerInput string
	//go:embed test_data/events/test_event_full_organizer.ical
	testIcalFullOrganizerInput string
	//go:embed test_data/events/test_event_invalid_start.ical
	testIcalInvalidStartInput string
	//go:embed test_data/events/test_event_invalid_end.ical
	testIcalInvalidEndInput string
	//go:embed test_data/events/test_event_content_after_end_block.ical
	testIcalContentAfterEndBlockInput string
	//go:embed test_data/events/test_event_duplicate_uid.ical
	testIcalDuplicateUIDInput string
	//go:embed test_data/events/test_event_duplicate_sequence.ical
	testIcalDuplicateSequenceInput string
	//go:embed test_data/events/test_event_both_duration_and_end.ical
	testIcalBothDurationAndEndInput string
	//go:embed test_data/events/test_event_both_duration_and_end_duration_first.ical
	testIcalBothDurationAndEndDurationFirstInput string
	//go:embed test_data/events/test_event_missing_colon.ical
	testIcalMissingColonInput string
	//go:embed test_data/events/test_event_missing_uid.ical
	testIcalMissingUIDInput string
	//go:embed test_data/events/test_event_missing_dtstart.ical
	testIcalMissingDTStartInput string

	//go:embed test_data/empty_calendar.ical
	testEmptyCalendarInput string
	//go:embed test_data/no_begin_calendar.ical
	testInvalidBeginCalendarInput string
	//go:embed test_data/no_end_calendar.ical
	testInvalidEndCalendarInput string
	//go:embed test_data/empty_line_calendar.ical
	testInvalidEmptyLineCalendarInput string
	//go:embed test_data/valid_calendar.ical
	testValidCalendarInput string
	//go:embed test_data/calendar_missing_version.ical
	testCalendarMissingVersionInput string
	//go:embed test_data/calendar_missing_prodid.ical
	testCalendarMissingProdIDInput string
)

func TestParseSuccess(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
		{
			name:  "Valid iCal event",
			input: testIcalInput,
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
							},
						},
					},
				},
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
			name:  "Valid organizer with all parameters set",
			input: testIcalFullOrganizerInput,
			expectedCalendar: &model.Calendar{
				ProdID:   "-//Event//Event Calendar//EN",
				Version:  "2.0",
				Method:   "REQUEST",
				CalScale: "GREGORIAN",
				Events: []model.Event{
					{
						DTStamp:     time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
						UID:         "13235@example.com",
						Start:       time.Date(2025, time.September, 28, 18, 30, 0, 0, time.UTC),
						End:         time.Date(2025, time.September, 28, 20, 30, 0, 0, time.UTC),
						Summary:     "Event Summary",
						Description: "Event Description",
						Location:    "555 Fake Street",
						Organizer: &model.Organizer{
							CommonName: "JohnSmith",
							Directory:  &url.URL{Scheme: "ldap", Host: "example.com:6666", Path: "/o=DC Associates,c=US", RawQuery: "??(cn=John%20Smith)"},
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "jsmith@example.com"},
							Language:   "en-us",
							SentBy:     &url.URL{Scheme: "mailto", Opaque: "mailtojsmith@example.com"},
							OtherParams: map[string]string{
								"MISCFIELD":  "TEST",
								"MISCFIELD2": "TEST2",
							},
						},
						Status:       model.EventStatusConfirmed,
						Sequence:     1,
						Comment:      []string{"I Am", "A Comment"},
						Categories:   []string{"first", "second", "third"},
						Geo:          []float64{37.386013, -122.082932},
						Transp:       model.EventTranspOpaque,
						Contacts:     []string{"Jim Dolittle, ABC Industries, +1-919-555-1234"},
						LastModified: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
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
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := IcalString(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, *tc.expectedCalendar, *calendar)
		})
	}
}

func TestParseError(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "Empty input",
			input:         "",
			expectedError: errNoCalendarFound,
		},
		{
			name:          "Invalid organizer",
			input:         testIcalInvalidOrganizerInput,
			expectedError: errInvalidProtocol,
		},
		{
			name:          "Calendar with no BEGIN:VCALENDAR",
			input:         testInvalidBeginCalendarInput,
			expectedError: errInvalidCalendarFormatMissingBegin,
		},
		{
			name:          "Calendar with no END:VCALENDAR",
			input:         testInvalidEndCalendarInput,
			expectedError: errInvalidCalendarFormatMissingEnd,
		},
		{
			name:          "Invalid start date",
			input:         testIcalInvalidStartInput,
			expectedError: errParseErrorInComponent,
		},
		{
			name:          "Invalid end date",
			input:         testIcalInvalidEndInput,
			expectedError: errParseErrorInComponent,
		},
		{
			name:          "Content after END:VCALENDAR",
			input:         testIcalContentAfterEndBlockInput,
			expectedError: errContentAfterEndBlock,
		},
		{
			name:          "Duplicate UID",
			input:         testIcalDuplicateUIDInput,
			expectedError: errDuplicateProperty,
		},
		{
			name:          "Duplicate sequence",
			input:         testIcalDuplicateSequenceInput,
			expectedError: fmt.Errorf(errDuplicatePropertyInComponentFormat, errDuplicatePropertyInComponent, model.EventTokenSequence, eventLocation),
		},
		{
			name:          "Both duration and end date are specified, DTEND first",
			input:         testIcalBothDurationAndEndInput,
			expectedError: errInvalidDurationPropertyDtend,
		},
		{
			name:          "Both duration and end date are specified, DURATION first",
			input:         testIcalBothDurationAndEndDurationFirstInput,
			expectedError: errInvalidDurationPropertyDtend,
		},
		{
			name:          "Missing colon in event property line",
			input:         testIcalMissingColonInput,
			expectedError: fmt.Errorf("%w: %s", errInvalidPropertyLine, "STATUSCONFIRMED"),
		},
		{
			name:          "Missing UID",
			input:         testIcalMissingUIDInput,
			expectedError: errMissingEventUIDProperty,
		},
		{
			name:          "Missing DTSTART",
			input:         testIcalMissingDTStartInput,
			expectedError: errMissingEventDTStartProperty,
		},
		{
			name:          "Empty line in calendar",
			input:         testInvalidEmptyLineCalendarInput,
			expectedError: errInvalidCalendarEmptyLine,
		},
		{
			name:          "Calendar missing VERSION property",
			input:         testCalendarMissingVersionInput,
			expectedError: errMissingCalendarVersionProperty,
		},
		{
			name:          "Calendar missing PRODID property",
			input:         testCalendarMissingProdIDInput,
			expectedError: errMissingCalendarProdIDProperty,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := IcalString(tc.input)
			assert.ErrorContains(t, err, tc.expectedError.Error())
			assert.Nil(t, calendar)
		})
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

func BenchmarkIcalString(b *testing.B) {
	for b.Loop() {
		_, _ = IcalString(testIcalInput)
	}
}

func BenchmarkParseOrganizer(b *testing.B) {
	params := map[string]string{"CN": "My Org"}
	value := "MAILTO:dc@example.com"
	for b.Loop() {
		_, _ = parseOrganizer(value, params)
	}
}
