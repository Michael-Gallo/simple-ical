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
	property, value, ok := strings.Cut(line, ":")
	if !ok {
		return fmt.Errorf("%w: %s", errInvalidEventPropertyLineMissingColon, line)
	}

	// Handle properties that might have parameters (like ORGANIZER;CN=...)
	baseProperty, _, _ := strings.Cut(property, ";")

	switch model.EventToken(baseProperty) {
	case model.EventTokenDtstart:
		return setOnceTimeProperty(&event.Start, value, baseProperty, errInvalidDatePropertyDtstart)
	case model.EventTokenDTStamp:
		return setOnceTimeProperty(&event.DTStamp, value, baseProperty, errInvalidDatePropertyDTStamp)

	// End and Duration are mutually exclusive
	case model.EventTokenDtend:
		if event.Duration != 0 {
			return errInvalidDurationPropertyDtend
		}
		return setOnceTimeProperty(&event.End, value, baseProperty, errInvalidDatePropertyDtend)
	case model.EventTokenDuration:
		if event.End != (time.Time{}) {
			return errInvalidDurationPropertyDtend
		}
		return setOnceDurationProperty(&event.Duration, value, baseProperty, errInvalidDurationProperty)
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
	case model.EventTokenGeo:
		if event.Geo != nil {
			return fmt.Errorf("%w: %s", errDuplicateProperty, baseProperty)
		}
		// Geo must be two floats separted by a colon
		geo := strings.Split(value, ";")
		if len(geo) != 2 {
			return errInvalidGeoProperty
		}
		latitude, err := strconv.ParseFloat(geo[0], 64)
		if err != nil {
			return errInvalidGeoPropertyLatitude
		}
		longitude, err := strconv.ParseFloat(geo[1], 64)
		if err != nil {
			return errInvalidGeoPropertyLongitude
		}
		event.Geo = append(event.Geo, latitude, longitude)
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
