package sorts

import "fmt"

// Pi - (x: A) -> (y: B(x)) similar to Arrow
// this is the universal quantifier
type Pi struct {
	A Sort
	B Dependent
}

func (s Pi) attr() sortAttr {
	x := dummyTerm(s.A, "x")
	level := max(Level(s.A), Level(s.B(x)))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("Π%s:%s. %s", Name(x), Name(s.A), Name(s.B(x))),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Pi:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				if !LessEqual(d.A, s.A) {
					return false
				}
				y := dummyTerm(d.A, "y")
				return LessEqual(s.B(x), d.B(y))
			default:
				return false
			}
		},
	}
}

// Intro - take a func that maps (a: A) into (b: B(a))  give (x: Πx:A. B(x))
func (s Pi) Intro(f func(Sort) Sort) Sort {
	// verify
	a := dummyTerm(s.A, "a")
	b := f(a)
	mustTermOf(b, s.B(a)) // TODO - think, shouldn't it have to check for every a of type A?

	return dummyTerm(s, fmt.Sprintf("(intro %s %p)", Name(s), f))
}

func (s Pi) Elim(a Sort) Sort {
	
}
