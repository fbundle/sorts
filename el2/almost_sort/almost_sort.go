package almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

func NewActualSort(sort sorts.Sort) ActualSort {
	return ActualSort{sort: sort}
}

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	AttrAlmostSort() // use for type safety
	TypeCheck(sa sorts.SortAttr, parent ActualSort) ActualSort
}

// ActualSort - a sort
type ActualSort struct {
	sort sorts.Sort
}

func (s ActualSort) AttrAlmostSort() {}

func (s ActualSort) TypeCheck(a sorts.SortAttr, parent ActualSort) ActualSort {
	if !a.LessEqual(a.Parent(s.sort), parent.sort) {
		panic(TypeErr)
	}
	return ActualSort{s.sort}
}
