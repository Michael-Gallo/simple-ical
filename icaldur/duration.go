// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package icaldur

import (
	"errors"
	"strconv"
	"time"
	"unicode"
)

var (
	errEmpty          = errors.New("empty duration")
	errBadPrefix      = errors.New("duration must start with P (optionally preceded by + or -)")
	errUnexpectedChar = errors.New("unexpected character")
	errMissingUnit    = errors.New("missing unit after number")
	errMixedWeeks     = errors.New("weeks form (PnW) cannot be mixed with other components")
	errTimeWithoutT   = errors.New("time components require a preceding 'T'")
	errDuplicateUnit  = errors.New("duplicate time unit")
)

// ParseICalDuration parses an iCal duration string according to RFC 5545 section 3.3.6 into a time.Duration.
// The string can be prefixed with a + or - sign to indicate a positive or negative duration
// The string can contain the following units:
// - D: days
// - H: hours
// - M: minutes
// - S: seconds
// - W: weeks.
func ParseICalDuration(s string) (time.Duration, error) {
	if len(s) == 0 {
		return 0, errEmpty
	}

	// Trim spaces (optional)
	start, end := 0, len(s)
	for start < end && unicode.IsSpace(rune(s[start])) {
		start++
	}
	for end > start && unicode.IsSpace(rune(s[end-1])) {
		end--
	}
	if start == end {
		return 0, errEmpty
	}
	s = s[start:end]

	sign := int64(1)
	i := 0

	// Optional sign
	switch s[i] {
	case '+':
		i++
	case '-':
		sign = -1
		i++
	}

	// Must start with 'P'
	if i >= len(s) || s[i] != 'P' {
		return 0, errBadPrefix
	}
	i++

	var (
		inTime              bool
		dur                 int64 // nanoseconds
		usedH, usedM, usedS bool
	)

	// Helper to read a positive integer
	readInt := func() (int64, bool) {
		if i >= len(s) || !unicode.IsDigit(rune(s[i])) {
			return 0, false
		}
		start := i
		for i < len(s) && unicode.IsDigit(rune(s[i])) {
			i++
		}
		v, err := strconv.ParseInt(s[start:i], 10, 64)
		if err != nil {
			return 0, false
		}
		return v, true
	}

	// Special-case weeks: PnW and nothing else
	// Detect if there's a 'W' anywhere; if present, it must be the only unit
	if wpos := indexByteFrom(s, 'W', i); wpos != -1 {
		// Ensure there are only digits between i and wpos, and nothing after W
		numStart := i
		if numStart >= wpos {
			return 0, errMissingUnit
		}
		for j := numStart; j < wpos; j++ {
			if !unicode.IsDigit(rune(s[j])) {
				return 0, errUnexpectedChar
			}
		}
		if wpos != len(s)-1 {
			return 0, errMixedWeeks
		}
		v, err := strconv.ParseInt(s[numStart:wpos], 10, 64)
		if err != nil {
			return 0, err
		}
		dur = v * 7 * 24 * int64(time.Hour)
		return time.Duration(sign * dur), nil
	}

	// Otherwise parse date/time components: P[nD][T[nH][nM][nS]]
	for i < len(s) {
		if s[i] == 'T' {
			inTime = true
			i++
			continue
		}

		v, ok := readInt()
		if !ok {
			return 0, errMissingUnit
		}
		if i >= len(s) {
			return 0, errMissingUnit
		}
		unit := s[i]
		i++

		switch unit {
		case 'D':
			if inTime {
				return 0, errUnexpectedChar
			}
			dur += v * 24 * int64(time.Hour)
		case 'H':
			if !inTime {
				return 0, errTimeWithoutT
			}
			if usedH {
				return 0, errDuplicateUnit
			}
			usedH = true
			dur += v * int64(time.Hour)
		case 'M':
			if !inTime {
				return 0, errTimeWithoutT
			}
			if usedM {
				return 0, errDuplicateUnit
			}
			usedM = true
			dur += v * int64(time.Minute)
		case 'S':
			if !inTime {
				return 0, errTimeWithoutT
			}
			if usedS {
				return 0, errDuplicateUnit
			}
			usedS = true
			dur += v * int64(time.Second)
		default:
			return 0, errUnexpectedChar
		}
	}

	return time.Duration(sign * dur), nil
}

// indexByteFrom finds the first index of b in s starting at from, or -1.
func indexByteFrom(s string, b byte, from int) int {
	for j := from; j < len(s); j++ {
		if s[j] == b {
			return j
		}
	}
	return -1
}
