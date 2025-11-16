package parse

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/michael-gallo/simpleical/model"
)

const journalLocation = "Journal"

// parseJournalProperty parses a single property line and adds it to the provided journal.
func parseJournalProperty(propertyName string, value string, params map[string]string, journal *model.Journal) error {
	switch model.JournalToken(propertyName) {
	case model.JournalTokenDTStamp:
		return setOnceTimeProperty(&journal.DTStamp, value, propertyName, journalLocation)
	case model.JournalTokenUID:
		return setOnceProperty(&journal.UID, value, propertyName, journalLocation)
	case model.JournalTokenClass:
		return setOnceProperty(&journal.Class, model.JournalClass(value), propertyName, journalLocation)
	case model.JournalTokenCreated:
		return setOnceTimeProperty(&journal.Created, value, propertyName, journalLocation)
	case model.JournalTokenDTStart:
		return setOnceTimeProperty(&journal.DTStart, value, propertyName, journalLocation)
	case model.JournalTokenLastModified:
		return setOnceTimeProperty(&journal.LastModified, value, propertyName, journalLocation)
	case model.JournalTokenOrganizer:
		organizer, err := parseOrganizer(value, params)
		if err != nil {
			return err
		}
		journal.Organizer = organizer
	case model.JournalTokenRecurrenceID:
		return setOnceTimeProperty(&journal.RecurrenceID, value, propertyName, journalLocation)
	case model.JournalTokenSequence:
		return setOnceIntProperty(&journal.Sequence, value, propertyName, journalLocation)
	case model.JournalTokenStatus:
		journal.Status = model.JournalStatus(value)
	case model.JournalTokenSummary:
		return setOnceProperty(&journal.Summary, value, propertyName, journalLocation)
	case model.JournalTokenURL:
		return setOnceProperty(&journal.URL, value, propertyName, journalLocation)

	// Repeatable properties
	case model.JournalTokenAttach:
		journal.Attach = append(journal.Attach, value)
		return nil
	case model.JournalTokenAttendee:
		parsedURL, err := url.Parse(value)
		if err != nil {
			return err
		}
		journal.Attendees = append(journal.Attendees, *parsedURL)
	case model.JournalTokenCategories:
		journal.Categories = append(journal.Categories, strings.Split(value, ",")...)
	case model.JournalTokenComment:
		journal.Comment = append(journal.Comment, value)
	case model.JournalTokenContact:
		journal.Contacts = append(journal.Contacts, value)
	case model.JournalTokenDescription:
		journal.Description = append(journal.Description, value)
	case model.JournalTokenExceptionDates:
		return appendTimeProperty(&journal.ExceptionDates, value, propertyName, journalLocation)
	case model.JournalTokenRelated:
		journal.Related = append(journal.Related, value)
	case model.JournalTokenRdate:
		return appendTimeProperty(&journal.Rdate, value, propertyName, journalLocation)
	case model.JournalTokenRequestStatus:
		journal.RequestStatus = append(journal.RequestStatus, value)
	default:
		return fmt.Errorf("%w: %s", errInvalidJournalProperty, propertyName)
	}
	return nil
}

// validateJournal ensures that all required values are present for a journal.
func validateJournal(journal *model.Journal) error {
	if journal.UID == "" {
		return errMissingJournalUIDProperty
	}
	if time.Time.IsZero(journal.DTStart) {
		return errMissingJournalDTStartProperty
	}
	return nil
}
