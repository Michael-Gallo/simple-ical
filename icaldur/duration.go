// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package icaldur provides functionality to convert between iCal duration strings and golang's native time.Duration
// https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.6
package icaldur

import (
	"errors"
	"strconv"
	"time"
	"unicode"
)

var (
	ErrEmpty          = errors.New("empty duration")
	ErrBadPrefix      = errors.New("duration must start with P (optionally preceded by + or -)")
	ErrUnexpectedChar = errors.New("unexpected character")
	ErrMissingUnit    = errors.New("missing unit after number")
	ErrMixedWeeks     = errors.New("weeks form (PnW) cannot be mixed with other components")
	ErrTimeWithoutT   = errors.New("time components require a preceding 'T'")
	ErrDuplicateUnit  = errors.New("duplicate time unit")
)

func ParseICalDuration(s string) (time.Duration, error) {
	if len(s) == 0 {
		return 0, ErrEmpty
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
		return 0, ErrEmpty
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
		return 0, ErrBadPrefix
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
			return 0, ErrMissingUnit
		}
		for j := numStart; j < wpos; j++ {
			if !unicode.IsDigit(rune(s[j])) {
				return 0, ErrUnexpectedChar
			}
		}
		if wpos != len(s)-1 {
			return 0, ErrMixedWeeks
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
			return 0, ErrMissingUnit
		}
		if i >= len(s) {
			return 0, ErrMissingUnit
		}
		unit := s[i]
		i++

		switch unit {
		case 'D':
			if inTime {
				return 0, ErrUnexpectedChar
			}
			dur += v * 24 * int64(time.Hour)
		case 'H':
			if !inTime {
				return 0, ErrTimeWithoutT
			}
			if usedH {
				return 0, ErrDuplicateUnit
			}
			usedH = true
			dur += v * int64(time.Hour)
		case 'M':
			if !inTime {
				return 0, ErrTimeWithoutT
			}
			if usedM {
				return 0, ErrDuplicateUnit
			}
			usedM = true
			dur += v * int64(time.Minute)
		case 'S':
			if !inTime {
				return 0, ErrTimeWithoutT
			}
			if usedS {
				return 0, ErrDuplicateUnit
			}
			usedS = true
			dur += v * int64(time.Second)
		default:
			return 0, ErrUnexpectedChar
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
