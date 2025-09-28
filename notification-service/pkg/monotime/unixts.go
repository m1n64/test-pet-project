package monotime

import (
	"sync/atomic"
	"time"
)

var lastTs int64
var seq uint32

func NowNanoUnique() int64 {
	ts := time.Now().UnixNano()
	prev := atomic.LoadInt64(&lastTs)

	if ts == prev {
		n := atomic.AddUint32(&seq, 1)
		return ts + int64(n)
	}

	atomic.StoreInt64(&lastTs, ts)
	atomic.StoreUint32(&seq, 0)
	return ts
}
