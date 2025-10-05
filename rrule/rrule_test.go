package rrule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: replace with calls to New once go 1.26 is released
func getPointer[T any](v T) *T {
	return &v
}

func TestParseRRule(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *RRule
		expectError error
	}{

		{
			name:  "Valid daily rule",
			input: "FREQ=DAILY;INTERVAL=1;COUNT=10",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  1,
				Count:     getPointer(10),
				Until:     nil,
			},
			expectError: nil,
		},
		{
			name:  "Valid daily rule with interval set",
			input: "FREQ=DAILY;INTERVAL=2;COUNT=10",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  2,
				Count:     getPointer(10),
				Until:     nil,
			},
			expectError: nil,
		},
		{
			name:  "Valid daily rule with interval not set",
			input: "FREQ=DAILY;COUNT=10",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  1,
				Count:     getPointer(10),
				Until:     nil,
			},
			expectError: nil,
		},
		{
			name:        "Invalid rule: missing frequency",
			input:       "INTERVAL=1;COUNT=10",
			want:        nil,
			expectError: ErrFrequencyRequired,
		},
		{
			name:        "Invalid rule: count and until cannot both be set",
			input:       "FREQ=DAILY;COUNT=10;UNTIL=19730429T070000Z",
			want:        nil,
			expectError: ErrCountAndUntilBothSet,
		},
		{
			name:        "Invalid rule: interval must be a positive integer",
			input:       "FREQ=DAILY;INTERVAL=0;COUNT=10",
			want:        nil,
			expectError: ErrInvalidInterval,
		},
		{
			name:        "Invalid rule: malformed rrule string",
			input:       "FREQ=DAILY;INVALID",
			want:        nil,
			expectError: ErrInvalidRRuleString,
		},
		{
			name:  "10th of June and July every year",
			input: "FREQ=YEARLY;COUNT=10;BYMONTH=6,7",
			want: &RRule{
				Frequency: FrequencyYearly,
				Interval:  1,
				Count:     getPointer(10),
				Month:     []int{6, 7},
			},
			expectError: nil,
		},
		{
			name:  "Monthly on the third-to-the-last day of the month, forever",
			input: "FREQ=MONTHLY;BYMONTHDAY=-3",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Monthday:  []int{-3},
			},
			expectError: nil,
		},
		{
			name:  "Monthly on the first and last day of the month for 10 occurrences",
			input: "FREQ=MONTHLY;COUNT=10;BYMONTHDAY=1,-1",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Count:     getPointer(10),
				Monthday:  []int{1, -1},
			},
			expectError: nil,
		},
		{
			name:  "Every Tuesday, every other month",
			input: "FREQ=MONTHLY;INTERVAL=2;BYDAY=TU",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  2,
				Weekday: []ByDay{{
					Weekday:  WeekdayTuesday,
					Interval: 1,
				}},
			},
			expectError: nil,
		},
		{
			name:  "Every Tuesday, every other month",
			input: "FREQ=MONTHLY;INTERVAL=2;BYDAY=TU",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  2,
				Weekday: []ByDay{{
					Weekday:  WeekdayTuesday,
					Interval: 1,
				}},
			},
			expectError: nil,
		},
		{
			name:  "Every third year on the 1st, 100th, and 200th day for 10 occurrences:",
			input: "FREQ=YEARLY;INTERVAL=3;COUNT=10;BYYEARDAY=1,100,200",
			want: &RRule{
				Frequency: FrequencyYearly,
				Interval:  3,
				Count:     getPointer(10),
				YearDay:   []int{1, 100, 200},
			},
			expectError: nil,
		},
		{
			name:  "Every 20th Monday of the year, forever",
			input: "FREQ=YEARLY;BYDAY=20MO",
			want: &RRule{
				Frequency: FrequencyYearly,
				Interval:  1,
				Weekday:   []ByDay{{Weekday: WeekdayMonday, Interval: 20}},
			},
			expectError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule, err := ParseRRule(test.input)
			if test.expectError != nil {
				assert.ErrorIs(t, err, test.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.want, rule)
		})
	}
}

func TestParseByDay(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedInt     int
		expectedWeekDay Weekday
		expectError     error
	}{
		{
			name:            "String with interval and weekday",
			input:           "20MO",
			expectedInt:     20,
			expectedWeekDay: WeekdayMonday,
			expectError:     nil,
		},
		{
			name:            "String with just weekday",
			input:           "MO",
			expectedInt:     1,
			expectedWeekDay: WeekdayMonday,
			expectError:     nil,
		},
		{
			name:            "String with interval and Tuesday",
			input:           "5TU",
			expectedInt:     5,
			expectedWeekDay: WeekdayTuesday,
			expectError:     nil,
		},
		{
			name:            "String with just Tuesday",
			input:           "TU",
			expectedInt:     1,
			expectedWeekDay: WeekdayTuesday,
			expectError:     nil,
		},
		{
			name:            "String with interval and Wednesday",
			input:           "3WE",
			expectedInt:     3,
			expectedWeekDay: WeekdayWednesday,
			expectError:     nil,
		},
		{
			name:            "String with just Wednesday",
			input:           "WE",
			expectedInt:     1,
			expectedWeekDay: WeekdayWednesday,
			expectError:     nil,
		},
		{
			name:            "String with interval and Thursday",
			input:           "7TH",
			expectedInt:     7,
			expectedWeekDay: WeekdayThursday,
			expectError:     nil,
		},
		{
			name:            "String with just Thursday",
			input:           "TH",
			expectedInt:     1,
			expectedWeekDay: WeekdayThursday,
			expectError:     nil,
		},
		{
			name:            "String with interval and Friday",
			input:           "2FR",
			expectedInt:     2,
			expectedWeekDay: WeekdayFriday,
			expectError:     nil,
		},
		{
			name:            "String with just Friday",
			input:           "FR",
			expectedInt:     1,
			expectedWeekDay: WeekdayFriday,
			expectError:     nil,
		},
		{
			name:            "String with interval and Saturday",
			input:           "4SA",
			expectedInt:     4,
			expectedWeekDay: WeekdaySaturday,
			expectError:     nil,
		},
		{
			name:            "String with just Saturday",
			input:           "SA",
			expectedInt:     1,
			expectedWeekDay: WeekdaySaturday,
			expectError:     nil,
		},
		{
			name:            "String with interval and Sunday",
			input:           "6SU",
			expectedInt:     6,
			expectedWeekDay: WeekdaySunday,
			expectError:     nil,
		},
		{
			name:            "String with just Sunday",
			input:           "SU",
			expectedInt:     1,
			expectedWeekDay: WeekdaySunday,
			expectError:     nil,
		},
		{
			name:        "Invalid string returns error",
			input:       "INVALID",
			expectedInt: 0,
			expectError: ErrInvalidByDayString,
		},
		{
			name:        "Empty string returns error",
			input:       "",
			expectedInt: 0,
			expectError: ErrInvalidByDayString,
		},
		{
			name:            "String with invalid weekday returns error",
			input:           "5XX",
			expectedInt:     0,
			expectedWeekDay: "",
			expectError:     ErrInvalidByDayString,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			interval, weekday, err := ParseByDay(test.input)
			if test.expectError != nil {
				assert.ErrorIs(t, err, test.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedInt, interval)
			assert.Equal(t, test.expectedWeekDay, weekday)
		})
	}
}
