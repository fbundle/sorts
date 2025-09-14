package el2_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/form"
)

var TypeErr = fmt.Errorf("type_err")

func mustMatchHead(H form.Name, list form.List) {
	if H != list[0] {
		panic(TypeErr)
	}
}

func mustSort(s almost_sort.AlmostSort) almost_sort.ActualSort {
	s1, ok := s.(almost_sort.ActualSort)
	if !ok {
		panic(TypeErr)
	}
	return s1
}

func not(b bool) bool {
	return !b
}
