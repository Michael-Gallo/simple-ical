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

	//go:embed test_data/todos/test_todo.ical
	testTodoInput string
	//go:embed test_data/todos/test_todo_missing_uid.ical
	testTodoMissingUIDInput string
	//go:embed test_data/todos/test_todo_both_due_and_duration.ical
	testTodoBothDueAndDurationInput string
	//go:embed test_data/todos/test_todo_duplicate_uid.ical
	testTodoDuplicateUIDInput string
	//go:embed test_data/todos/test_todo_invalid_geo.ical
	testTodoInvalidGeoInput string
)

func TestValidTodo(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedCalendar *model.Calendar
	}{
		{
			name:  "Valid VTODO",
			input: testTodoInput,
			expectedCalendar: &model.Calendar{
				ProdID:  "-//Test//Todo Calendar//EN",
				Version: "2.0",
				Todos: []model.Todo{
					{
						UID:             "todo123@example.com",
						DTStamp:         time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						Summary:         "Complete project documentation",
						Description:     []string{"Write comprehensive documentation for the new API", "Include examples and usage patterns"},
						Location:        "Office",
						Class:           model.TodoClassConfidential,
						Status:          model.TodoStatusInProcess,
						Priority:        1,
						PercentComplete: 75,
						Created:         time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
						LastModified:    time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC),
						DTStart:         time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC),
						Due:             time.Date(2024, time.January, 30, 17, 0, 0, 0, time.UTC),
						Organizer: &model.Organizer{
							CommonName: "Project Manager",
							CalAddress: &url.URL{Scheme: "mailto", Opaque: "pm@example.com"},
						},
						Attendees:  []url.URL{{Scheme: "mailto", Opaque: "dev1@example.com"}, {Scheme: "mailto", Opaque: "dev2@example.com"}},
						Contacts:   []string{"John Doe, Engineering Team, +1-555-0123"},
						Categories: []string{"work", "urgent", "project"},
						Comment:    []string{"This is a critical task for the Q1 release"},
						Resources:  []string{"laptop", "meeting-room"},
						Geo:        []float64{37.7749, -122.4194},
						URL:        "https://project.example.com/todo/123",
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

func TestInvalidTodo(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "VTODO missing UID",
			input: testTodoMissingUIDInput,
		},
		{
			name:  "VTODO both DUE and DURATION",
			input: testTodoBothDueAndDurationInput,
		},
		{
			name:  "VTODO invalid GEO",
			input: testTodoInvalidGeoInput,
		},
		{
			name:  "VTODO duplicate UID",
			input: testTodoDuplicateUIDInput,
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
