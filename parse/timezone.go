package parse

import (
	"fmt"
	"net/url"
	"time"

	"github.com/michael-gallo/simple-ical/model"
)

const timezoneLocation = "TimeZone"

// parseTimezoneProperty parses a single property line and adds it to the provided timezone.
func parseTimezoneProperty(propertyName string, value string, params map[string]string, ctx *parseContext) error {
	// Handle sub-components (STANDARD and DAYLIGHT)
	if ctx.state.inStandard || ctx.state.inDaylight {
		return parseTimeZonePropertySubComponent(propertyName, value, params, ctx.currentTimeZoneProperty)
	}

	// Handle timezone-level properties
	switch model.TimezoneToken(propertyName) {
	case model.TimezoneTokenTimeZoneID:
		return setOnceProperty(&ctx.currentTimezone.TimeZoneID, value, propertyName, timezoneLocation)
	case model.TimezoneTokenLastMod:
		return setOnceTimeProperty(&ctx.currentTimezone.LastMod, value, propertyName, timezoneLocation)
	case model.TimezoneTokenTimeZoneURL:
		parsedURL, err := url.Parse(value)
		if err != nil {
			return err
		}
		return setOnceProperty(&ctx.currentTimezone.TimeZoneURL, parsedURL, propertyName, timezoneLocation)
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
		parsedTime, err := parseTimezoneTime(value)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidTimezoneProperty, err.Error())
		}
		tzProp.DTStart = parsedTime
	case model.TimezoneTokenComment:
		tzProp.Comment = append(tzProp.Comment, value)
	case model.TimezoneTokenRdate:
		parsedTime, err := parseTimezoneTime(value)
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

// parseTimezoneTime parses a datetime value that may or may not have a Z suffix.
func parseTimezoneTime(value string) (time.Time, error) {
	// Try with Z format first
	if time, err := time.Parse("20060102T150405Z", value); err == nil {
		return time, nil
	}
	// Try without Z format
	if time, err := time.Parse("20060102T150405", value); err == nil {
		return time, nil
	}
	return time.Time{}, fmt.Errorf("%w: %s", errInvalidTimezoneDatetimeFormat, value)
}

// validateTimeZone ensures that all required values are present for a timezone.
func validateTimeZone(ctx *parseContext) error {
	if ctx.currentTimezone.TimeZoneID == "" {
		return errMissingTimezoneTZIDProperty
	}
	return nil
}
