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
}

// parseContext holds all the current parsing state for different components.
type parseContext struct {
	state                   *stateMachine
	currentEvent            *model.Event
	currentTimezone         *model.TimeZone
	currentTimeZoneProperty *model.TimeZoneProperty
	currentTodo             *model.Todo
	currentCalendar         *model.Calendar
	// Add more current* fields as needed for other components
}

// IcalString takes the string representation of an ICAL and parses it into a Calendar.
// It returns an error if the input is not a valid ICAL string.
func IcalString(input string) (*model.Calendar, error) {
	// Handle empty input
	if input == "" {
		return nil, errNoCalendarFound
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
		currentCalendar: calendar,
		state: &stateMachine{
			// We can save an Allocation by assuming the first line is a BEGIN:VCALENDAR
			// because we will immediately be checking for it
			inCalendar: true,
		},
	}
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, errNoCalendarFound
	}

	line := strings.TrimRight(scanner.Text(), "\r")
	if line != "BEGIN:VCALENDAR" {
		return nil, errInvalidCalendarFormatMissingBegin
	}

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")

		if line == "" {
			return nil, errInvalidCalendarEmptyLine
		}
		propertyName, params, value, err := parseIcalLine(line)
		if err != nil {
			return nil, err
		}
		switch propertyName {
		case "BEGIN":
			if err := handleBeginBlock(value, ctx); err != nil {
				return nil, err
			}
			continue
		case "END":
			if !ctx.state.inCalendar {
				return nil, errContentAfterEndBlock
			}
			if err := handleEndBlock(value, ctx, calendar); err != nil {
				return nil, err
			}
			continue
		default:
			if !ctx.state.inCalendar {
				return nil, errContentAfterEndBlock
			}
			if err := parsePropertyLine(propertyName, value, params, ctx); err != nil {
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
	if ctx.state.inCalendar {
		return nil, errInvalidCalendarFormatMissingEnd
	}

	return calendar, nil
}

// parsePropertyLine parses a single property line and adds it to the appropriate component based on current state.
func parsePropertyLine(propertyName string, value string, params []string, ctx *parseContext) error {
	// Route to appropriate parser based on current state
	if ctx.state.inEvent {
		return parseEventProperty(propertyName, value, params, ctx.currentEvent)
	}
	if ctx.state.inTimezone {
		return parseTimezoneProperty(propertyName, value, params, ctx)
	}
	if ctx.state.inTodo {
		return parseTodoProperty(propertyName, value, params, ctx.currentTodo)
	}

	return parseCalendarProperty(propertyName, value, params, ctx.currentCalendar)
	// Add more state checks as needed

}

// handleBeginBlock processes BEGIN blocks and updates the parser state.
func handleBeginBlock(beginValue string, ctx *parseContext) error {
	switch beginValue {
	case string(model.SectionTokenVEvent):
		ctx.state.inEvent = true
		ctx.currentEvent = &model.Event{}
	case string(model.SectionTokenVCalendar):
		ctx.state.inCalendar = true
	case string(model.SectionTokenVTimezone):
		ctx.state.inTimezone = true
		ctx.currentTimezone = &model.TimeZone{}
	case string(model.SectionTokenVFreebusy):
		ctx.state.inFreebusy = true
		// TODO: add freebusy parsing
	case string(model.SectionTokenVAlarm):
		ctx.state.inAlarm = true
		// TODO: add alarm parsing
	case string(model.SectionTokenVJournal):
		ctx.state.inJournal = true
		// TODO: add journal parsing
	case string(model.SectionTokenVTodo):
		ctx.state.inTodo = true
		*ctx.currentTodo = model.Todo{}
	case string(model.SectionTokenVStandard):
		ctx.state.inStandard = true
		ctx.currentTimeZoneProperty = &model.TimeZoneProperty{}
		// TODO: add standard parsing
	default:
		return fmt.Errorf("%w: %s", errTemplateInvalidStartBlock, beginValue)
	}
	return nil
}

// handleEndBlock processes END blocks and updates the parser state.
func handleEndBlock(endLineValue string, ctx *parseContext, calendar *model.Calendar) error {
	switch endLineValue {
	case string(model.SectionTokenVEvent):
		if err := validateEvent(ctx.currentEvent); err != nil {
			return err
		}
		ctx.state.inEvent = false
		calendar.Events = append(calendar.Events, *ctx.currentEvent)
	case string(model.SectionTokenVCalendar):
		if err := validateCalendar(calendar); err != nil {
			return err
		}
		ctx.state.inCalendar = false
	case string(model.SectionTokenVTimezone):
		ctx.state.inTimezone = false
		calendar.TimeZones = append(calendar.TimeZones, *ctx.currentTimezone)
	case string(model.SectionTokenVFreebusy):
		ctx.state.inFreebusy = false
		// TODO: add freebusy parsing
	case string(model.SectionTokenVAlarm):
		ctx.state.inAlarm = false
		// TODO: add alarm parsing
	case string(model.SectionTokenVJournal):
		ctx.state.inJournal = false
		// TODO: add journal parsing
	case string(model.SectionTokenVTodo):
		ctx.state.inTodo = false
		calendar.Todos = append(calendar.Todos, *ctx.currentTodo)
	case string(model.SectionTokenVStandard):
		ctx.state.inStandard = false
		ctx.currentTimezone.Standard = append(ctx.currentTimezone.Standard, *ctx.currentTimeZoneProperty)
	default:
		return fmt.Errorf("%w: %s", errTemplateInvalidEndBlock, endLineValue)
	}
	return nil
}
