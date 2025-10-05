// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package rrule implements the recurrence rules defined in RFC 5545
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.10
package rrule

import (
	"strconv"
	"strings"
	"time"

	"github.com/michael-gallo/simple-ical/parse"
)

type Frequency string

const (
	FrequencySecondly Frequency = "SECONDLY"
	FrequencyMinutely Frequency = "MINUTELY"
	FrequencyHourly   Frequency = "HOURLY"
	FrequencyDaily    Frequency = "DAILY"
	FrequencyWeekly   Frequency = "WEEKLY"
	FrequencyMonthly  Frequency = "MONTHLY"
	FrequencyYearly   Frequency = "YEARLY"
)

type RRule struct {
	// The frequency of the event
	// This MUST be specified
	Frequency Frequency
	// The interval between occurrences of the event
	// eg: an interval of 2 for a daily rule means the event will happen every other day
	// Not mandatory, but treated as 1 if not present
	Interval int
	// The number of occurrences of the event
	// Can not occur with the Until property
	// DTStart always counts as the first occurrence
	Count *int
	// The date and time until the rule ends, inclusive
	// Can not occur with the Count property
	Until *time.Time
}

// ParseRRule takes an iCal reccurence rule string and parses it into a RRule struct
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.10
// Example for an event that happens daily for 10 days:
// Input:
// RRULE:FREQ=DAILY;INTERVAL=1;COUNT=10
// Output:
// RRule{Frequency: FrequencyDaily, Interval: 1, Count: 10, Until: time.Time{}}
func ParseRRule(rruleString string) (*RRule, error) {
	rrule := &RRule{
		// Default to 1 if not present
		Interval: 1,
	}
	for part := range strings.SplitSeq(rruleString, ";") {
		parts := strings.Split(part, "=")
		switch parts[0] {
		case "FREQ":
			rrule.Frequency = Frequency(parts[1])
		case "INTERVAL":
			interval, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			rrule.Interval = interval
		case "COUNT":
			count, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			rrule.Count = &count
		case "UNTIL":
			until, err := parse.ParseIcalTime(parts[1])
			if err != nil {
				return nil, err
			}
			rrule.Until = &until
		}
		if len(parts) != 2 {
			return nil, ErrInvalidRRuleString
		}
	}
	if err := validateRRule(rrule); err != nil {
		return nil, err
	}
	return rrule, nil
}

func validateRRule(rrule *RRule) error {
	if rrule.Frequency == "" {
		return ErrFrequencyRequired
	}
	if rrule.Count != nil && rrule.Until != nil {
		return ErrCountAndUntilBothSet
	}
	if rrule.Interval <= 0 {
		return ErrInvalidInterval
	}
	return nil
}
