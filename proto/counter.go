package proto

import (
	"math/rand"
	"sync/atomic"
	"time"
)

var counter int64

func init() {
	rand.Seed(time.Now().UnixNano())

	counter = time.Now().UnixNano()
}

// NextCnt returns a incremented atomic counter
func NextCnt() int64 {
	return atomic.AddInt64(&counter, 1)
}
