package census

import (
	"sync"
	"sync/atomic"
)

const (
	RenewDuration          = 30
	EvictDuration          = 60
	SelfProtectedThreshold = 0.8
)

type Census struct {
	count        int64 // renew count
	needCount    int64 //   need renew count
	lastestCount int64 // latest count
	threshold    int64 //  renew least renew count
	sync.RWMutex
}

// increment renew count
func (c *Census) IncrCount() {
	atomic.AddInt64(&c.count, 1)
}

// reset  renew count
func (c *Census) ResetCount() {
	atomic.StoreInt64(&c.lastestCount, atomic.SwapInt64(&c.count, 0))
}

// increment need renew count
func (c *Census) IncrNeedCount() {
	c.Lock()
	defer c.Unlock()
	c.needCount += int64(float64(EvictDuration) / float64(RenewDuration))
	c.threshold = int64(float64(c.needCount) * SelfProtectedThreshold)
}

// decrement need renew count
func (c *Census) DecrNeedCount() {
	c.Lock()
	defer c.Unlock()
	c.needCount -= int64(float64(EvictDuration) / float64(RenewDuration))
	c.threshold = int64(float64(c.needCount) * SelfProtectedThreshold)
}

// seek need renew count
func (c *Census) SeekNeedCount(count int64) {
	c.Lock()
	defer c.Unlock()
	c.needCount = count * int64(float64(EvictDuration)/float64(RenewDuration))
	c.threshold = int64(float64(c.needCount) * SelfProtectedThreshold)
}

// check protected status
func (c *Census) ProtectedStatus() bool {

	return atomic.LoadInt64(&c.threshold) > atomic.LoadInt64(&c.lastestCount)
}
