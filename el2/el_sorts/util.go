package el_sorts

import (
	"sync/atomic"
)

func not(b bool) bool {
	return !b
}

func mustTermOf(ctx Context, x Sort, X Sort) {
	if !ctx.LessEqual(ctx.Parent(x), X) {
		panic(TypeErr)
	}
}

var valueCount uint64

func nextValue() uint64 {
	return atomic.AddUint64(&valueCount, 1)
}
