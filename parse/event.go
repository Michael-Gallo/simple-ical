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
func parseEventProperty(propertyName string, value string, params map[string]string, eventIndex int16, calendar *model.Calendar) error {
	switch model.EventToken(propertyName) {
	case model.EventTokenDtstart:
		return setOnceTimeProperty(&calendar.Events[eventIndex].Start, value, propertyName, eventLocation)
	case model.EventTokenDTStamp:
		return setOnceTimeProperty(&calendar.Events[eventIndex].DTStamp, value, propertyName, eventLocation)

	// End and Duration are mutually exclusive
	case model.EventTokenDtend:
		if calendar.Events[eventIndex].Duration != 0 {
			return ErrInvalidDurationPropertyDtend
		}
		return setOnceTimeProperty(&calendar.Events[eventIndex].End, value, propertyName, eventLocation)
	case model.EventTokenDuration:
		if calendar.Events[eventIndex].End != (time.Time{}) {
			return ErrInvalidDurationPropertyDtend
		}
		return setOnceDurationProperty(&calendar.Events[eventIndex].Duration, value, propertyName, eventLocation)
	case model.EventTokenLastModified:
		return setOnceTimeProperty(&calendar.Events[eventIndex].LastModified, value, propertyName, eventLocation)

	case model.EventTokenSummary:
		return setOnceProperty(&calendar.Events[eventIndex].Summary, value, propertyName, eventLocation)
	case model.EventTokenDescription:
		return setOnceProperty(&calendar.Events[eventIndex].Description, value, propertyName, eventLocation)
	case model.EventTokenLocation:
		return setOnceProperty(&calendar.Events[eventIndex].Location, value, propertyName, eventLocation)
	case model.EventTokenUID:
		return setOnceProperty(&calendar.Events[eventIndex].UID, value, propertyName, eventLocation)
	case model.EventTokenContact:
		calendar.Events[eventIndex].Contacts = append(calendar.Events[eventIndex].Contacts, value)
		return nil

	case model.EventTokenStatus:
		calendar.Events[eventIndex].Status = model.EventStatus(value)
	case model.EventTokenTransp:
		return setOnceProperty(&calendar.Events[eventIndex].Transp, model.EventTransp(value), propertyName, eventLocation)
	case model.EventTokenSequence:
		return setOnceIntProperty(&calendar.Events[eventIndex].Sequence, value, propertyName, eventLocation)
	case model.EventTokenOrganizer:
		organizer, err := parseOrganizer(value, params)
		if err != nil {
			return err
		}
		calendar.Events[eventIndex].Organizer = organizer
	case model.EventTokenComment:
		calendar.Events[eventIndex].Comment = append(calendar.Events[eventIndex].Comment, value)
	case model.EventTokenCategories:
		calendar.Events[eventIndex].Categories = append(calendar.Events[eventIndex].Categories, strings.Split(value, ",")...)
	case model.EventTokenGeo:
		if calendar.Events[eventIndex].Geo != nil {
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
		calendar.Events[eventIndex].Geo = append(calendar.Events[eventIndex].Geo, latitude, longitude)
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
