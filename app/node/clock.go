package main

import "SDCCProject/app/utility"

var (
	scalarClock utility.ScalarClock
)

func startClocks() {
	scalarClock.Start()
}

func incrementClock(clock utility.Clock, id int) {
	clock.Increment(id - 1)
}

func updateClock(clock utility.Clock, timestamp []uint64) {
	clock.Update(timestamp)
}

func getValueClock(clock utility.Clock) []uint64 {
	return clock.GetValue()
}
