package counters

import (
	"sync/atomic"
)

type InFlight struct {
	count *int64
}

func NewInFlightCounter() *InFlight {
	return &InFlight{
		count: new(int64),
	}
}

func (i *InFlight) Increase() {
	atomic.AddInt64(i.count, 1)
}

func (i *InFlight) Decrease() {
	atomic.AddInt64(i.count, -1)
}

func (i *InFlight) Get() int64 {
	return atomic.LoadInt64(i.count)
}
