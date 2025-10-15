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

	"github.com/michael-gallo/simple-ical/icaldur"
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

type Weekday string

const (
	WeekdayMonday    Weekday = "MO"
	WeekdayTuesday   Weekday = "TU"
	WeekdayWednesday Weekday = "WE"
	WeekdayThursday  Weekday = "TH"
	WeekdayFriday    Weekday = "FR"
	WeekdaySaturday  Weekday = "SA"
	WeekdaySunday    Weekday = "SU"
)

type ByDay struct {
	// The day of the week that the event occurs on
	Weekday Weekday
	// The interval between occurrences of the event
	// eg: If Weekday is Tuesday, and Interval is 2, then the event will happen every other Tuesday
	Interval int
}

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
	// The day of the week that the event occurs on
	// This is optional and repeatable
	Weekday []ByDay

	// The Month(s) of the year that the event occurs on
	Month []int

	// The day of the month that the event occurs on
	// eg: 10th of the month, negative numbers are allowed to indicate the last day of the month
	// for example, -3 is the third-to-last-day of the month
	Monthday []int

	// The day of the year that the event occurs on
	// eg: 100th day of the year, negative numbers are allowed to indicate the last day of the year
	YearDay []int
}

// ParseRRule takes an iCal reccurence rule string and parses it into a RRule struct
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.10
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.5.3
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
		tag, value, found := strings.Cut(part, "=")
		if !found {
			return nil, ErrInvalidRRuleString
		}
		switch tag {
		case "FREQ":
			rrule.Frequency = Frequency(value)
		case "INTERVAL":
			interval, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			rrule.Interval = interval
		case "COUNT":
			count, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			rrule.Count = &count
		case "UNTIL":
			until, err := icaldur.ParseIcalTime(value)
			if err != nil {
				return nil, err
			}
			rrule.Until = &until
		case "BYDAY":
			weekdays := strings.Split(value, ",")
			rrule.Weekday = make([]ByDay, 0, len(weekdays))
			for _, weekday := range weekdays {
				// if there is an interval other than 1, it can be expressed as the number at the start of the string
				interval, weekday, err := ParseByDay(weekday)
				if err != nil {
					return nil, err
				}
				rrule.Weekday = append(rrule.Weekday, ByDay{Weekday: weekday, Interval: interval})
			}
		case "BYMONTH":
			months := strings.Split(value, ",")
			rrule.Month = make([]int, 0, len(months))
			for _, month := range months {
				monthInt, err := strconv.Atoi(month)
				if err != nil {
					return nil, err
				}
				rrule.Month = append(rrule.Month, monthInt)
			}
		case "BYMONTHDAY":
			monthdays := strings.Split(value, ",")
			rrule.Monthday = make([]int, 0, len(monthdays))
			for _, monthday := range monthdays {
				monthdayInt, err := strconv.Atoi(monthday)
				if err != nil {
					return nil, err
				}
				rrule.Monthday = append(rrule.Monthday, monthdayInt)
			}
		case "BYYEARDAY":
			yeardays := strings.Split(value, ",")
			rrule.YearDay = make([]int, 0, len(yeardays))
			for _, yearday := range yeardays {
				yeardayInt, err := strconv.Atoi(yearday)
				if err != nil {
					return nil, err
				}
				rrule.YearDay = append(rrule.YearDay, yeardayInt)
			}
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

// ParseByDay parses a BYDAY value string and returns the interval and weekday.
// The string can be in the format "20MO" (interval + weekday) or just "MO" (weekday only).
// If no interval is specified, the interval defaults to 1.
// Valid weekdays are: MO, TU, WE, TH, FR, SA, SU.
// Returns (interval, weekday, error) where interval is an integer and weekday is a string.
func ParseByDay(byDayString string) (int, Weekday, error) {
	if byDayString == "" {
		return 0, "", ErrInvalidByDayString
	}

	// Check if string starts with a digit or minus sign
	if len(byDayString) > 0 && (byDayString[0] >= '0' && byDayString[0] <= '9' || byDayString[0] == '-') {
		// Find where the digits end (including negative sign)
		digitEnd := 0
		for i, char := range byDayString {
			if char < '0' || char > '9' {
				// Allow minus sign at the beginning
				if char == '-' && i == 0 {
					continue
				}
				digitEnd = i
				break
			}
			digitEnd = i + 1
		}

		// Extract interval and weekday
		intervalStr := byDayString[:digitEnd]
		weekday := Weekday(byDayString[digitEnd:])

		// Validate weekday
		if !isValidWeekday(weekday) {
			return 0, "", ErrInvalidByDayString
		}

		// Parse interval (can be negative)
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			return 0, "", ErrInvalidByDayString
		}

		return interval, weekday, nil
	}

	// No interval prefix, check if it's a valid weekday
	if !isValidWeekday(Weekday(byDayString)) {
		return 0, "", ErrInvalidByDayString
	}

	return 1, Weekday(byDayString), nil
}

// isValidWeekday checks if the string is a valid weekday abbreviation.
func isValidWeekday(weekday Weekday) bool {
	switch weekday {
	case WeekdayMonday, WeekdayTuesday, WeekdayWednesday, WeekdayThursday, WeekdayFriday, WeekdaySaturday, WeekdaySunday:
		return true
	default:
		return false
	}
}
