package parse

import (
	"fmt"
	"net/url"

	"github.com/michael-gallo/simpleical/icaldur"
	"github.com/michael-gallo/simpleical/model"
)

const timezoneLocation = "TimeZone"

// parseTimezoneProperty parses a single property line and adds it to the provided timezone.
func parseTimezoneProperty(propertyName string, value string, params map[string]string, currentState parserState, timezone *model.TimeZone) error {
	// Handle sub-components (STANDARD and DAYLIGHT)
	if currentState == stateStandard || currentState == stateDaylight {
		var tzProp *model.TimeZoneProperty
		if currentState == stateStandard {
			tzProp = &timezone.Standard[len(timezone.Standard)-1]
		} else {
			tzProp = &timezone.Daylight[len(timezone.Daylight)-1]
		}
		return parseTimeZonePropertySubComponent(propertyName, value, params, tzProp)
	}

	// Handle timezone-level properties
	switch model.TimezoneToken(propertyName) {
	case model.TimezoneTokenTimeZoneID:
		return setOnceProperty(&timezone.TimeZoneID, value, propertyName, timezoneLocation)
	case model.TimezoneTokenLastMod:
		return setOnceTimeProperty(&timezone.LastMod, value, propertyName, timezoneLocation)
	case model.TimezoneTokenTimeZoneURL:
		parsedURL, err := url.Parse(value)
		if err != nil {
			return err
		}
		return setOnceProperty(&timezone.TimeZoneURL, parsedURL, propertyName, timezoneLocation)
	default:
		return fmt.Errorf("%w: %s", errInvalidTimezoneProperty, propertyName)
	}
}

// parseTimeZonePropertySubComponent parses a single property line for STANDARD or DAYLIGHT sub-components.
func parseTimeZonePropertySubComponent(propertyName string, value string, _ map[string]string, tzProp *model.TimeZoneProperty) error {
	switch model.TimezoneToken(propertyName) {
	case model.TimezoneTokenTimeZoneOffsetFrom:
		tzProp.TimeZoneOffsetFrom = value
	case model.TimezoneTokenTimeZoneOffsetTo:
		tzProp.TimeZoneOffsetTo = value
	case model.TimezoneTokenDTStart:
		return setOnceTimeProperty(&tzProp.DTStart, value, propertyName, timezoneLocation)
	case model.TimezoneTokenComment:
		tzProp.Comment = append(tzProp.Comment, value)
	case model.TimezoneTokenRdate:
		parsedTime, err := icaldur.ParseIcalTime(value)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidTimezoneProperty, err.Error())
		}
		tzProp.Rdate = append(tzProp.Rdate, parsedTime)
	case model.TimezoneTokenTimeZoneName:
		tzProp.TimeZoneName = append(tzProp.TimeZoneName, value)
	default:
		return fmt.Errorf("%w: %s", errInvalidTimezoneProperty, propertyName)
	}
	return nil
}

// validateTimeZone ensures that all required values are present for a timezone.
func validateTimeZone(timezone *model.TimeZone) error {
	if timezone.TimeZoneID == "" {
		return errMissingTimezoneTZIDProperty
	}
	return nil
}
