package sorts

import "github.com/fbundle/sorts/adt"

func newInhabited(ss SortSystem, sort Sort, child Sort) adt.Option[InhabitedSort] {
	if child.Level() != sort.Level()-1 {
		return adt.None[InhabitedSort]()
	}
	return adt.Some[InhabitedSort](inhabited{
		sort:  sort,
		child: child,
		ss:    ss,
	})
}

type inhabited struct {
	sort  Sort
	child Sort
	ss    SortSystem
}

func (i inhabited) Sort() Sort {
	return i.sort
}

func (i inhabited) Child() Sort {
	return i.child
}
