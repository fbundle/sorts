package almost_sort_extra

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

func NewActualSort(sort sorts.Sort) ActualSort {
	return ActualSort{_sort: sort}
}

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	attrAlmostSort() // use for type safety
	TypeCheck(ctx Context, parent ActualSort) ActualSort
}

// ActualSort - a sort
type ActualSort struct {
	_sort sorts.Sort
}

func (s ActualSort) AttrAlmostSort() {}

func (s ActualSort) Repr() sorts.Sort {
	return s._sort
}

func (s ActualSort) TypeCheck(a sorts.SortAttr, parent ActualSort) ActualSort {
	if !a.LessEqual(a.Parent(s._sort), parent._sort) {
		panic(TypeErr)
	}
	return ActualSort{s._sort}
}
