package allocate

import (
        "math/big"
        "math/rand"
)
const (
        min       = 1
        max       = 200 // max possible allocations
        allocated = 1   // allocated bit mark
)
type Allocator struct {
        pool *big.Int
        last int
        low  int
        high int
}

func NewAllocator() *Allocator {
        al := &Allocator{
                pool: big.NewInt(0),
                last: min,
                low:  min,
                high: max,
        }
        return al
}

func (a *Allocator) Next() (int, bool) {
        for ; a.last <= a.high; a.last++ {
                if a.reserve(a.last) {
                        return a.last, true
                }
        }
        // See if time permits implementation of freeing bit and allocating
        // free bits for channel id
        return 0, false
}

func (a *Allocator) reserve(n int) bool {
        if a.reserved(n) {
                return false
        }
        a.pool.SetBit(a.pool, n-a.low, allocated)
        return true
}

func (a *Allocator) reserved(n int) bool {
        return a.pool.Bit(n-a.low) == allocated
}

var char = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandomID() string {
        size := 32
        id := make([]rune, size)
        for i := range id {
                id[i] = char[rand.Intn(size)]
        }
        return string(id)
}