package parse

import (
	"simple-ical/model"
	"strings"
	"time"
)

// iCalDateTimeFormat represents the standard iCal datetime format
// Format: YYYYMMDDTHHMMSSZ (e.g., 20250928T183000Z)
const iCalDateTimeFormat = "20060102T150405Z"

// ParseIcalString takes the string representation of an ICAL and parses it into an event
// It returns an error if the input is not a valid ICAL string
func ParseIcalString(input string) (*model.Event, error) {
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
			if beginValue == "VEVENT" {
				inEvent = true
			}
			continue
		}

		// Handle END blocks
		endLineValue, _ := strings.CutPrefix(line, "END:")
		if endLineValue == "VEVENT" {
			inEvent = false

			continue
		}

		// Only process lines when we're inside a VEVENT
		if inEvent {
			// Parse event properties (DTSTART, DTEND, SUMMARY, etc.)
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
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
					case "ORGANIZER":
						organizer, err := parseOrganizer(line)
						if err != nil {
							return nil, err
						}
						event.Organizer = organizer
					}
				}
			}
		}
	}

	return event, nil
}

// parseOrganizer takes a calendar line starting with ORGANIZER
func parseOrganizer(line string) (*model.Organizer, error) {
	value, isOrganizerLine := strings.CutPrefix(line, "ORGANIZER")

	if !isOrganizerLine {
		return nil, ErrLineShouldStartWithOrganizerError
	}

	organizer := &model.Organizer{}
	sections := strings.Split(value, ":")
	commonName, hasCommonName := strings.CutPrefix(sections[0], ";CN=")
	if hasCommonName {
		organizer.CommonName = commonName
	}

	hasMailto := sections[1] == "MAILTO"
	if !hasMailto {
		organizer.CalAddress.IsMailTo = false
		organizer.CalAddress.URI = strings.Join(sections[1:], ":")

		return organizer, nil
	}

	// Join remaining sections, // we can not just use sections[2] because the mailing address may contain a port
	organizer.CalAddress.URI = strings.Join(sections[2:], ":")
	organizer.CalAddress.IsMailTo = true

	return organizer, nil
}
