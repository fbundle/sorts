package el_sorts

import (
	"sync/atomic"
)

func not(b bool) bool {
	return !b
}

var valueCount uint64

func nextValue() uint64 {
	return atomic.AddUint64(&valueCount, 1)
}
