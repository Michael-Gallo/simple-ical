// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rrule

import "errors"

// Predefined errors for the rrule package.
var (
	// ErrInvalidRRuleString is returned when the rrule string format is invalid.
	ErrInvalidRRuleString = errors.New("invalid rrule string")

	// ErrFrequencyRequired is returned when the frequency property is missing.
	ErrFrequencyRequired = errors.New("frequency is required")

	// ErrCountAndUntilBothSet is returned when both count and until properties are set.
	ErrCountAndUntilBothSet = errors.New("count and until cannot both be set")

	// ErrInvalidInterval is returned when the interval is not a positive integer.
	ErrInvalidInterval = errors.New("interval must be a positive integer")

	// ErrInvalidByDayString is returned when the BYDAY string format is invalid.
	ErrInvalidByDayString = errors.New("invalid BYDAY string")

	errInvalidFrequency = errors.New("invalid frequency")
)
