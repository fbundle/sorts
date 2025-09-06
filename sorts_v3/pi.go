package sorts

import "fmt"

// Pi - (x: A) -> (y: B(x))
type Pi struct {
	A Inhabited
	B Dependent
}

func (s Pi) attr() sortAttr {
	x := s.A.Child
	level := max(Level(s.A), Level(s.B(x)))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("Π%s:%s. %s", Name(x), Name(s.A), Name(s.B(x))),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Atom:
				return false
			case Arrow:
				return false
			case Pi:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				if !LessEqual(d.A, s.A) {
					return false
				}
				y := d.A.Child
				return LessEqual(s.B(x), d.B(y))
			default:
				panic("type_error - should catch all types")
			}
		},
	}
}
