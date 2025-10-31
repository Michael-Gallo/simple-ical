package parse

import (
	"fmt"
	"strconv"
	"time"

	"github.com/michael-gallo/simpleical/icaldur"
)

// setOnceProperty ensures that set-once properties have consistent error handling
func setOnceProperty[T comparable](field *T, value T, propertyName string, componentType string) error {
	var zero T
	if *field != zero {
		return fmt.Errorf(ErrDuplicatePropertyInComponentFormat, ErrDuplicatePropertyInComponent, propertyName, componentType)
	}
	*field = value
	return nil
}

// setOnceIntProperty sets an int field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceIntProperty(field *int, value, propertyName string, componentType string) error {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", ErrParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, intValue, propertyName, componentType)
}

// setOnceTimeProperty sets a time.Time field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceTimeProperty(field *time.Time, value, propertyName string, componentType string) error {
	time, err := icaldur.ParseIcalTime(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", ErrParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, time, propertyName, componentType)
}

// setOnceDurationProperty sets a duration field only if it hasn't been set before.
// this is intended for properties that according to the spec must only be set once
func setOnceDurationProperty(field *time.Duration, value, propertyName string, componentType string) error {
	duration, err := icaldur.ParseICalDuration(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", ErrParseErrorInComponent, componentType, propertyName)
	}
	return setOnceProperty(field, duration, propertyName, componentType)
}

func appendTimeProperty(field *[]time.Time, value, propertyName string, componentType string) error {
	time, err := icaldur.ParseIcalTime(value)
	if err != nil {
		return fmt.Errorf("%w: %s property %s in iCal", ErrParseErrorInComponent, componentType, propertyName)
	}
	*field = append(*field, time)
	return nil
}
