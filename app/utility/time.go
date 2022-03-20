package utility

import "sync"

type Clock interface {
	Start()     // Start not concurrent operation
	Increment() //id: my node identifier, do not care if ScalarClock
	Update(timestamp int)
	GetValue() []int
}

type ScalarClock struct {
	counter int
	mutex   sync.Mutex
}

func (clock *ScalarClock) Start() {
	clock.counter = 0
}

func (clock *ScalarClock) Increment() {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter++
}
func (clock *ScalarClock) Update(timestamp int) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter = Max(clock.counter, timestamp)
}
func (clock *ScalarClock) GetValue() []int {
	ret := make([]int, 1)
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	ret[0] = clock.counter
	return ret
}

func Max(vars ...int) int {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}
