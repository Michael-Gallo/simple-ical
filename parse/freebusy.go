package parse

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/michael-gallo/simpleical/icaldur"
	"github.com/michael-gallo/simpleical/model"
)

const freeBusyLocation = "FreeBusy"

// parseFreeBusyProperty parses a single property line and adds it to the provided freebusy.
func parseFreeBusyProperty(propertyName string, value string, params map[string]string, freeBusy *model.FreeBusy) error {
	switch model.FreeBusyToken(propertyName) {
	case model.FreeBusyTokenDTStamp:
		return setOnceTimeProperty(&freeBusy.DTStamp, value, propertyName, freeBusyLocation)
	case model.FreeBusyTokenUID:
		return setOnceProperty(&freeBusy.UID, value, propertyName, freeBusyLocation)
	case model.FreeBusyTokenContact:
		return setOnceProperty(&freeBusy.Contact, value, propertyName, freeBusyLocation)
	case model.FreeBusyTokenDTStart:
		return setOnceTimeProperty(&freeBusy.DTStart, value, propertyName, freeBusyLocation)
	case model.FreeBusyTokenDTEnd:
		return setOnceTimeProperty(&freeBusy.DTEnd, value, propertyName, freeBusyLocation)
	case model.FreeBusyTokenOrganizer:
		organizer, err := parseOrganizer(value, params)
		if err != nil {
			return err
		}
		freeBusy.Organizer = organizer
	case model.FreeBusyTokenURL:
		return setOnceProperty(&freeBusy.URL, value, propertyName, freeBusyLocation)

	// Repeatable properties
	case model.FreeBusyTokenAttendee:
		parsedURL, err := url.Parse(value)
		if err != nil {
			return err
		}
		freeBusy.Attendees = append(freeBusy.Attendees, *parsedURL)
	case model.FreeBusyTokenComment:
		freeBusy.Comment = append(freeBusy.Comment, value)
	case model.FreeBusyTokenFreeBusy:
		fbTime, err := parseFreeBusyTime(value)
		if err != nil {
			return err
		}
		freeBusy.FreeBusy = append(freeBusy.FreeBusy, fbTime)
	case model.FreeBusyTokenRequestStatus:
		freeBusy.RequestStatus = append(freeBusy.RequestStatus, value)
	default:
		return fmt.Errorf("%w: %s", ErrInvalidFreeBusyProperty, propertyName)
	}
	return nil
}

// parseFreeBusyTime parses a FREEBUSY property value into a FreeBusyTime struct.
// Format: "/" separated start/end datetime pair, optionally followed by "/" and status.
// Example: "19970101T180000Z/19970102T070000Z" or "19970101T180000Z/19970102T070000Z/BUSY"
func parseFreeBusyTime(value string) (model.FreeBusyTime, error) {
	// Extract start time (everything before first '/')
	startStr, remaining, found := strings.Cut(value, "/")
	if !found {
		return model.FreeBusyTime{}, fmt.Errorf("%w: %s", ErrInvalidFreeBusyFormat, value)
	}

	startTime, err := icaldur.ParseIcalTime(startStr)
	if err != nil {
		return model.FreeBusyTime{}, fmt.Errorf("invalid start time in FREEBUSY property: %w", err)
	}

	// Extract end time and optional status (everything after first '/')
	endStr, statusStr, hasStatus := strings.Cut(remaining, "/")
	endTime, err := icaldur.ParseIcalTime(endStr)
	if err != nil {
		return model.FreeBusyTime{}, fmt.Errorf("invalid end time in FREEBUSY property: %w", err)
	}

	fbTime := model.FreeBusyTime{
		Start: startTime,
		End:   endTime,
	}

	// Optional status parameter
	if hasStatus {
		fbTime.Status = model.FreeBusyStatus(statusStr)
	} else {
		// Default to BUSY if no status specified
		fbTime.Status = model.FreeBusyStatusBusy
	}

	return fbTime, nil
}

// validateFreeBusy ensures that all required values are present for a freebusy.
func validateFreeBusy(freeBusy *model.FreeBusy) error {
	if freeBusy.UID == "" {
		return ErrMissingFreeBusyUIDProperty
	}
	if time.Time.IsZero(freeBusy.DTStart) {
		return ErrMissingFreeBusyDTStartProperty
	}
	return nil
}
