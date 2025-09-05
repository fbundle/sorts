package sorts

import "github.com/fbundle/sorts/adt"

func newAtom(ss SortSystem, level int, name string, parent Sort) adt.Option[Sort] {
	if parent != nil && parent.Level() != level+1 {
		return adt.None[Sort]()
	}
	return adt.Some[Sort](atom{
		level:  level,
		name:   name,
		parent: parent,
		ss:     ss,
	})
}

type atom struct {
	level  int
	name   string
	parent Sort
	ss     SortSystem
}

func (s atom) Level() int {
	return s.level
}

func (s atom) Name() string {
	return s.name
}

func (s atom) Parent() Sort {
	if s.parent != nil {
		return s.parent
	}
	// default parent
	return s.ss.Default(s.level + 1)
}

func (s atom) LessEqual(dst Sort) bool {
	switch d := dst.(type) {
	case atom:
		if s.level != d.level {
			return false
		}
		return s.ss.LessEqual(s.name, d.name)
	case arrow:
		// cannot compare atom and arrow
		return false
	default:
		panic("unreachable")
	}
}
