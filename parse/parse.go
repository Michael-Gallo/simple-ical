// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parse contains the logic for parsing iCalendar files and strings into Go structs
package parse

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strconv"
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
	state                   *stateMachine
	currentEvent            *model.Event
	currentTimezone         *model.TimeZone
	currentTimeZoneProperty *model.TimeZoneProperty
	currentTodo             *model.Todo
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
// This is more memory-efficient for large files as it processes data line by line.
func IcalReader(reader io.Reader) (*model.Calendar, error) {
	// Create parse context with all current parsing state
	ctx := &parseContext{
		state: &stateMachine{
			// We can save an Allocation by assuming the first line is a BEGIN:VCALENDAR
			// because we will immediately be checking for it
			inCalendar: true,
		},
	}

	calendar := &model.Calendar{}
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, errNoCalendarFound
	}

	line := strings.TrimSpace(scanner.Text())
	if line != "BEGIN:VCALENDAR" {
		return nil, errInvalidCalendarFormatMissingBegin
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Handle BEGIN blocks
		if beginValue, isBeginLine := strings.CutPrefix(line, "BEGIN:"); isBeginLine {
			if err := handleBeginBlock(beginValue, ctx); err != nil {
				return nil, err
			}
			continue
		}

		// Verify that this line is within a VCALENDAR
		if !ctx.state.inCalendar {
			return nil, errContentAfterEndBlock
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
func parsePropertyLine(line string, ctx *parseContext) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}

	// Route to appropriate parser based on current state
	if ctx.state.inEvent {
		return parseEventProperty(line, ctx.currentEvent)
	}
	if ctx.state.inTimezone {
		return parseTimezoneProperty(line, ctx)
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

	property, value, ok := strings.Cut(line, ":")
	if !ok {
		return errInvalidPropertyLine
	}

	// Handle properties that might have parameters (like ORGANIZER;CN=...)
	baseProperty, _, _ := strings.Cut(property, ";")

	switch model.EventToken(baseProperty) {
	case model.EventTokenDtstart:
		parsedTime, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyDtstart
		}
		event.Start = parsedTime
	case model.EventTokenDTStamp:
		parsedTime, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyDTStamp
		}
		event.DTStamp = parsedTime
	case model.EventTokenDtend:
		parsedTime, err := time.Parse(iCalDateTimeFormat, value)
		if err != nil {
			return errInvalidDatePropertyDtend
		}

		event.End = parsedTime

	case model.EventTokenSummary:
		event.Summary = value
	case model.EventTokenDescription:
		event.Description = value
	case model.EventTokenLocation:
		event.Location = value
	case model.EventTokenStatus:
		event.Status = model.EventStatus(value)
	case model.EventTokenOrganizer:
		organizer, err := parseOrganizer(line)
		if err != nil {
			return err
		}
		event.Organizer = organizer
	case model.EventTokenUID:
		event.UID = value
	case model.EventTokenSequence:
		sequence, err := strconv.Atoi(value)
		if err != nil {
			return errInvalidEventPropertySequence
		}
		event.Sequence = sequence
	case model.EventTokenTransp:
		event.Transp = model.EventTransp(value)
	default:
		return fmt.Errorf("%w: %s", errInvalidEventProperty, baseProperty)
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
func parseTimezoneProperty(line string, ctx *parseContext) error {
	if !strings.Contains(line, ":") {
		return errInvalidPropertyLine
	}
	if ctx.state.inStandard {
		return parseStandardTimeZoneProperty(line, ctx.currentTimeZoneProperty)
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
		ctx.currentTimezone.TimeZoneID = value
	default:
		return fmt.Errorf("%w: %s", errInvalidPropertyLine, baseProperty)
	}

	return nil
}

// parseStandardTimeZoneProperty parses a single property line and adds it to the provided standard timezone property.
func parseStandardTimeZoneProperty(line string, standard *model.TimeZoneProperty) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return errInvalidPropertyLine
	}

	property := parts[0]
	value := parts[1]

	// Handle properties that might have parameters
	baseProperty := strings.Split(property, ";")[0]

	switch baseProperty {
	case "TZOFFSETFROM":
		standard.TimeZoneOffsetFrom = value
	case "TZOFFSETTO":
		standard.TimeZoneOffsetTo = value
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
		ctx.currentTimezone.Standard = append(ctx.currentTimezone.Standard, *ctx.currentTimeZoneProperty)
	default:
		return fmt.Errorf("%w: %s", errTemplateInvalidEndBlock, endLineValue)
	}
	return nil
}
