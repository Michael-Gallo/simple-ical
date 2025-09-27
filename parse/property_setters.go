package parse

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/michael-gallo/simple-ical/icaldur"
)

var errDuplicatePropertyInComponent = errors.New("duplicate property error")
var errParseErrorInComponent = errors.New("parse error in component")

func setOnceProperty[T comparable](field *T, value T, propertyName string, componentType string) error {
	var zero T
	if *field != zero {
		return fmt.Errorf("%w: %s set twice in component %s", errDuplicatePropertyInComponent, propertyName, componentType)
	}
	*field = value
	return nil
}

// setOnceIntProperty sets an int field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceIntProperty(field *int, value, propertyName string, componentType string) error {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", errParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, parsedValue, propertyName, componentType)
}

// setOnceTimeProperty sets a time.Time field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceTimeProperty(field *time.Time, value, propertyName string, componentType string) error {
	parsedTime, err := time.Parse(iCalDateTimeFormat, value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", errParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, parsedTime, propertyName, componentType)
}

// setOnceDurationProperty sets a duration field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceDurationProperty(field *time.Duration, value, propertyName string, componentType string) error {
	parsedDuration, err := icaldur.ParseICalDuration(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", errParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, parsedDuration, propertyName, componentType)
}
