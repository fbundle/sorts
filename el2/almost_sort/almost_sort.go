package el2_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type ParseFunc = func(form form.Form) AlmostSort
type ListParseFunc = func(parse ParseFunc, list form.List) AlmostSort
type ListParseFuncWithHead = func(H form.Name) ListParseFunc

var TypeErr = fmt.Errorf("type_error")

func MustSort(as AlmostSort) sorts.Sort {
	return as.(ActualSort).Sort
}

func IsSort(as AlmostSort) bool {
	_, ok := as.(ActualSort)
	return ok
}

// AlmostSort - almost a Sort - for example, a lambda
type AlmostSort interface {
	almostSortAttr()
	TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort // not nullable
}

// ActualSort - a Sort
type ActualSort struct {
	Sort sorts.Sort
}

func (s ActualSort) almostSortAttr() {}

func (s ActualSort) TypeCheck(sa sorts.SortAttr, parent sorts.Sort) sorts.Sort {
	must(sa).termOf(s.Sort, parent)
	return s.Sort
}
