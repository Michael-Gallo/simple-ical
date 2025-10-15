// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
	"time"
)

// AlarmAction represents the possible values for a VALARM's ACTION field.
// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.6.1
type AlarmAction string

const (
	AlarmActionAudio     AlarmAction = "AUDIO"
	AlarmActionDisplay   AlarmAction = "DISPLAY"
	AlarmActionEmail     AlarmAction = "EMAIL"
	AlarmActionProcedure AlarmAction = "PROCEDURE"
)

// Alarm represents a VALARM component in the iCalendar format.
// A VALARM is a grouping of component properties that defines an alarm.
// VALARM components are sub-components of VEVENT, VTODO, or VJOURNAL.
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.6.6
type Alarm struct {
	// REQUIRED, MUST NOT occur more than once
	// Defines the action to be invoked when an alarm is triggered.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.6.1
	Action AlarmAction

	// REQUIRED, MUST NOT occur more than once
	// Specifies when an alarm will trigger.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.6.3
	Trigger string

	// OPTIONAL, MUST NOT occur more than once (for AUDIO and EMAIL actions)
	// Provides the capability to associate a document object with an alarm.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.1
	Attach []string

	// OPTIONAL, MUST NOT occur more than once (for AUDIO and EMAIL actions)
	// Specifies a positive duration of time for repeating alarms.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.6.5
	Duration time.Duration

	// OPTIONAL, MUST NOT occur more than once (for DISPLAY and EMAIL actions)
	// Provides a more complete description of the alarm than that provided by the SUMMARY property.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.5
	Description []string

	// OPTIONAL, MUST NOT occur more than once (for AUDIO and EMAIL actions)
	// Defines the number of times the alarm should be repeated.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.6.4
	Repeat int

	// OPTIONAL, MUST NOT occur more than once (for EMAIL action)
	// Defines a short summary or subject for the alarm.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.1.12
	Summary string

	// OPTIONAL, MAY occur more than once (for EMAIL action, at least one required)
	// Specifies the participants that are invited to the alarm.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.1
	Attendees []url.URL

	// OPTIONAL, MAY occur more than once
	// A Non-Standard Property. Can be represented by any name with a X-prefix.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.2
	XProp map[string]string

	// OPTIONAL, MAY occur more than once
	// An IANA registered property name.
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.8.1
	IANAProp map[string]string
}
