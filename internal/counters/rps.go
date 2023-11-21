package counters

import (
	"sync"
	"sync/atomic"
	"time"
)

type Rps struct {
	mutex         sync.Mutex
	startTime     time.Time
	totalRequests uint64
}

func NewRpsCounter() *Rps {
	return &Rps{
		startTime: time.Now(),
	}
}

func (r *Rps) Increase() {
	atomic.AddUint64(&r.totalRequests, 1)
}

func (r *Rps) GetCounter() uint64 {
	return atomic.LoadUint64(&r.totalRequests)
}

func (r *Rps) GetAndReset() float64 {
	r.mutex.Lock()
	duration := time.Since(r.startTime)
	rps := float64(r.totalRequests) / duration.Seconds()
	r.startTime = time.Now()
	r.totalRequests = 0
	r.mutex.Unlock()
	return rps
}
