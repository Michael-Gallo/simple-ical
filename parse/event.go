package parse

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/michael-gallo/simple-ical/model"
)

const eventLocation = "Event"

// parseEventProperty parses a single property line and adds it to the provided vevent.
func parseEventProperty(propertyName string, value string, params []string, event *model.Event) error {

	switch model.EventToken(propertyName) {
	case model.EventTokenDtstart:
		return setOnceTimeProperty(&event.Start, value, propertyName, eventLocation)
	case model.EventTokenDTStamp:
		return setOnceTimeProperty(&event.DTStamp, value, propertyName, eventLocation)

	// End and Duration are mutually exclusive
	case model.EventTokenDtend:
		if event.Duration != 0 {
			return errInvalidDurationPropertyDtend
		}
		return setOnceTimeProperty(&event.End, value, propertyName, eventLocation)
	case model.EventTokenDuration:
		if event.End != (time.Time{}) {
			return errInvalidDurationPropertyDtend
		}
		return setOnceDurationProperty(&event.Duration, value, propertyName, eventLocation)
	case model.EventTokenLastModified:
		return setOnceTimeProperty(&event.LastModified, value, propertyName, eventLocation)

	case model.EventTokenSummary:
		return setOnceProperty(&event.Summary, value, propertyName, eventLocation)
	case model.EventTokenDescription:
		return setOnceProperty(&event.Description, value, propertyName, eventLocation)
	case model.EventTokenLocation:
		return setOnceProperty(&event.Location, value, propertyName, eventLocation)
	case model.EventTokenUID:
		return setOnceProperty(&event.UID, value, propertyName, eventLocation)
	case model.EventTokenContact:
		return setOnceProperty(&event.Contact, value, propertyName, eventLocation)

	case model.EventTokenStatus:
		event.Status = model.EventStatus(value)
	case model.EventTokenTransp:
		event.Transp = model.EventTransp(value)
	case model.EventTokenSequence:
		return setOnceIntProperty(&event.Sequence, value, propertyName, eventLocation)
	case model.EventTokenOrganizer:
		organizer, err := parseOrganizer(value, params)
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
			return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
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
		return fmt.Errorf("%w: %s", errInvalidEventProperty, propertyName)
	}
	return nil
}

// parseOrganizer parses a calendar line starting with ORGANIZER.
func parseOrganizer(value string, params []string) (*model.Organizer, error) {

	organizer := &model.Organizer{}
	if params != nil {
		commonName, hasCommonName := strings.CutPrefix(params[0], "CN=")
		if hasCommonName {
			organizer.CommonName = commonName
		}
	}

	parsedURI, err := url.Parse(value)
	if err != nil {
		return nil, err
	}
	organizer.CalAddress = parsedURI

	return organizer, nil
}
