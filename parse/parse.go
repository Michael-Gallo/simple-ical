// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parse contains the logic for parsing iCalendar files and strings into Go structs
package parse

import (
	"net/url"
	"strings"
	"time"

	"github.com/michael-gallo/simple-ical/model"
)

// iCalDateTimeFormat represents the standard iCal datetime format
// Format: YYYYMMDDTHHMMSSZ (e.g., 20250928T183000Z).
const iCalDateTimeFormat = "20060102T150405Z"

// IcalString takes the string representation of an ICAL and parses it into an event
// It returns an error if the input is not a valid ICAL string.
func IcalString(input string) (*model.Event, error) {
	// TODO: add more checks for invalid calendar data
	event := &model.Event{}

	// Use a state machine approach for efficiency
	var inEvent bool

	lines := strings.SplitSeq(input, "\n")

	for s := range lines {
		line := strings.TrimSpace(s)
		if line == "" {
			continue
		}

		// Handle BEGIN blocks
		beginValue, isBeginLine := strings.CutPrefix(line, "BEGIN:")
		if isBeginLine {
			if beginValue == string(model.SectionTokenVEvent) {
				inEvent = true
			}
			continue
		}

		// Handle END blocks
		endLineValue, _ := strings.CutPrefix(line, "END:")
		if endLineValue == string(model.SectionTokenVEvent) {
			inEvent = false

			continue
		}

		// Only process lines when we're inside a VEVENT
		if inEvent {
			err := parseEventProperty(line, event)
			if err != nil {
				return nil, err
			}
			// TODO: add these to a parser for timezones
			// case "TZID":
			// 	event.TimeZoneId = value
			// case "TZOFFSETFROM":
			// 	event.TimeZoneOffsetFrom = value
			// case "TZOFFSETTO":
			// 	event.TimeZoneOffsetTo = value
		}
	}

	return event, nil
}

// parseEventProperty parses a singel property line and adds it to the provided vevent.
func parseEventProperty(line string, event *model.Event) error {
	if !strings.Contains(line, ":") {
		return ErrInvalidPropertyLine
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return ErrInvalidPropertyLine
	}

	property := parts[0]
	value := parts[1]

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
		return nil, ErrLineShouldStartWithOrganizer
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
