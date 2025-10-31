// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parse contains the logic for parsing iCalendar files and strings into Go structs
package parse

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/michael-gallo/simpleical/model"
)

// ParserState represents the current parsing state using a single integer.
type ParserState uint8

const (
	StateCalendar ParserState = iota
	StateEvent
	StateTimezone
	StateTodo
	StateJournal
	StateFreebusy
	StateEventAlarm
	StateTodoAlarm
	StateStandard
	StateDaylight
	StateFinished
)

// IcalString takes the string representation of an ICAL and parses it into a Calendar.
// It returns an error if the input is not a valid ICAL string.
func IcalString(input string) (*model.Calendar, error) {
	// Handle empty input
	if input == "" {
		return nil, ErrNoCalendarFound
	}

	// Use the reader-based parser for consistency
	reader := strings.NewReader(input)
	return IcalReader(reader)
}

// IcalReader takes an io.Reader containing iCalendar data and parses it into a Calendar.
func IcalReader(reader io.Reader) (*model.Calendar, error) {
	calendar := &model.Calendar{}
	currentState := StateCalendar
	// Reusable parameter map to avoid allocations on every property
	reusableParams := make(map[string]string, 2)
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, ErrNoCalendarFound
	}

	line := strings.TrimRight(scanner.Text(), " ")
	if line != "BEGIN:VCALENDAR" {
		return nil, ErrInvalidCalendarFormatMissingBegin
	}

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " ")

		if line == "" {
			return nil, ErrInvalidCalendarEmptyLine
		}

		// Clear the reusable parameter map before each use
		for k := range reusableParams {
			delete(reusableParams, k)
		}

		propertyName, params, value, err := parseIcalLineWithReusableMap(line, reusableParams)
		if err != nil {
			return nil, err
		}
		switch propertyName {
		case "BEGIN":
			if err := handleBeginBlock(value, &currentState, calendar); err != nil {
				return nil, err
			}
			continue
		case "END":
			if currentState == StateFinished {
				return nil, ErrContentAfterEndBlock
			}
			if err := handleEndBlock(value, &currentState, calendar); err != nil {
				return nil, err
			}
			continue
		default:
			if currentState == StateFinished {
				return nil, ErrContentAfterEndBlock
			}
			if err := parsePropertyLine(propertyName, value, params, currentState, calendar); err != nil {
				return nil, err
			}
			continue
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading iCalendar data: %w", err)
	}

	// Verify that the last line was a END:VCALENDAR
	if currentState != StateFinished {
		return nil, ErrInvalidCalendarFormatMissingEnd
	}

	return calendar, nil
}

// parsePropertyLine parses a single property line and adds it to the appropriate component based on current state.
func parsePropertyLine(propertyName string, value string, params map[string]string, currentState ParserState, calendar *model.Calendar) error {
	// Route to appropriate parser based on current state
	switch currentState {
	case StateEventAlarm:
		currentAlarm := &calendar.Events[len(calendar.Events)-1].Alarms[len(calendar.Events[len(calendar.Events)-1].Alarms)-1]
		return parseAlarmProperty(propertyName, value, params, currentAlarm)
	case StateTodoAlarm:
		currentAlarm := &calendar.Todos[len(calendar.Todos)-1].Alarms[len(calendar.Todos[len(calendar.Todos)-1].Alarms)-1]
		return parseAlarmProperty(propertyName, value, params, currentAlarm)
	case StateEvent:
		return parseEventProperty(propertyName, value, params, &calendar.Events[len(calendar.Events)-1])
	case StateTimezone:
		return parseTimezoneProperty(propertyName, value, params, currentState, &calendar.TimeZones[len(calendar.TimeZones)-1])
	case StateTodo:
		return parseTodoProperty(propertyName, value, params, &calendar.Todos[len(calendar.Todos)-1])
	case StateJournal:
		return parseJournalProperty(propertyName, value, params, &calendar.Journals[len(calendar.Journals)-1])
	case StateFreebusy:
		return parseFreeBusyProperty(propertyName, value, params, &calendar.FreeBusys[len(calendar.FreeBusys)-1])
	case StateStandard, StateDaylight:
		// These are handled within timezone parsing
		return parseTimezoneProperty(propertyName, value, params, currentState, &calendar.TimeZones[len(calendar.TimeZones)-1])
	default: // StateCalendar
		return parseCalendarProperty(propertyName, value, params, calendar)
	}
}

