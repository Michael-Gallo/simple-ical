// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rrule

import "errors"

// Predefined errors for the rrule package.
var (
	// errInvalidRRuleString is returned when the rrule string format is invalid.
	errInvalidRRuleString = errors.New("invalid rrule string")

	// errFrequencyRequired is returned when the frequency property is missing.
	errFrequencyRequired = errors.New("frequency is required")

	// errCountAndUntilBothSet is returned when both count and until properties are set.
	errCountAndUntilBothSet = errors.New("count and until cannot both be set")

	// errInvalidInterval is returned when the interval is not a positive integer.
	errInvalidInterval = errors.New("interval must be a positive integer")

	// errInvalidByDayString is returned when the BYDAY string format is invalid.
	errInvalidByDayString = errors.New("invalid BYDAY string")

	errInvalidFrequency = errors.New("invalid frequency")
)
