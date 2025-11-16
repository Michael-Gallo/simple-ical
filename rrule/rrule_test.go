package rrule

import (
	"fmt"
	"testing"
	"time"

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
			name:        "Invalid frequency",
			input:       "FREQ=DALLY;INTERVAL=2;COUNT=10",
			want:        nil,
			expectError: fmt.Errorf("%w: %s", errInvalidFrequency, "DALLY"),
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
			expectError: errFrequencyRequired,
		},
		{
			name:        "Invalid rule: count and until cannot both be set",
			input:       "FREQ=DAILY;COUNT=10;UNTIL=19730429T070000Z",
			want:        nil,
			expectError: errCountAndUntilBothSet,
		},
		{
			name:        "Invalid rule: interval must be a positive integer",
			input:       "FREQ=DAILY;INTERVAL=0;COUNT=10",
			want:        nil,
			expectError: errInvalidInterval,
		},
		{
			name:        "Invalid rule: malformed rrule string",
			input:       "FREQ=DAILY;INVALID",
			want:        nil,
			expectError: errInvalidRRuleString,
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
		// DAILY examples from RFC 5545
		{
			name:  "Daily for 10 occurrences",
			input: "FREQ=DAILY;COUNT=10",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  1,
				Count:     getPointer(10),
			},
			expectError: nil,
		},
		{
			name:  "Daily until December 24, 1997",
			input: "FREQ=DAILY;UNTIL=19971224T000000Z",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  1,
				Until:     getPointer(time.Date(1997, 12, 24, 0, 0, 0, 0, time.UTC)),
			},
			expectError: nil,
		},
		{
			name:  "Every other day - forever",
			input: "FREQ=DAILY;INTERVAL=2",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  2,
			},
			expectError: nil,
		},
		{
			name:  "Every 10 days, 5 occurrences",
			input: "FREQ=DAILY;INTERVAL=10;COUNT=5",
			want: &RRule{
				Frequency: FrequencyDaily,
				Interval:  10,
				Count:     getPointer(5),
			},
			expectError: nil,
		},
		// WEEKLY examples from RFC 5545
		{
			name:  "Weekly for 10 occurrences",
			input: "FREQ=WEEKLY;COUNT=10",
			want: &RRule{
				Frequency: FrequencyWeekly,
				Interval:  1,
				Count:     getPointer(10),
			},
			expectError: nil,
		},
		{
			name:  "Weekly until December 24, 1997",
			input: "FREQ=WEEKLY;UNTIL=19971224T000000Z",
			want: &RRule{
				Frequency: FrequencyWeekly,
				Interval:  1,
				Until:     getPointer(time.Date(1997, 12, 24, 0, 0, 0, 0, time.UTC)),
			},
			expectError: nil,
		},
		{
			name:  "Every other week - forever",
			input: "FREQ=WEEKLY;INTERVAL=2",
			want: &RRule{
				Frequency: FrequencyWeekly,
				Interval:  2,
			},
			expectError: nil,
		},
		{
			name:  "Weekly on Tuesday and Thursday for five weeks",
			input: "FREQ=WEEKLY;COUNT=10;BYDAY=TU,TH",
			want: &RRule{
				Frequency: FrequencyWeekly,
				Interval:  1,
				Count:     getPointer(10),
				Weekday: []ByDay{
					{Weekday: WeekdayTuesday, Interval: 1},
					{Weekday: WeekdayThursday, Interval: 1},
				},
			},
			expectError: nil,
		},
		{
			name:  "Every other week on Monday, Wednesday, and Friday until December 24, 1997",
			input: "FREQ=WEEKLY;INTERVAL=2;UNTIL=19971224T000000Z;BYDAY=MO,WE,FR",
			want: &RRule{
				Frequency: FrequencyWeekly,
				Interval:  2,
				Until:     getPointer(time.Date(1997, 12, 24, 0, 0, 0, 0, time.UTC)),
				Weekday: []ByDay{
					{Weekday: WeekdayMonday, Interval: 1},
					{Weekday: WeekdayWednesday, Interval: 1},
					{Weekday: WeekdayFriday, Interval: 1},
				},
			},
			expectError: nil,
		},
		{
			name:  "Every other week on Tuesday and Thursday, for 8 occurrences",
			input: "FREQ=WEEKLY;INTERVAL=2;COUNT=8;BYDAY=TU,TH",
			want: &RRule{
				Frequency: FrequencyWeekly,
				Interval:  2,
				Count:     getPointer(8),
				Weekday: []ByDay{
					{Weekday: WeekdayTuesday, Interval: 1},
					{Weekday: WeekdayThursday, Interval: 1},
				},
			},
			expectError: nil,
		},
		// MONTHLY examples from RFC 5545
		{
			name:  "Monthly on the first Friday for 10 occurrences",
			input: "FREQ=MONTHLY;COUNT=10;BYDAY=1FR",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Count:     getPointer(10),
				Weekday:   []ByDay{{Weekday: WeekdayFriday, Interval: 1}},
			},
			expectError: nil,
		},
		{
			name:  "Monthly on the first Friday until December 24, 1997",
			input: "FREQ=MONTHLY;UNTIL=19971224T000000Z;BYDAY=1FR",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Until:     getPointer(time.Date(1997, 12, 24, 0, 0, 0, 0, time.UTC)),
				Weekday:   []ByDay{{Weekday: WeekdayFriday, Interval: 1}},
			},
			expectError: nil,
		},
		{
			name:  "Every other month on the first and last Sunday of the month for 10 occurrences",
			input: "FREQ=MONTHLY;INTERVAL=2;COUNT=10;BYDAY=1SU,-1SU",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  2,
				Count:     getPointer(10),
				Weekday: []ByDay{
					{Weekday: WeekdaySunday, Interval: 1},
					{Weekday: WeekdaySunday, Interval: -1},
				},
			},
			expectError: nil,
		},
		{
			name:  "Monthly on the second-to-last Monday of the month for 6 months",
			input: "FREQ=MONTHLY;COUNT=6;BYDAY=-2MO",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Count:     getPointer(6),
				Weekday:   []ByDay{{Weekday: WeekdayMonday, Interval: -2}},
			},
			expectError: nil,
		},
		{
			name:  "Monthly on the 2nd and 15th of the month for 10 occurrences",
			input: "FREQ=MONTHLY;COUNT=10;BYMONTHDAY=2,15",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Count:     getPointer(10),
				Monthday:  []int{2, 15},
			},
			expectError: nil,
		},
		{
			name:  "Every 18 months on the 10th thru 15th of the month for 10 occurrences",
			input: "FREQ=MONTHLY;INTERVAL=18;COUNT=10;BYMONTHDAY=10,11,12,13,14,15",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  18,
				Count:     getPointer(10),
				Monthday:  []int{10, 11, 12, 13, 14, 15},
			},
			expectError: nil,
		},
		// YEARLY examples from RFC 5545
		{
			name:  "Yearly in June and July for 10 occurrences",
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
			name:  "Every other year on January, February, and March for 10 occurrences",
			input: "FREQ=YEARLY;INTERVAL=2;COUNT=10;BYMONTH=1,2,3",
			want: &RRule{
				Frequency: FrequencyYearly,
				Interval:  2,
				Count:     getPointer(10),
				Month:     []int{1, 2, 3},
			},
			expectError: nil,
		},
		{
			name:  "Every Thursday in March, forever",
			input: "FREQ=YEARLY;BYMONTH=3;BYDAY=TH",
			want: &RRule{
				Frequency: FrequencyYearly,
				Interval:  1,
				Month:     []int{3},
				Weekday:   []ByDay{{Weekday: WeekdayThursday, Interval: 1}},
			},
			expectError: nil,
		},
		{
			name:  "Every Thursday, but only during June, July, and August, forever",
			input: "FREQ=YEARLY;BYDAY=TH;BYMONTH=6,7,8",
			want: &RRule{
				Frequency: FrequencyYearly,
				Interval:  1,
				Month:     []int{6, 7, 8},
				Weekday:   []ByDay{{Weekday: WeekdayThursday, Interval: 1}},
			},
			expectError: nil,
		},
		{
			name:  "Every Friday the 13th, forever",
			input: "FREQ=MONTHLY;BYDAY=FR;BYMONTHDAY=13",
			want: &RRule{
				Frequency: FrequencyMonthly,
				Interval:  1,
				Weekday:   []ByDay{{Weekday: WeekdayFriday, Interval: 1}},
				Monthday:  []int{13},
			},
			expectError: nil,
		},
		// HOURLY and MINUTELY examples from RFC 5545
		{
			name:  "Every 3 hours from 9:00 AM to 5:00 PM on a specific day",
			input: "FREQ=HOURLY;INTERVAL=3;UNTIL=19970902T170000Z",
			want: &RRule{
				Frequency: FrequencyHourly,
				Interval:  3,
				Until:     getPointer(time.Date(1997, 9, 2, 17, 0, 0, 0, time.UTC)),
			},
			expectError: nil,
		},
		{
			name:  "Every 15 minutes for 6 occurrences",
			input: "FREQ=MINUTELY;INTERVAL=15;COUNT=6",
			want: &RRule{
				Frequency: FrequencyMinutely,
				Interval:  15,
				Count:     getPointer(6),
			},
			expectError: nil,
		},
		{
			name:  "Every hour and a half for 4 occurrences",
			input: "FREQ=MINUTELY;INTERVAL=90;COUNT=4",
			want: &RRule{
				Frequency: FrequencyMinutely,
				Interval:  90,
				Count:     getPointer(4),
			},
			expectError: nil,
		},
		// Missing RFC 5545 examples that need to be implemented
		// TODO: Uncomment when WKST property is implemented
		// {
		// 	name:  "Every other week - forever with Sunday as week start",
		// 	input: "FREQ=WEEKLY;INTERVAL=2;WKST=SU",
		// 	want: &RRule{
		// 		Frequency: FrequencyWeekly,
		// 		Interval:  2,
		// 		WeekStart: WeekdaySunday,
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "Weekly on Tuesday and Thursday for five weeks with Sunday as week start",
		// 	input: "FREQ=WEEKLY;UNTIL=19971007T000000Z;WKST=SU;BYDAY=TU,TH",
		// 	want: &RRule{
		// 		Frequency: FrequencyWeekly,
		// 		Interval:  1,
		// 		Until:     getPointer(time.Date(1997, 10, 7, 0, 0, 0, 0, time.UTC)),
		// 		WeekStart: WeekdaySunday,
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayTuesday, Interval: 1},
		// 			{Weekday: WeekdayThursday, Interval: 1},
		// 		},
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "Every other week on Monday, Wednesday, and Friday until December 24, 1997 with Sunday as week start",
		// 	input: "FREQ=WEEKLY;INTERVAL=2;UNTIL=19971224T000000Z;WKST=SU;BYDAY=MO,WE,FR",
		// 	want: &RRule{
		// 		Frequency: FrequencyWeekly,
		// 		Interval:  2,
		// 		Until:     getPointer(time.Date(1997, 12, 24, 0, 0, 0, 0, time.UTC)),
		// 		WeekStart: WeekdaySunday,
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayMonday, Interval: 1},
		// 			{Weekday: WeekdayWednesday, Interval: 1},
		// 			{Weekday: WeekdayFriday, Interval: 1},
		// 		},
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "Every other week on Tuesday and Thursday, for 8 occurrences with Sunday as week start",
		// 	input: "FREQ=WEEKLY;INTERVAL=2;COUNT=8;WKST=SU;BYDAY=TU,TH",
		// 	want: &RRule{
		// 		Frequency: FrequencyWeekly,
		// 		Interval:  2,
		// 		Count:     getPointer(8),
		// 		WeekStart: WeekdaySunday,
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayTuesday, Interval: 1},
		// 			{Weekday: WeekdayThursday, Interval: 1},
		// 		},
		// 	},
		// 	expectError: nil,
		// },

		// TODO: Uncomment when BYWEEKNO property is implemented
		// {
		// 	name:  "Monday of week number 20 (where the default start of the week is Monday), forever",
		// 	input: "FREQ=YEARLY;BYWEEKNO=20;BYDAY=MO",
		// 	want: &RRule{
		// 		Frequency: FrequencyYearly,
		// 		Interval:  1,
		// 		WeekNo:    []int{20},
		// 		Weekday:   []ByDay{{Weekday: WeekdayMonday, Interval: 1}},
		// 	},
		// 	expectError: nil,
		// },

		// TODO: Uncomment when BYSETPOS property is implemented
		// {
		// 	name:  "The third instance into the month of one of Tuesday, Wednesday, or Thursday, for the next 3 months",
		// 	input: "FREQ=MONTHLY;COUNT=3;BYDAY=TU,WE,TH;BYSETPOS=3",
		// 	want: &RRule{
		// 		Frequency: FrequencyMonthly,
		// 		Interval:  1,
		// 		Count:     getPointer(3),
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayTuesday, Interval: 1},
		// 			{Weekday: WeekdayWednesday, Interval: 1},
		// 			{Weekday: WeekdayThursday, Interval: 1},
		// 		},
		// 		SetPos: []int{3},
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "The second-to-last weekday of the month",
		// 	input: "FREQ=MONTHLY;BYDAY=MO,TU,WE,TH,FR;BYSETPOS=-2",
		// 	want: &RRule{
		// 		Frequency: FrequencyMonthly,
		// 		Interval:  1,
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayMonday, Interval: 1},
		// 			{Weekday: WeekdayTuesday, Interval: 1},
		// 			{Weekday: WeekdayWednesday, Interval: 1},
		// 			{Weekday: WeekdayThursday, Interval: 1},
		// 			{Weekday: WeekdayFriday, Interval: 1},
		// 		},
		// 		SetPos: []int{-2},
		// 	},
		// 	expectError: nil,
		// },

		// TODO: Uncomment when complex combinations with multiple BY* properties are implemented
		// {
		// 	name:  "Every 4 years, the first Tuesday after a Monday in November, forever (U.S. Presidential Election day)",
		// 	input: "FREQ=YEARLY;INTERVAL=4;BYMONTH=11;BYDAY=TU;BYMONTHDAY=2,3,4,5,6,7,8",
		// 	want: &RRule{
		// 		Frequency: FrequencyYearly,
		// 		Interval:  4,
		// 		Month:     []int{11},
		// 		Weekday:   []ByDay{{Weekday: WeekdayTuesday, Interval: 1}},
		// 		Monthday:  []int{2, 3, 4, 5, 6, 7, 8},
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "The first Saturday that follows the first Sunday of the month, forever",
		// 	input: "FREQ=MONTHLY;BYDAY=SA;BYMONTHDAY=7,8,9,10,11,12,13",
		// 	want: &RRule{
		// 		Frequency: FrequencyMonthly,
		// 		Interval:  1,
		// 		Weekday:   []ByDay{{Weekday: WeekdaySaturday, Interval: 1}},
		// 		Monthday:  []int{7, 8, 9, 10, 11, 12, 13},
		// 	},
		// 	expectError: nil,
		// },

		// TODO: Uncomment when BYHOUR and BYMINUTE properties are implemented
		// {
		// 	name:  "Every 20 minutes from 9:00 AM to 4:40 PM every day",
		// 	input: "FREQ=DAILY;BYHOUR=9,10,11,12,13,14,15,16;BYMINUTE=0,20,40",
		// 	want: &RRule{
		// 		Frequency: FrequencyDaily,
		// 		Interval:  1,
		// 		Hour:      []int{9, 10, 11, 12, 13, 14, 15, 16},
		// 		Minute:    []int{0, 20, 40},
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "Every 20 minutes from 9:00 AM to 4:40 PM every day (alternative with MINUTELY)",
		// 	input: "FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16",
		// 	want: &RRule{
		// 		Frequency: FrequencyMinutely,
		// 		Interval:  20,
		// 		Hour:      []int{9, 10, 11, 12, 13, 14, 15, 16},
		// 	},
		// 	expectError: nil,
		// },

		// TODO: Uncomment when WKST property is implemented
		// {
		// 	name:  "An example where the days generated makes a difference because of WKST (Monday start)",
		// 	input: "FREQ=WEEKLY;INTERVAL=2;COUNT=4;BYDAY=TU,SU;WKST=MO",
		// 	want: &RRule{
		// 		Frequency: FrequencyWeekly,
		// 		Interval:  2,
		// 		Count:     getPointer(4),
		// 		WeekStart: WeekdayMonday,
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayTuesday, Interval: 1},
		// 			{Weekday: WeekdaySunday, Interval: 1},
		// 		},
		// 	},
		// 	expectError: nil,
		// },
		// {
		// 	name:  "An example where the days generated makes a difference because of WKST (Sunday start)",
		// 	input: "FREQ=WEEKLY;INTERVAL=2;COUNT=4;BYDAY=TU,SU;WKST=SU",
		// 	want: &RRule{
		// 		Frequency: FrequencyWeekly,
		// 		Interval:  2,
		// 		Count:     getPointer(4),
		// 		WeekStart: WeekdaySunday,
		// 		Weekday: []ByDay{
		// 			{Weekday: WeekdayTuesday, Interval: 1},
		// 			{Weekday: WeekdaySunday, Interval: 1},
		// 		},
		// 	},
		// 	expectError: nil,
		// },

		// TODO: Uncomment when complex validation is implemented
		// {
		// 	name:  "An example where an invalid date (i.e., February 30) is ignored",
		// 	input: "FREQ=MONTHLY;BYMONTHDAY=15,30;COUNT=5",
		// 	want: &RRule{
		// 		Frequency: FrequencyMonthly,
		// 		Interval:  1,
		// 		Count:     getPointer(5),
		// 		Monthday:  []int{15, 30},
		// 	},
		// 	expectError: nil,
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule, err := ParseRRule(test.input)
			if test.expectError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.expectError.Error())
				assert.Nil(t, rule)
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
			expectError: errInvalidByDayString,
		},
		{
			name:        "Empty string returns error",
			input:       "",
			expectedInt: 0,
			expectError: errInvalidByDayString,
		},
		{
			name:            "String with invalid weekday returns error",
			input:           "5XX",
			expectedInt:     0,
			expectedWeekDay: "",
			expectError:     errInvalidByDayString,
		},
		{
			name:            "String with negative interval and weekday",
			input:           "-1SU",
			expectedInt:     -1,
			expectedWeekDay: WeekdaySunday,
			expectError:     nil,
		},
		{
			name:            "String with negative interval and Monday",
			input:           "-2MO",
			expectedInt:     -2,
			expectedWeekDay: WeekdayMonday,
			expectError:     nil,
		},
		{
			name:            "String with negative interval and Friday",
			input:           "-3FR",
			expectedInt:     -3,
			expectedWeekDay: WeekdayFriday,
			expectError:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			interval, weekday, err := parseByDay(test.input)
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
