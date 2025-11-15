package rrule_test

import (
	"fmt"

	"github.com/michael-gallo/simpleical/rrule"
)

func ExampleParseRRule() {
	rrule, err := rrule.ParseRRule("FREQ=DAILY;INTERVAL=1;COUNT=10")
	if err != nil {
		panic(err)
	}
	fmt.Println(rrule.Frequency)
	fmt.Println(rrule.Interval)
	fmt.Println(*rrule.Count)
	// Output: DAILY
	// 1
	// 10
}
