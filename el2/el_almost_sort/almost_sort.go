package el_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type ParseFunc = func(form form.Form) AlmostSort
type ListParseFunc = func(parse ParseFunc, list form.List) AlmostSort

var TypeErr = fmt.Errorf("type_error")

func MustSort(as AlmostSort) sorts.Sort {
	return as.(ActualSort).Sort
}

// AlmostSort - almost a Sort - for example, a lambda
type AlmostSort interface {
	TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort // not nullable
}

// ActualSort - a Sort
type ActualSort struct {
	Sort sorts.Sort
}

func (s ActualSort) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	must(sa).termOf(s.Sort, parent)
	return s.Sort
}
