package census

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	RenewDuration                = time.Second * 30
	ScanEvictDuration            = time.Second * 60
	SelfProtectedThreshold       = 0.8
	InstanceEvictExpiredDuration = time.Second * 90
	InstanceMaxExpiredDuration   = time.Second * 3600
)

// census
type Census struct {
	count       int64 // renew count
	needCount   int64 //   need renew count
	latestCount int64 // latest count
	threshold   int64 //  renew least renew count
	sync.RWMutex
}

// increment renew count
func (c *Census) IncrCount() {
	atomic.AddInt64(&c.count, 1)
}

// reset  renew count
func (c *Census) ResetCount() {
	atomic.StoreInt64(&c.latestCount, atomic.SwapInt64(&c.count, 0))
}

// increment need renew count
func (c *Census) IncrNeedCount() {
	c.Lock()
	defer c.Unlock()
	c.needCount += int64(float64(ScanEvictDuration) / float64(RenewDuration))
	c.threshold = int64(float64(c.needCount) * SelfProtectedThreshold)
}

// decrement need renew count
func (c *Census) DecrNeedCount() {
	c.Lock()
	defer c.Unlock()
	c.needCount -= int64(float64(ScanEvictDuration) / float64(RenewDuration))
	c.threshold = int64(float64(c.needCount) * SelfProtectedThreshold)
}

// seek need renew count
func (c *Census) SeekNeedCount(count int64) {
	c.Lock()
	defer c.Unlock()
	c.needCount = count * int64(float64(ScanEvictDuration)/float64(RenewDuration))
	c.threshold = int64(float64(c.needCount) * SelfProtectedThreshold)
}

// check protected status
func (c *Census) ProtectedStatus() bool {

	return atomic.LoadInt64(&c.threshold) > atomic.LoadInt64(&c.latestCount)
}
