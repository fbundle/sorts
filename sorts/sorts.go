package sorts

import (
	"fmt"
)

var TypeErr = fmt.Errorf("type_err") // cannot recover

type ListCompileFunc = func(parse func(form Form) Sort, list List) Sort

type sortAttr struct {
	form      Form                // every Sort is identified with a Form
	level     int                 // universe Level
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}

type Sort interface {
	sortAttr(a SortAttr) sortAttr
}

// SortAttr - an almost_sort that can provide sort information
type SortAttr interface {
	Form(s any) Form
	Level(s Sort) int
	Parent(s Sort) Sort
	LessEqual(x Sort, y Sort) bool

	GetRule(src Name, dst Name) bool
}

func GetForm(a SortAttr, s any) Form {
	switch s := s.(type) {
	case Sort:
		return s.sortAttr(a).form
	case Dependent:
		return s.Repr
	default:
		panic(TypeErr)
	}
}

func GetLevel(a SortAttr, s Sort) int {
	return s.sortAttr(a).level
}
func GetParent(a SortAttr, s Sort) Sort {
	return s.sortAttr(a).parent
}
func GetLessEqual(a SortAttr, x Sort, y Sort) bool {
	return x.sortAttr(a).lessEqual(y)
}
