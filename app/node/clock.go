package main

import "SDCCProject/app/utility"

var (
	scalarClock utility.ScalarClock
)

func startClocks() {
	scalarClock.Start()
}

func incrementClock(clock utility.Clock) {
	clock.Increment()
}

func updateClock(clock utility.Clock, timestamp int) {
	clock.Update(timestamp)
}

func getValueClock(clock utility.Clock) []int {
	return clock.GetValue()
}
