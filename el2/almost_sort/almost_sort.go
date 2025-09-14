package el2_almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

type Runtime interface {
	Parse(form form.Form) (Runtime, AlmostSort)
	SortAttr() sorts.SortAttr
	Get(name form.Name) ActualSort
	Set(name form.Name, sort ActualSort) Runtime
}

type ListParseFunc = func(ctx Runtime, list form.List) (Runtime, AlmostSort)

var TypeErr = fmt.Errorf("type_error")

func MustSort(as AlmostSort) ActualSort {
	return as.(ActualSort)
}

func IsSort(as AlmostSort) bool {
	_, ok := as.(ActualSort)
	return ok
}

// AlmostSort - almost a Sort - for example, a lambda
type AlmostSort interface {
	almostSortAttr()
	TypeCheck(sa sorts.SortAttr, parent ActualSort) ActualSort
}

// ActualSort - a Sort
type ActualSort struct {
	Sort sorts.Sort
}

func (s ActualSort) almostSortAttr() {}

func (s ActualSort) TypeCheck(sa sorts.SortAttr, parent ActualSort) ActualSort {
	must(sa).termOf(s.Sort, parent.Sort)
	return ActualSort{s.Sort}
}
