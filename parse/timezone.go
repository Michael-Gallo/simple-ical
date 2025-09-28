package parse

import (
	"fmt"

	"github.com/michael-gallo/simple-ical/model"
)

// parseTimezoneProperty parses a single property line and adds it to the provided timezone.
func parseTimezoneProperty(propertyName string, value string, params []string, ctx *parseContext) error {
	if ctx.state.inStandard {
		return parseStandardTimeZoneProperty(propertyName, value, params, ctx.currentTimeZoneProperty)
	}

	switch propertyName {
	case "TZID":
		ctx.currentTimezone.TimeZoneID = value
	default:
		return fmt.Errorf("%w: %s", errInvalidPropertyLine, propertyName)
	}

	return nil
}

// parseStandardTimeZoneProperty parses a single property line and adds it to the provided standard timezone property.
func parseStandardTimeZoneProperty(propertyName string, value string, _ []string, standard *model.TimeZoneProperty) error {
	switch propertyName {
	case "TZOFFSETFROM":
		standard.TimeZoneOffsetFrom = value
	case "TZOFFSETTO":
		standard.TimeZoneOffsetTo = value
	}

	return nil
}
