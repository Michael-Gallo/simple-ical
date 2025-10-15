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
		return fmt.Errorf("%w: %s", errInvalidAlarmProperty, propertyName)
	}
	return nil
}

// validateAlarm ensures that all required values are present for an alarm.
func validateAlarm(ctx *parseContext) error {
	if ctx.currentAlarm.Action == "" {
		return errMissingAlarmActionProperty
	}
	if ctx.currentAlarm.Trigger == "" {
		return errMissingAlarmTriggerProperty
	}

	// Validate action-specific requirements
	switch ctx.currentAlarm.Action {
	case model.AlarmActionDisplay:
		if len(ctx.currentAlarm.Description) == 0 {
			return errMissingAlarmDescriptionForDisplay
		}
	case model.AlarmActionEmail:
		if len(ctx.currentAlarm.Description) == 0 {
			return errMissingAlarmDescriptionForEmail
		}
		if ctx.currentAlarm.Summary == "" {
			return errMissingAlarmSummaryForEmail
		}
		if len(ctx.currentAlarm.Attendees) == 0 {
			return errMissingAlarmAttendeesForEmail
		}
	}

	return nil
}
