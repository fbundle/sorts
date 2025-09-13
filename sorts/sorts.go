package sorts

import (
	"errors"
)

var TypeErr = errors.New("type_error") // cannot recover

// Form - Union[Name, List]
type Form interface {
	mustForm()
}

type Name string

func (n Name) mustForm() {}

type List []Form

func (l List) mustForm() {}

type sortAttr struct {
	form      Form                // every Sort is identified with a Form
	level     int                 // universe Level
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}

type Sort interface {
	sortAttr(a SortAttr) sortAttr
}

// SortAttr - an object that can provide sort information
type SortAttr interface {
	Form(s any) Form
	Level(s Sort) int
	Parent(s Sort) Sort
	LessEqual(x Sort, y Sort) bool

	NameLessEqual(src Name, dst Name) bool
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
