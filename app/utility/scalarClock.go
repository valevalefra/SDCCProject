package utility

import "sync"

type Clock interface {
	Start()           // Start not concurrent operation
	Increment(id int) //id: my node identifier, do not care if ScalarClock
	Update(timestamp []uint64)
	GetValue() []uint64
}

type ScalarClock struct {
	counter uint64
	mutex   sync.Mutex
}

func (clock *ScalarClock) Start() {
	clock.counter = 0
}

func (clock *ScalarClock) Increment(_id int) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter++
}
func (clock *ScalarClock) Update(timestamp []uint64) {
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter = Max(clock.counter, timestamp[0])
}
func (clock *ScalarClock) GetValue() []uint64 {
	ret := make([]uint64, 1)
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	ret[0] = clock.counter
	return ret
}

func Max(vars ...uint64) uint64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}