// handleBeginBlock processes BEGIN blocks and updates the parser state.
func handleBeginBlock(beginValue string, currentState *ParserState, calendar *model.Calendar) error {
	switch beginValue {
	case string(model.SectionTokenVEvent):
		*currentState = StateEvent
		calendar.Events = append(calendar.Events, model.Event{})
	case string(model.SectionTokenVCalendar):
		*currentState = StateCalendar
	case string(model.SectionTokenVTimezone):
		*currentState = StateTimezone
		calendar.TimeZones = append(calendar.TimeZones, model.TimeZone{})
	case string(model.SectionTokenVFreebusy):
		*currentState = StateFreebusy
		calendar.FreeBusys = append(calendar.FreeBusys, model.FreeBusy{})
	case string(model.SectionTokenVAlarm):
		// Determine which parent component to add the alarm to based on current state
		switch *currentState {
		case StateEvent:
			*currentState = StateEventAlarm
			calendar.Events[len(calendar.Events)-1].Alarms = append(calendar.Events[len(calendar.Events)-1].Alarms, model.Alarm{})
		case StateTodo:
			*currentState = StateTodoAlarm
			calendar.Todos[len(calendar.Todos)-1].Alarms = append(calendar.Todos[len(calendar.Todos)-1].Alarms, model.Alarm{})
		case StateJournal:
			// Journal alarms are not supported in the current model, but we'll handle gracefully
			calendar.Journals[len(calendar.Journals)-1].Alarms = append(calendar.Journals[len(calendar.Journals)-1].Alarms, model.Alarm{})
		}
	case string(model.SectionTokenVJournal):
		*currentState = StateJournal
		calendar.Journals = append(calendar.Journals, model.Journal{})
	case string(model.SectionTokenVTodo):
		*currentState = StateTodo
		calendar.Todos = append(calendar.Todos, model.Todo{})
	case string(model.SectionTokenVStandard):
		*currentState = StateStandard
		calendar.TimeZones[len(calendar.TimeZones)-1].Standard = append(calendar.TimeZones[len(calendar.TimeZones)-1].Standard, model.TimeZoneProperty{})
	case string(model.SectionTokenVDaylight):
		*currentState = StateDaylight
		calendar.TimeZones[len(calendar.TimeZones)-1].Daylight = append(calendar.TimeZones[len(calendar.TimeZones)-1].Daylight, model.TimeZoneProperty{})
	default:
		return fmt.Errorf("%w: %s", ErrTemplateInvalidStartBlock, beginValue)
	}
	return nil
}

// handleEndBlock processes END blocks and updates the parser state.
func handleEndBlock(endLineValue string, currentState *ParserState, calendar *model.Calendar) error {
	switch endLineValue {
	case string(model.SectionTokenVEvent):
		if err := validateEvent(calendar.Events[len(calendar.Events)-1]); err != nil {
			return err
		}
		*currentState = StateCalendar
	case string(model.SectionTokenVCalendar):
		if err := validateCalendar(calendar); err != nil {
			return err
		}
		*currentState = StateFinished
	case string(model.SectionTokenVTimezone):
		if err := validateTimeZone(&calendar.TimeZones[len(calendar.TimeZones)-1]); err != nil {
			return err
		}
		*currentState = StateCalendar
	case string(model.SectionTokenVFreebusy):
		if err := validateFreeBusy(&calendar.FreeBusys[len(calendar.FreeBusys)-1]); err != nil {
			return err
		}
		*currentState = StateCalendar
	case string(model.SectionTokenVAlarm):
		// Validate alarm based on current state
		switch *currentState {
		case StateEventAlarm:
			if err := validateAlarm(&calendar.Events[len(calendar.Events)-1].Alarms[len(calendar.Events[len(calendar.Events)-1].Alarms)-1]); err != nil {
				return err
			}
			*currentState = StateEvent // Return to parent state
		case StateTodoAlarm:
			if err := validateAlarm(&calendar.Todos[len(calendar.Todos)-1].Alarms[len(calendar.Todos[len(calendar.Todos)-1].Alarms)-1]); err != nil {
				return err
			}
			*currentState = StateTodo // Return to parent state
		}
	case string(model.SectionTokenVJournal):
		if err := validateJournal(&calendar.Journals[len(calendar.Journals)-1]); err != nil {
			return err
		}
		*currentState = StateCalendar
	case string(model.SectionTokenVTodo):
		if err := validateTodo(&calendar.Todos[len(calendar.Todos)-1]); err != nil {
			return err
		}
		*currentState = StateCalendar
	case string(model.SectionTokenVStandard):
		*currentState = StateTimezone
	case string(model.SectionTokenVDaylight):
		*currentState = StateTimezone
	default:
		return fmt.Errorf("%w: %s", ErrTemplateInvalidEndBlock, endLineValue)
	}
	return nil
}
