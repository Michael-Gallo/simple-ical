package parse

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/michael-gallo/simpleical/model"
)

const eventLocation = "Event"

// parseEventProperty parses a single property line and adds it to the provided vevent.
func parseEventProperty(propertyName string, value string, params map[string]string, event *model.Event) error {
	switch model.EventToken(propertyName) {
	case model.EventTokenDtstart:
		return setOnceTimeProperty(&event.Start, value, propertyName, eventLocation)
	case model.EventTokenDTStamp:
		return setOnceTimeProperty(&event.DTStamp, value, propertyName, eventLocation)

	// End and Duration are mutually exclusive
	case model.EventTokenDtend:
		if event.Duration != 0 {
			return ErrInvalidDurationPropertyDtend
		}
		return setOnceTimeProperty(&event.End, value, propertyName, eventLocation)
	case model.EventTokenDuration:
		if event.End != (time.Time{}) {
			return ErrInvalidDurationPropertyDtend
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
		event.Contacts = append(event.Contacts, value)
		return nil

	case model.EventTokenStatus:
		event.Status = model.EventStatus(value)
	case model.EventTokenTransp:
		return setOnceProperty(&event.Transp, model.EventTransp(value), propertyName, eventLocation)
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
			return fmt.Errorf("%w: %s", ErrDuplicateProperty, propertyName)
		}
		// Geo must be two floats separted by a colon
		latitudeString, longitudeString, found := strings.Cut(value, ";")
		if !found {
			return ErrInvalidGeoProperty
		}
		latitude, err := strconv.ParseFloat(latitudeString, 64)
		if err != nil {
			return ErrInvalidGeoPropertyLatitude
		}
		longitude, err := strconv.ParseFloat(longitudeString, 64)
		if err != nil {
			return ErrInvalidGeoPropertyLongitude
		}
		event.Geo = append(event.Geo, latitude, longitude)
	default:
		return fmt.Errorf("%w: %s", ErrInvalidEventProperty, propertyName)
	}
	return nil
}

// parseOrganizer parses a calendar line starting with ORGANIZER.
func parseOrganizer(value string, params map[string]string) (*model.Organizer, error) {
	organizer := &model.Organizer{}
	for propName, propValue := range params {
		switch propName {
		case "CN":
			organizer.CommonName = propValue
		case "DIR":
			parsedURI, err := url.Parse(propValue)
			if err != nil {
				return nil, err
			}
			organizer.Directory = parsedURI
		case "LANGUAGE":
			organizer.Language = propValue
		case "SENT-BY":
			parsedURI, err := url.Parse(propValue)
			if err != nil {
				return nil, err
			}
			organizer.SentBy = parsedURI
		default:
			if organizer.OtherParams == nil {
				organizer.OtherParams = make(map[string]string)
			}
			organizer.OtherParams[propName] = propValue
		}
	}

	parsedURI, err := url.Parse(value)
	if err != nil {
		return nil, err
	}
	organizer.CalAddress = parsedURI

	return organizer, nil
}

// validateEvent ensures that all required values are present for an event
func validateEvent(event model.Event) error {
	if event.UID == "" {
		return ErrMissingEventUIDProperty
	}
	if event.Start.IsZero() {
		return ErrMissingEventDTStartProperty
	}
	return nil
}
