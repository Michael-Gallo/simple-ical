package parse

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/michael-gallo/simple-ical/model"
)

// parseEventProperty parses a single property line and adds it to the provided vevent.
func parseEventProperty(line string, event *model.Event) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}

	property, value, ok := strings.Cut(line, ":")
	if !ok {
		return errInvalidPropertyLine
	}

	// Handle properties that might have parameters (like ORGANIZER;CN=...)
	baseProperty, _, _ := strings.Cut(property, ";")

	switch model.EventToken(baseProperty) {
	case model.EventTokenDtstart:
		parsedTime, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyDtstart
		}
		event.Start = parsedTime
	case model.EventTokenDTStamp:
		parsedTime, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyDTStamp
		}
		event.DTStamp = parsedTime
	case model.EventTokenDtend:
		parsedTime, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyDtend
		}

		event.End = parsedTime

	case model.EventTokenSummary:
		event.Summary = value
	case model.EventTokenDescription:
		event.Description = value
	case model.EventTokenLocation:
		event.Location = value
	case model.EventTokenStatus:
		event.Status = model.EventStatus(value)
	case model.EventTokenOrganizer:
		organizer, err := parseOrganizer(line)
		if err != nil {
			return err
		}
		event.Organizer = organizer
	case model.EventTokenUID:
		event.UID = value
	case model.EventTokenSequence:
		sequence, err := strconv.Atoi(value)
		if err != nil {
			return errInvalidEventPropertySequence
		}
		event.Sequence = sequence
	case model.EventTokenTransp:
		event.Transp = model.EventTransp(value)
	case model.EventTokenContact:
		event.Contact = value
	case model.EventTokenLastModified:
		lastModified, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyLastModified
		}
		event.LastModified = lastModified
	default:
		return fmt.Errorf("%w: %s", errInvalidEventProperty, baseProperty)
	}
	return nil
}

// parseOrganizer parses a calendar line starting with ORGANIZER.
func parseOrganizer(line string) (*model.Organizer, error) {
	value, isOrganizerLine := strings.CutPrefix(line, "ORGANIZER")

	if !isOrganizerLine {
		return nil, errLineShouldStartWithOrganizer
	}

	organizer := &model.Organizer{}
	sections := strings.Split(value, ":")
	commonName, hasCommonName := strings.CutPrefix(sections[0], ";CN=")
	if hasCommonName {
		organizer.CommonName = commonName
	}

	uri := strings.Join(sections[1:], ":")
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	organizer.CalAddress = parsedURI

	return organizer, nil
}
