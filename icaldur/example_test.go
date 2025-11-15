package icaldur_test

import (
	"fmt"

	"github.com/michael-gallo/simpleical/icaldur"
)

func ExampleParseICalDuration() {
	duration, err := icaldur.ParseICalDuration("P15DT5H0M20S")
	if err != nil {
		panic(err)
	}
	fmt.Println(duration.String())
	// Output: 365h0m20s
}

func ExampleParseIcalTime() {
	time, err := icaldur.ParseIcalTime("20250928T183000Z")
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Year())
	fmt.Println(time.Month())
	fmt.Println(time.Day())
	fmt.Println(time.Hour())
	fmt.Println(time.Minute())
	fmt.Println(time.Second())
	// Output: 2025
	// September
	// 28
	// 18
	// 30
	// 0
}
