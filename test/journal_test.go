package test

import (
	_ "embed"
	"net/url"
	"testing"
	"time"

	"github.com/michael-gallo/simpleical/model"
	"github.com/michael-gallo/simpleical/parse"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed test_data/journals/test_journal.ical
	testJournalInput string
	//go:embed test_data/journals/test_journal_missing_uid.ical
	testJournalMissingUIDInput string
	//go:embed test_data/journals/test_journal_duplicate_uid.ical
	testJournalDuplicateUIDInput string
	//go:embed test_data/journals/test_journal_multiple_exdates.ical
	testJournalMultipleExdatesInput string
)

func TestValidJournal(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
		{
			name:  "Valid VJOURNAL",
			input: testJournalInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Journal Calendar//EN",
				Version: "2.0",
				Journals: []model.Journal{
					{
						UID:          "journal123@example.com",
						DTStamp:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						Summary:      "Project status update",
						Description:  []string{"Completed the initial research phase", "Identified key stakeholders and requirements"},
						Class:        model.JournalClassConfidential,
						Status:       model.JournalStatusFinal,
						Created:      time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						LastModified: time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC),
						DTStart:      time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						Organizer: &model.Organizer{
							CommonName: "Project Lead",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "lead@example.com"},
						},
						Attendees:  []url.URL{{Scheme: "mailto", Opaque: "stakeholder1@example.com"}, {Scheme: "mailto", Opaque: "stakeholder2@example.com"}},
						Contacts:   []string{"Jane Doe, Project Manager, +1-555-0456"},
						Categories: []string{"work", "project", "status"},
						Comment:    []string{"This journal entry documents the completion of Phase 1"},
						URL:        "https://project.example.com/journal/123",
					},
				},
			},
		},
		{
			name:  "Valid VJOURNAL with Multiple Exception Dates",
			input: testJournalMultipleExdatesInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Journal Calendar//EN",
				Version: "2.0",
				Journals: []model.Journal{
					{
						UID:         "journal123@example.com",
						DTStamp:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						DTStart:     time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						Summary:     "Journal with Multiple Exception Dates",
						Description: []string{"This journal has multiple exception dates to test the append functionality"},
						Class:       model.JournalClassConfidential,
						Status:      model.JournalStatusFinal,
						ExceptionDates: []time.Time{
							time.Date(2024, time.January, 15, 9, 0, 0, 0, time.UTC),
							time.Date(2024, time.January, 22, 9, 0, 0, 0, time.UTC),
							time.Date(2024, time.January, 29, 9, 0, 0, 0, time.UTC),
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

func TestInvalidJournal(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "VJOURNAL missing UID",
			input: testJournalMissingUIDInput,
		},
		{
			name:  "VJOURNAL duplicate UID",
			input: testJournalDuplicateUIDInput,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendar, err := parse.IcalString(tc.input)
			assert.Error(t, err)
			assert.Nil(t, calendar)
		})
	}
}
