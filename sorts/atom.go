package sorts

func newAtomChain[T term](level int, name func(int) T) Atom[T] {
	return Atom[T]{
		level: level,
		name:  name(level),
		parent: func() Sort[T] {
			return newAtomChain[T](level+1, name)
		},
	}
}
func newAtomTerm[T term](u Universe[T], name T, parent Sort[T]) Atom[T] {
	return Atom[T]{
		level: u.Level(parent) - 1,
		name:  name,
		parent: func() Sort[T] {
			return parent
		},
	}
}

type Atom[T term] struct {
	level  int
	name   T
	parent func() Sort[T]
}

func (s Atom[T]) sortAttr() sortAttr[T] {
	return sortAttr[T]{
		repr:   Node[T]{Value: s.name},
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(u Universe[T], dst Sort[T]) bool {
			switch d := dst.(type) {
			case Atom[T]:
				if s.level != d.level {
					return false
				}
				return u.lessEqual(s.name, d.name)
			default:
				return false
			}
		},
	}
}
