package parse

import (
	"fmt"
	"strings"

	"github.com/michael-gallo/simple-ical/model"
)

// parseTimezoneProperty parses a single property line and adds it to the provided timezone.
func parseTimezoneProperty(line string, ctx *parseContext) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}
	if ctx.state.inStandard {
		return parseStandardTimeZoneProperty(line, ctx.currentTimeZoneProperty)
	}

	property, value, ok := strings.Cut(line, ":")
	if !ok {
		return errInvalidPropertyLine
	}

	// Handle properties that might have parameters
	baseProperty, _, _ := strings.Cut(property, ";")

	switch baseProperty {
	case "TZID":
		ctx.currentTimezone.TimeZoneID = value
	default:
		return fmt.Errorf("%w: %s", errInvalidPropertyLine, baseProperty)
	}

	return nil
}

// parseStandardTimeZoneProperty parses a single property line and adds it to the provided standard timezone property.
func parseStandardTimeZoneProperty(line string, standard *model.TimeZoneProperty) error {
	property, value, ok := strings.Cut(line, ":")
	if !ok {
		return errInvalidPropertyLine
	}

	// Handle properties that might have parameters
	baseProperty, _, _ := strings.Cut(property, ";")

	switch baseProperty {
	case "TZOFFSETFROM":
		standard.TimeZoneOffsetFrom = value
	case "TZOFFSETTO":
		standard.TimeZoneOffsetTo = value
	}

	return nil
}
