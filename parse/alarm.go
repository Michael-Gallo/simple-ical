package parse

import (
	"fmt"
	"net/url"

	"github.com/michael-gallo/simple-ical/model"
)

const alarmLocation = "Alarm"

// parseAlarmProperty parses a single property line and adds it to the provided alarm.
func parseAlarmProperty(propertyName string, value string, params map[string]string, alarm *model.Alarm) error {
	switch model.AlarmToken(propertyName) {
	case model.AlarmTokenAction:
		return setOnceProperty(&alarm.Action, model.AlarmAction(value), propertyName, alarmLocation)
	case model.AlarmTokenTrigger:
		return setOnceProperty(&alarm.Trigger, value, propertyName, alarmLocation)
	case model.AlarmTokenAttach:
		alarm.Attach = append(alarm.Attach, value)
		return nil
	case model.AlarmTokenDuration:
		return setOnceDurationProperty(&alarm.Duration, value, propertyName, alarmLocation)
	case model.AlarmTokenDescription:
		alarm.Description = append(alarm.Description, value)
		return nil
	case model.AlarmTokenRepeat:
		return setOnceIntProperty(&alarm.Repeat, value, propertyName, alarmLocation)
	case model.AlarmTokenSummary:
		return setOnceProperty(&alarm.Summary, value, propertyName, alarmLocation)
	case model.AlarmTokenAttendee:
		parsedURL, err := url.Parse(value)
		if err != nil {
			return err
		}
		alarm.Attendees = append(alarm.Attendees, *parsedURL)
	default:
		return fmt.Errorf("%w: %s", ErrInvalidAlarmProperty, propertyName)
	}
	return nil
}

// validateAlarm ensures that all required values are present for an alarm.
func validateAlarm(alarm *model.Alarm) error {
	if alarm.Action == "" {
		return ErrMissingAlarmActionProperty
	}
	if alarm.Trigger == "" {
		return ErrMissingAlarmTriggerProperty
	}

	// Validate action-specific requirements
	switch alarm.Action {
	case model.AlarmActionDisplay:
		if len(alarm.Description) == 0 {
			return ErrMissingAlarmDescriptionForDisplay
		}
	case model.AlarmActionEmail:
		if len(alarm.Description) == 0 {
			return ErrMissingAlarmDescriptionForEmail
		}
		if alarm.Summary == "" {
			return ErrMissingAlarmSummaryForEmail
		}
		if len(alarm.Attendees) == 0 {
			return ErrMissingAlarmAttendeesForEmail
		}
	}

	return nil
}
