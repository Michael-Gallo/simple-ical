package parse

import (
	"fmt"
	"strconv"
	"time"

	"github.com/michael-gallo/simple-ical/icaldur"
)

func setOnceIntProperty(field *int, value, propertyName string, parseError error) error {
	if *field != 0 {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return parseError
	}
	*field = parsedValue
	return nil
}

// setOnceTimeProperty sets a time.Time field only if it hasn't been set before.
func setOnceTimeProperty(field *time.Time, value, propertyName string, parseError error) error {
	if *field != (time.Time{}) {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	parsedTime, err := time.Parse(iCalDateTimeFormat, value)
	if err != nil {
		return parseError
	}
	*field = parsedTime
	return nil
}

// setOnceStringProperty sets a string field only if it hasn't been set before.
func setOnceStringProperty(field *string, value, propertyName string) error {
	if *field != "" {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	*field = value
	return nil
}

// setOnceDurationProperty sets a duration field only if it hasn't been set before.
func setOnceDurationProperty(field *time.Duration, value, propertyName string, parseError error) error {
	if *field != 0 {
		return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
	}
	parsedDuration, err := icaldur.ParseICalDuration(value)
	if err != nil {
		return parseError
	}
	*field = parsedDuration
	return nil
}
