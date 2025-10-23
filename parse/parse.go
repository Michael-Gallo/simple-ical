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

	"github.com/michael-gallo/simple-ical/model"
)

// stateMachine tracks where the parser is in a vcalendar file.
type stateMachine struct {
	inEvent    bool
	inCalendar bool
	inTimezone bool
	inTodo     bool
	inAlarm    bool
	inJournal  bool
	inFreebusy bool
	inStandard bool
	inDaylight bool
}

// parseContext holds all the current parsing state for different components.
type parseContext struct {
	stateMachine
	currentEventIndex        int16
	currentTimezoneIndex     int16
	currentTimeZonePropIndex int16
	currentTodoIndex         int16
	currentJournalIndex      int16
	currentFreeBusyIndex     int16
	// Because Alarms are sub-components of VEVENT and VTODOs, the currentAlarmIndex will reset to 0 whenever we exit an event or todo
	currentAlarmIndex int16
}

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
	// Create parse context with all current parsing state
	ctx := &parseContext{
		stateMachine: stateMachine{
			// We can save an Allocation by assuming the first line is a BEGIN:VCALENDAR
			// because we will immediately be checking for it
			inCalendar: true,
		},
	}
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, ErrNoCalendarFound
	}

	line := strings.TrimRight(scanner.Text(), "\r")
	if line != "BEGIN:VCALENDAR" {
		return nil, ErrInvalidCalendarFormatMissingBegin
	}

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")

		if line == "" {
			return nil, ErrInvalidCalendarEmptyLine
		}
		propertyName, params, value, err := parseIcalLine(line)
		if err != nil {
			return nil, err
		}
		switch propertyName {
		case "BEGIN":
			if err := handleBeginBlock(value, ctx, calendar); err != nil {
				return nil, err
			}
			continue
		case "END":
			if !ctx.inCalendar {
				return nil, ErrContentAfterEndBlock
			}
			if err := handleEndBlock(value, ctx, calendar); err != nil {
				return nil, err
			}
			continue
		default:
			if !ctx.inCalendar {
				return nil, ErrContentAfterEndBlock
			}
			if err := parsePropertyLine(propertyName, value, params, ctx, calendar); err != nil {
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
	if ctx.inCalendar {
		return nil, ErrInvalidCalendarFormatMissingEnd
	}

	return calendar, nil
}

// parsePropertyLine parses a single property line and adds it to the appropriate component based on current state.
func parsePropertyLine(propertyName string, value string, params map[string]string, ctx *parseContext, calendar *model.Calendar) error {
	// Route to appropriate parser based on current state
	// Check sub-components first (alarms can be inside events, todos, or journals)
	if ctx.inAlarm {
		var currentAlarm *model.Alarm
		if ctx.inEvent {
			currentAlarm = &calendar.Events[ctx.currentEventIndex].Alarms[ctx.currentAlarmIndex]
		} else if ctx.inTodo {
			currentAlarm = &calendar.Todos[ctx.currentTodoIndex].Alarms[ctx.currentAlarmIndex]
		}
		return parseAlarmProperty(propertyName, value, params, currentAlarm)
	}
	if ctx.inEvent {
		return parseEventProperty(propertyName, value, params, &calendar.Events[ctx.currentEventIndex])
	}
	if ctx.inTimezone {
		return parseTimezoneProperty(propertyName, value, params, ctx, &calendar.TimeZones[ctx.currentTimezoneIndex])
	}
	if ctx.inTodo {
		return parseTodoProperty(propertyName, value, params, &calendar.Todos[ctx.currentTodoIndex])
	}
	if ctx.inJournal {
		return parseJournalProperty(propertyName, value, params, &calendar.Journals[ctx.currentJournalIndex])
	}
	if ctx.inFreebusy {
		return parseFreeBusyProperty(propertyName, value, params, &calendar.FreeBusys[ctx.currentFreeBusyIndex])
	}

	return parseCalendarProperty(propertyName, value, params, calendar)
}

// handleBeginBlock processes BEGIN blocks and updates the parser state.
func handleBeginBlock(beginValue string, ctx *parseContext, calendar *model.Calendar) error {
	switch beginValue {
	case string(model.SectionTokenVEvent):
		ctx.inEvent = true
		calendar.Events = append(calendar.Events, model.Event{})
	case string(model.SectionTokenVCalendar):
		ctx.inCalendar = true
	case string(model.SectionTokenVTimezone):
		ctx.inTimezone = true
		calendar.TimeZones = append(calendar.TimeZones, model.TimeZone{})
	case string(model.SectionTokenVFreebusy):
		ctx.inFreebusy = true
		calendar.FreeBusys = append(calendar.FreeBusys, model.FreeBusy{})
	case string(model.SectionTokenVAlarm):
		ctx.inAlarm = true
		// Determine which parent component to add the alarm to
		if ctx.inEvent {
			calendar.Events[ctx.currentEventIndex].Alarms = append(calendar.Events[ctx.currentEventIndex].Alarms, model.Alarm{})
		} else if ctx.inTodo {
			calendar.Todos[ctx.currentTodoIndex].Alarms = append(calendar.Todos[ctx.currentTodoIndex].Alarms, model.Alarm{})
		} else if ctx.inJournal {
			calendar.Journals[ctx.currentJournalIndex].Alarms = append(calendar.Journals[ctx.currentJournalIndex].Alarms, model.Alarm{})
		}
	case string(model.SectionTokenVJournal):
		ctx.inJournal = true
		calendar.Journals = append(calendar.Journals, model.Journal{})
	case string(model.SectionTokenVTodo):
		ctx.inTodo = true
		calendar.Todos = append(calendar.Todos, model.Todo{})
	case string(model.SectionTokenVStandard):
		ctx.inStandard = true
		calendar.TimeZones[ctx.currentTimezoneIndex].Standard = append(calendar.TimeZones[ctx.currentTimezoneIndex].Standard, model.TimeZoneProperty{})
	case string(model.SectionTokenVDaylight):
		ctx.inDaylight = true
		calendar.TimeZones[ctx.currentTimezoneIndex].Daylight = append(calendar.TimeZones[ctx.currentTimezoneIndex].Daylight, model.TimeZoneProperty{})
	default:
		return fmt.Errorf("%w: %s", ErrTemplateInvalidStartBlock, beginValue)
	}
	return nil
}

// handleEndBlock processes END blocks and updates the parser state.
func handleEndBlock(endLineValue string, ctx *parseContext, calendar *model.Calendar) error {
	switch endLineValue {
	case string(model.SectionTokenVEvent):
		if err := validateEvent(calendar.Events[ctx.currentEventIndex]); err != nil {
			return err
		}
		ctx.currentEventIndex++
		ctx.currentAlarmIndex = 0
		ctx.inEvent = false
	case string(model.SectionTokenVCalendar):
		if err := validateCalendar(calendar); err != nil {
			return err
		}
		ctx.inCalendar = false
	case string(model.SectionTokenVTimezone):
		if err := validateTimeZone(ctx, calendar); err != nil {
			return err
		}
		ctx.inTimezone = false
		ctx.currentTimezoneIndex++
	case string(model.SectionTokenVFreebusy):
		if err := validateFreeBusy(ctx, calendar); err != nil {
			return err
		}
		ctx.inFreebusy = false
		ctx.currentFreeBusyIndex++
	case string(model.SectionTokenVAlarm):
		if err := validateAlarm(&calendar.Events[ctx.currentEventIndex].Alarms[ctx.currentAlarmIndex]); err != nil {
			return err
		}
		ctx.inAlarm = false
		ctx.currentAlarmIndex++
	case string(model.SectionTokenVJournal):
		if err := validateJournal(&calendar.Journals[ctx.currentJournalIndex]); err != nil {
			return err
		}
		ctx.inJournal = false
		ctx.currentJournalIndex++
	case string(model.SectionTokenVTodo):
		if err := validateTodo(&calendar.Todos[ctx.currentTodoIndex]); err != nil {
			return err
		}
		ctx.inTodo = false
		ctx.currentTodoIndex++
		ctx.currentAlarmIndex = 0
	case string(model.SectionTokenVStandard):
		ctx.inStandard = false
	case string(model.SectionTokenVDaylight):
		ctx.inDaylight = false
	default:
		return fmt.Errorf("%w: %s", ErrTemplateInvalidEndBlock, endLineValue)
	}
	return nil
}
