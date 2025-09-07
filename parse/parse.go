// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parse contains the logic for parsing iCalendar files and strings into Go structs
package parse

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/michael-gallo/simple-ical/model"
)

// iCalDateTimeFormat represents the standard iCal datetime format
// Format: YYYYMMDDTHHMMSSZ (e.g., 20250928T183000Z).
const iCalDateTimeFormat = "20060102T150405Z"

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
	state           *stateMachine
	currentEvent    *model.Event
	currentTimezone *model.TimeZone
	currentTodo     *model.Todo
	// Add more current* fields as needed for other components
}

// IcalString takes the string representation of an ICAL and parses it into an event
// It returns an error if the input is not a valid ICAL string.
func IcalString(input string) (*model.Calendar, error) {
	// Create parse context with all current parsing state
	ctx := &parseContext{
		state: &stateMachine{},
	}

	// TODO: add more checks for invalid calendar data
	calendar := &model.Calendar{}

	lines := strings.Split(input, "\n")

	// Handle empty input - return empty event
	if len(lines) == 0 || input == "" {
		return nil, errNoCalendarFound
	}

	// Use a state machine approach for efficiency
	for _, s := range lines {
		line := strings.TrimSpace(s)
		if line == "" || line == "\n" {
			continue
		}

		// Handle BEGIN blocks
		if beginValue, isBeginLine := strings.CutPrefix(line, "BEGIN:"); isBeginLine {
			if err := handleBeginBlock(beginValue, ctx); err != nil {
				return nil, err
			}
			continue
		}

		// Verify that the first line was a BEGIN:VCALENDAR
		if !ctx.state.inCalendar {
			return nil, errInvalidCalendarFormatMissingBegin
		}
		// Handle END blocks
		if endLineValue, isEndLine := strings.CutPrefix(line, "END:"); isEndLine {
			if err := handleEndBlock(endLineValue, ctx, calendar); err != nil {
				return nil, err
			}
			continue
		}

		// Process property lines based on current state
		if err := parsePropertyLine(line, ctx); err != nil {
			return nil, err
		}
	}

	// Verify that the last line was a END:VCALENDAR
	if ctx.state.inCalendar {
		return nil, errInvalidCalendarFormatMissingEnd
	}

	return calendar, nil
}

// parsePropertyLine parses a single property line and adds it to the appropriate component based on current state.
func parsePropertyLine(line string, ctx *parseContext) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}

	// Route to appropriate parser based on current state
	if ctx.state.inEvent {
		return parseEventProperty(line, ctx.currentEvent)
	}
	if ctx.state.inTimezone {
		return parseTimezoneProperty(line, ctx.currentTimezone)
	}
	if ctx.state.inTodo {
		return parseTodoProperty(line, ctx.currentTodo)
	}
	// Add more state checks as needed

	return nil
}

// parseEventProperty parses a single property line and adds it to the provided vevent.
func parseEventProperty(line string, event *model.Event) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}

	// parts := strings.SplitN(line, ":", 2)
	property, value, ok := strings.Cut(line, ":")
	if !ok {
		return errInvalidPropertyLine
	}

	// property := parts[0]
	// value := parts[1]

	// Handle properties that might have parameters (like ORGANIZER;CN=...)
	baseProperty := strings.Split(property, ";")[0]

	switch baseProperty {
	case "DTSTART":
		if parsedTime, err := time.Parse(iCalDateTimeFormat, value); err == nil {
			event.Start = parsedTime
		}
	case "DTEND":
		if parsedTime, err := time.Parse(iCalDateTimeFormat, value); err == nil {
			event.End = parsedTime
		}
	case "SUMMARY":
		event.Summary = value
	case "DESCRIPTION":
		event.Description = value
	case "LOCATION":
		event.Location = value
	case "STATUS":
		event.Status = model.EventStatus(value)
	case "ORGANIZER":
		organizer, err := parseOrganizer(line)
		if err != nil {
			return err
		}
		event.Organizer = organizer
	}

	return nil
}

// parseOrganizer parses a calendar line starting with ORGANIZER.
func parseOrganizer(line string) (*model.Organizer, error) {
	value, isOrganizerLine := strings.CutPrefix(line, "ORGANIZER")

	if !isOrganizerLine {
		return nil, errLineShouldStartWithOrganizer
	}

	organizer := &model.Organizer{}
	sections := strings.Split(value, ":")
	commonName, hasCommonName := strings.CutPrefix(sections[0], ";CN=")
	if hasCommonName {
		organizer.CommonName = commonName
	}

	uri := strings.Join(sections[1:], ":")
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	organizer.CalAddress = parsedURI

	return organizer, nil
}

// parseTimezoneProperty parses a single property line and adds it to the provided timezone.
func parseTimezoneProperty(line string, timezone *model.TimeZone) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return errInvalidPropertyLine
	}

	property := parts[0]
	value := parts[1]

	// Handle properties that might have parameters
	baseProperty := strings.Split(property, ";")[0]

	switch baseProperty {
	case "TZID":
		timezone.TimeZoneID = value
	case "TZOFFSETFROM":
		timezone.TimeZoneOffsetFrom = value
	case "TZOFFSETTO":
		timezone.TimeZoneOffsetTo = value
	}

	return nil
}

// parseTodoProperty parses a single property line and adds it to the provided todo.
func parseTodoProperty(line string, todo *model.Todo) error {
	// TODO: Implement todo property parsing
	// This is a placeholder for future implementation
	return nil
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
		ctx.state.inEvent = false
		calendar.Events = append(calendar.Events, *ctx.currentEvent)
	case string(model.SectionTokenVCalendar):
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
		// TODO: add standard parsing
	default:
		return fmt.Errorf("%w: %s", errTemplateInvalidEndBlock, endLineValue)
	}
	return nil
}
