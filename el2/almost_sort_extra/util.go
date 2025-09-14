package el2_almost_sort

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func mustMatchHead(H form.Name, list form.List) {
	if H != list[0] {
		panic(TypeErr)
	}
}

func not(b bool) bool {
	return !b
}
