package broker

import (
        "math/big"
)
const (
        min       = 1
        max       = 200 // max possible allocations
        allocated = 1   // allocated bit mark
)
type allocator struct {
        pool *big.Int
        last int
        low  int
        high int
}

func newAllocator() *allocator {
        al := &allocator{
                pool: big.NewInt(0),
                last: min,
                low:  min,
                high: max,
        }
        return al
}

func (a *allocator) next() (int, bool) {
        for ; a.last <= a.high; a.last++ {
                if a.reserve(a.last) {
                        return a.last, true
                }
        }
        // See if time permits implementation of freeing bit and allocating
        // free bits for channel id
        return 0, false
}

func (a *allocator) reserve(n int) bool {
        if a.reserved(n) {
                return false
        }
        a.pool.SetBit(a.pool, n-a.low, allocated)
        return true
}

func (a *allocator) reserved(n int) bool {
        return a.pool.Bit(n-a.low) == allocated
}