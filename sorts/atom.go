package sorts

type Atom[T comparable] struct {
	level  int
	name   T
	parent func() Sort[T]
}

func (s Atom[T]) sortAttr() sortAttr[T] {
	return sortAttr[T]{
		repr:   Node[T]{Value: s.name},
		level:  s.level,
		parent: s.parent(),
		lessEqual: func(u *Universe[T], dst Sort[T]) bool {
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
