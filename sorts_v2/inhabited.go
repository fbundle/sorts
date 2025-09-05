package sorts

import "github.com/fbundle/sorts/adt"

func newInhabited(ss SortSystem, sort Sort, child Sort) adt.Option[InhabitedSort] {
	if child.Level() != sort.Level()-1 {
		return adt.None[InhabitedSort]()
	}
	return adt.Some[InhabitedSort](inhabited{
		Sort:  sort,
		child: child,
		ss:    ss,
	})
}

type inhabited struct {
	Sort
	child Sort
	ss    SortSystem
}

func (i inhabited) Child() Sort {
	return i.child
}
