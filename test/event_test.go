package test

import (
	_ "embed"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/michael-gallo/simple-ical/parse"
	"github.com/stretchr/testify/assert"
)

var (

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
)

func TestValidEvent(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
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
								DTStart:            time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, *tc.expectedCalendar, *calendar)
		})
	}
}

func TestInvalidEvent(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "Invalid organizer",
			input:         testIcalInvalidOrganizerInput,
			expectedError: parse.ErrInvalidProtocol,
		},
		{
			name:          "Invalid start date",
			input:         testIcalInvalidStartInput,
			expectedError: parse.ErrParseErrorInComponent,
		},
		{
			name:          "Invalid end date",
			input:         testIcalInvalidEndInput,
			expectedError: parse.ErrParseErrorInComponent,
		},
		{
			name:          "Content after END:VCALENDAR",
			input:         testIcalContentAfterEndBlockInput,
			expectedError: parse.ErrContentAfterEndBlock,
		},
		{
			name:          "Duplicate UID",
			input:         testIcalDuplicateUIDInput,
			expectedError: parse.ErrDuplicateProperty,
		},
		{
			name:          "Duplicate sequence",
			input:         testIcalDuplicateSequenceInput,
			expectedError: fmt.Errorf(parse.ErrDuplicatePropertyInComponentFormat, parse.ErrDuplicatePropertyInComponent, model.EventTokenSequence, "Event"),
		},
		{
			name:          "Both duration and end date are specified, DTEND first",
			input:         testIcalBothDurationAndEndInput,
			expectedError: parse.ErrInvalidDurationPropertyDtend,
		},
		{
			name:          "Both duration and end date are specified, DURATION first",
			input:         testIcalBothDurationAndEndDurationFirstInput,
			expectedError: parse.ErrInvalidDurationPropertyDtend,
		},
		{
			name:          "Missing colon in event property line",
			input:         testIcalMissingColonInput,
			expectedError: fmt.Errorf("%w: %s", parse.ErrInvalidPropertyLine, "STATUSCONFIRMED"),
		},
		{
			name:          "Missing UID",
			input:         testIcalMissingUIDInput,
			expectedError: parse.ErrMissingEventUIDProperty,
		},
		{
			name:          "Missing DTSTART",
			input:         testIcalMissingDTStartInput,
			expectedError: parse.ErrMissingEventDTStartProperty,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expectedError.Error())
			assert.Nil(t, calendar)
		})
	}
}
