package el_sorts

import (
	"fmt"
	"sync/atomic"

	"github.com/fbundle/sorts/form"
)

var TypeErr = fmt.Errorf("type_err")

func mustMatchHead(Head form.Name, list form.List) {
	if len(list) == 0 || Head != list[0] {
		panic(TypeErr)
	}
}

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
