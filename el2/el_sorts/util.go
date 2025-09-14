package almost_sort_extra

import (
	"fmt"

	"github.com/fbundle/sorts/form"
)

var TypeErr = fmt.Errorf("type_err")

func mustMatchHead(H form.Name, list form.List) {
	if H != list[0] {
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
