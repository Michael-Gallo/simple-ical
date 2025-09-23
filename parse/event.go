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
		return setOnceTimeProperty(&event.Start, value, baseProperty, errInvalidDatePropertyDtstart)
	case model.EventTokenDTStamp:
		return setOnceTimeProperty(&event.DTStamp, value, baseProperty, errInvalidDatePropertyDTStamp)
	case model.EventTokenDtend:
		return setOnceTimeProperty(&event.End, value, baseProperty, errInvalidDatePropertyDtend)
	case model.EventTokenLastModified:
		return setOnceTimeProperty(&event.LastModified, value, baseProperty, errInvalidDatePropertyLastModified)

	case model.EventTokenSummary:
		return setOnceStringProperty(&event.Summary, value, baseProperty)
	case model.EventTokenDescription:
		return setOnceStringProperty(&event.Description, value, baseProperty)
	case model.EventTokenLocation:
		return setOnceStringProperty(&event.Location, value, baseProperty)
	case model.EventTokenUID:
		return setOnceStringProperty(&event.UID, value, baseProperty)
	case model.EventTokenContact:
		return setOnceStringProperty(&event.Contact, value, baseProperty)

	case model.EventTokenStatus:
		event.Status = model.EventStatus(value)
	case model.EventTokenTransp:
		event.Transp = model.EventTransp(value)
	case model.EventTokenSequence:
		return setOnceIntProperty(&event.Sequence, value, baseProperty, errInvalidEventPropertySequence)
	case model.EventTokenOrganizer:
		organizer, err := parseOrganizer(line)
		if err != nil {
			return err
		}
		event.Organizer = organizer
	case model.EventTokenComment:
		event.Comment = append(event.Comment, value)
	case model.EventTokenCategories:
		event.Categories = append(event.Categories, strings.Split(value, ",")...)
	default:
		return fmt.Errorf("%w: %s", errInvalidEventProperty, baseProperty)
	}
	return nil
}

func setOnceIntProperty(field *int, value, propertyName string, parseError error) error {
	if *field != 0 {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return parseError
	}
	*field = parsedValue
	return nil
}

// setOnceTimeProperty sets a time.Time field only if it hasn't been set before.
func setOnceTimeProperty(field *time.Time, value, propertyName string, parseError error) error {
	if *field != (time.Time{}) {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	parsedTime, err := time.Parse(iCalDateTimeFormat, value)
	if err != nil {
		return parseError
	}
	*field = parsedTime
	return nil
}

// setOnceStringProperty sets a string field only if it hasn't been set before.
func setOnceStringProperty(field *string, value, propertyName string) error {
	if *field != "" {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	*field = value
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
