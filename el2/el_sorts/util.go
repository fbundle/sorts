package almost_sort_extra

import (
	"github.com/fbundle/sorts/form"
)

func mustMatchHead(H form.Name, list form.List) {
	if H != list[0] {
		panic(TypeErr)
	}
}

func mustSort(s Sort) typeSort {
	s1, ok := s.(typeSort)
	if !ok {
		panic(TypeErr)
	}
	return s1
}

func not(b bool) bool {
	return !b
}

func Form(ctx Context, s Sort) form.Form {
	return s.attrSort(ctx).form
}
