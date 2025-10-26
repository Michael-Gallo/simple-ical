package test

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simple-ical/model"
	"github.com/michael-gallo/simple-ical/parse"
	"github.com/stretchr/testify/assert"
)

var (

	//go:embed test_data/freebusy/test_freebusy.ical
	testFreeBusyInput string

	//go:embed test_data/freebusy/test_freebusy_missing_uid.ical
	testFreeBusyMissingUIDInput string
	//go:embed test_data/freebusy/test_freebusy_duplicate_uid.ical
	testFreeBusyDuplicateUIDInput string
	//go:embed test_data/freebusy/test_freebusy_invalid_freebusy.ical
	testFreeBusyInvalidFreeBusyInput string
)

func TestValidFreeBusy(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
		{
			name:  "Valid free busy",
			input: testFreeBusyInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//FreeBusy Calendar//EN",
				Version: "2.0",
				FreeBusys: []model.FreeBusy{
					{
						UID:     "freebusy123@example.com",
						DTStamp: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						Contact: "John Doe, Scheduling Assistant, +1-555-0123",
						DTStart: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						DTEnd:   time.Date(2024, time.January, 31, 23, 59, 59, 0, time.UTC),
						Organizer: &model.Organizer{
							CommonName: "Calendar Owner",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "owner@example.com"},
						},
						Attendees: []url.URL{{Scheme: "mailto", Opaque: "user1@example.com"}, {Scheme: "mailto", Opaque: "user2@example.com"}},
						Comment:   []string{"Available for meetings during business hours"},
						FreeBusy: []model.FreeBusyTime{
							{
								Start:  time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
								End:    time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
								Status: model.FreeBusyStatusBusy,
							},
							{
								Start:  time.Date(2024, time.January, 1, 13, 0, 0, 0, time.UTC),
								End:    time.Date(2024, time.January, 1, 17, 0, 0, 0, time.UTC),
								Status: model.FreeBusyStatusBusy,
							},
							{
								Start:  time.Date(2024, time.January, 2, 10, 0, 0, 0, time.UTC),
								End:    time.Date(2024, time.January, 2, 11, 0, 0, 0, time.UTC),
								Status: model.FreeBusyStatusBusyTentative,
							},
						},
						URL: "https://calendar.example.com/freebusy/123",
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

func TestInvalidFreeBusy(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "VFREEBUSY missing UID",
			input:         testFreeBusyMissingUIDInput,
			expectedError: parse.ErrMissingFreeBusyUIDProperty,
		},
		{
			name:          "VFREEBUSY invalid FREEBUSY format",
			input:         testFreeBusyInvalidFreeBusyInput,
			expectedError: parse.ErrInvalidFreeBusyFormat,
		},
		{
			name:          "VFREEBUSY duplicate UID",
			input:         testFreeBusyDuplicateUIDInput,
			expectedError: parse.ErrDuplicatePropertyInComponent,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.Nil(t, calendar)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
