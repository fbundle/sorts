package almost_sort

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

var TypeErr = fmt.Errorf("type_error")

func NewActualSort(sort sorts.Sort) ActualSort {
	return ActualSort{sort: sort}
}

func MustSort(as AlmostSort) ActualSort {
	return as.(ActualSort)
}

func IsSort(as AlmostSort) bool {
	_, ok := as.(ActualSort)
	return ok
}

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	almostSortAttr()
	TypeCheck(sa sorts.SortAttr, parent ActualSort) ActualSort
}

// ActualSort - a sort
type ActualSort struct {
	sort sorts.Sort
}

func (s ActualSort) almostSortAttr() {}

func (s ActualSort) TypeCheck(a sorts.SortAttr, parent ActualSort) ActualSort {
	if !a.LessEqual(a.Parent(s.sort), parent.sort) {
		panic(TypeErr)
	}
	return ActualSort{s.sort}
}
