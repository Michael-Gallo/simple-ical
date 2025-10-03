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
)

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
			input: testEmptyCalendarInput,
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
			input:            testInvalidBeginCalendarInput,
			expectedCalendar: nil,
			expectedError:    errInvalidCalendarFormatMissingBegin,
		},
		{
			name:             "Calendar with no END:VCALENDAR",
			input:            testInvalidEndCalendarInput,
			expectedCalendar: nil,
			expectedError:    errInvalidCalendarFormatMissingEnd,
		},
		{
			name:             "Invalid start date",
			input:            testIcalInvalidStartInput,
			expectedCalendar: nil,
			expectedError:    errParseErrorInComponent,
		},
		{
			name:             "Invalid end date",
			input:            testIcalInvalidEndInput,
			expectedCalendar: nil,
			expectedError:    errParseErrorInComponent,
		},
		{
			name:             "Content after END:VCALENDAR",
			input:            testIcalContentAfterEndBlockInput,
			expectedCalendar: nil,
			expectedError:    errContentAfterEndBlock,
		},
		{
			name:             "Duplicate UID",
			input:            testIcalDuplicateUIDInput,
			expectedCalendar: nil,
			expectedError:    errDuplicateProperty,
		},
		{
			name:             "Duplicate sequence",
			input:            testIcalDuplicateSequenceInput,
			expectedCalendar: nil,
			expectedError:    fmt.Errorf(errDuplicatePropertyInComponentFormat, errDuplicatePropertyInComponent, model.EventTokenSequence, eventLocation),
		},
		{
			name:             "Both duration and end date are specified, DTEND first",
			input:            testIcalBothDurationAndEndInput,
			expectedCalendar: nil,
			expectedError:    errInvalidDurationPropertyDtend,
		},
		{
			name:             "Both duration and end date are specified, DURATION first",
			input:            testIcalBothDurationAndEndDurationFirstInput,
			expectedCalendar: nil,
			expectedError:    errInvalidDurationPropertyDtend,
		},
		{
			name:             "Missing colon in event property line",
			input:            testIcalMissingColonInput,
			expectedCalendar: nil,
			expectedError:    fmt.Errorf("%w: %s", errInvalidPropertyLine, "STATUSCONFIRMED"),
		},
		{
			name:             "Missing UID",
			input:            testIcalMissingUIDInput,
			expectedCalendar: nil,
			expectedError:    errMissingEventUIDProperty,
		},
		{
			name:             "Missing DTSTART",
			input:            testIcalMissingDTStartInput,
			expectedCalendar: nil,
			expectedError:    errMissingEventDTStartProperty,
		},
		{
			name:             "Empty line in calendar",
			input:            testInvalidEmptyLineCalendarInput,
			expectedCalendar: nil,
			expectedError:    errInvalidCalendarEmptyLine,
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
		value              string
		params             []string
		expectedURIScheme  string
		expectedCommonName string
		expectedError      error
	}{
		{
			name:               "Valid organizer line",
			value:              "MAILTO:dc@example.com",
			params:             []string{"CN=My Org"},
			expectedCommonName: "My Org",
			expectedURIScheme:  "mailto",
			expectedError:      nil,
		},
		{
			name:               "Valid organizer line with no common name",
			value:              "MAILTO:dc@example.com",
			expectedCommonName: "",
			expectedURIScheme:  "mailto",
			expectedError:      nil,
		},
		{
			name:               "Mailto has a port",
			value:              "MAILTO:dc@example.com:8080",
			params:             []string{"CN=My Org"},
			expectedCommonName: "My Org",
			expectedURIScheme:  "mailto",
			expectedError:      nil,
		},
		{
			name:               "Valid organizer line with non MAILTO URI",
			value:              "http://www.ietf.org/rfc/rfc2396.txt",
			params:             []string{"CN=My Org"},
			expectedCommonName: "My Org",
			expectedURIScheme:  "http",
			expectedError:      nil,
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
	params := []string{"CN=My Org"}
	value := "MAILTO:dc@example.com"
	for b.Loop() {
		_, _ = parseOrganizer(value, params)
	}
}
