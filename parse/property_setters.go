package parse

import (
	"fmt"
	"strconv"
	"time"

	"github.com/michael-gallo/simple-ical/icaldur"
)

// iCalDateTimeFormat represents the standard iCal datetime format
// Format: YYYYMMDDTHHMMSSZ (e.g., 20250928T183000Z).
const iCalDateTimeFormat = "20060102T150405Z"

func parseIcalTime(value string) (time.Time, error) {
	return time.Parse(iCalDateTimeFormat, value)
}

// setOncePropertyWithParse ensures that set-once properties that require string parsing have consistent error handling
func setOncePropertyWithParse[T comparable](field *T, value string, propertyName string, componentType string, parseFunc func(string) (T, error)) error {
	parsedValue, err := parseFunc(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", errParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, parsedValue, propertyName, componentType)
}

// setOnceProperty ensures that set-once properties have consistent error handling
func setOnceProperty[T comparable](field *T, value T, propertyName string, componentType string) error {
	var zero T
	if *field != zero {
		return fmt.Errorf(errDuplicatePropertyInComponentFormat, errDuplicatePropertyInComponent, propertyName, componentType)
	}
	*field = value
	return nil
}

// setOnceIntProperty sets an int field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceIntProperty(field *int, value, propertyName string, componentType string) error {
	return setOncePropertyWithParse(field, value, propertyName, componentType, strconv.Atoi)
}

// setOnceTimeProperty sets a time.Time field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceTimeProperty(field *time.Time, value, propertyName string, componentType string) error {
	return setOncePropertyWithParse(field, value, propertyName, componentType, parseIcalTime)
}

// setOnceDurationProperty sets a duration field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceDurationProperty(field *time.Duration, value, propertyName string, componentType string) error {
	return setOncePropertyWithParse(field, value, propertyName, componentType, icaldur.ParseICalDuration)
}
