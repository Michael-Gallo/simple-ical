// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parse

import "errors"

var (
	ErrNoVEventFound                     = errors.New("no VEVENT found in iCal data")
	ErrInvalidDateFormat                 = errors.New("invalid date format in iCal data")
	ErrMissingDTSTART                    = errors.New("missing DTSTART in VEVENT")
	ErrMissingDTEND                      = errors.New("missing DTEND in VEVENT")
	ErrInvalidTimeFormat                 = errors.New("invalid time format in iCal data")
	ErrInvalidPropertyLine               = errors.New("invalid property line in iCal data")
	ErrLineShouldStartWithOrganizer      = errors.New("line should start with ORGANIZER")
	ErrNoCalendarFound                   = errors.New("empty calendar sent")
	ErrInvalidCalendarFormatMissingBegin = errors.New("invalid calendar format: must start with BEGIN:VCALENDAR")
	ErrInvalidCalendarFormatMissingEnd   = errors.New("invalid calendar format: must end with END:VCALENDAR")

	// ErrInvalidProtocol is one of the errors that could be returned when parsing a URI with the standard library.
	ErrInvalidProtocol = errors.New("parse \"://invalid\": missing protocol scheme")
)
