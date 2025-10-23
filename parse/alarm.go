package parse

import (
	"fmt"
	"net/url"

	"github.com/michael-gallo/simple-ical/model"
)

const alarmLocation = "Alarm"

// parseAlarmProperty parses a single property line and adds it to the provided alarm.
func parseAlarmProperty(propertyName string, value string, params map[string]string, ctx *parseContext, calendar *model.Calendar) error {
	// Get the current alarm based on context
	var alarm *model.Alarm
	if ctx.inEvent {
		alarm = &calendar.Events[ctx.currentEventIndex].Alarms[ctx.currentAlarmIndex]
	} else if ctx.inTodo {
		alarm = &calendar.Todos[ctx.currentTodoIndex].Alarms[ctx.currentAlarmIndex]
	} else if ctx.inJournal {
		alarm = &calendar.Journals[ctx.currentJournalIndex].Alarms[ctx.currentAlarmIndex]
	}

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
func validateAlarm(ctx *parseContext, calendar *model.Calendar) error {
	// Get the current alarm based on context
	var currentAlarm *model.Alarm
	if ctx.inEvent {
		currentAlarm = &calendar.Events[ctx.currentEventIndex].Alarms[ctx.currentAlarmIndex]
	} else if ctx.inTodo {
		currentAlarm = &calendar.Todos[ctx.currentTodoIndex].Alarms[ctx.currentAlarmIndex]
	} else if ctx.inJournal {
		currentAlarm = &calendar.Journals[ctx.currentJournalIndex].Alarms[ctx.currentAlarmIndex]
	}

	if currentAlarm.Action == "" {
		return ErrMissingAlarmActionProperty
	}
	if currentAlarm.Trigger == "" {
		return ErrMissingAlarmTriggerProperty
	}

	// Validate action-specific requirements
	switch currentAlarm.Action {
	case model.AlarmActionDisplay:
		if len(currentAlarm.Description) == 0 {
			return ErrMissingAlarmDescriptionForDisplay
		}
	case model.AlarmActionEmail:
		if len(currentAlarm.Description) == 0 {
			return ErrMissingAlarmDescriptionForEmail
		}
		if currentAlarm.Summary == "" {
			return ErrMissingAlarmSummaryForEmail
		}
		if len(currentAlarm.Attendees) == 0 {
			return ErrMissingAlarmAttendeesForEmail
		}
	}

	return nil
}
