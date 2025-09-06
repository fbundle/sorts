package sorts

import "fmt"

// Pi - (x: A) -> (y: B(x)) similar to Arrow
// this is the universal quantifier
type Pi struct {
	A WithSort
	B Dependent
}

func (s Pi) sortAttr() sortAttr {
	x := dummyTerm(s.A, "x")
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		name:   fmt.Sprintf("Π(x:%s)%s(x)", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst WithSort) bool {
			switch d := dst.(type) {
			case Pi:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				if !LessEqual(d.A, s.A) {
					return false
				}
				y := dummyTerm(d.A, "y")
				dBy := d.B.Apply(y)
				return LessEqual(sBx, dBy)
			default:
				return false
			}
		},
	}
}

// Intro - take a func that maps (a: A) into (b: B(a))  give (f: Π(x:A)B(x))
func (s Pi) Intro(f func(WithSort) WithSort) WithSort {
	// verify
	a := dummyTerm(s.A, "a")
	b := f(a)
	sBa := s.B.Apply(a)
	mustTermOf(b, sBa) // TODO - think, shouldn't it have to check for every a of type A?

	return dummyTerm(s, fmt.Sprintf("(intro %s %p)", Name(s), f))
}

// Elim - take (f: Π(x:A)B(x)) (a: A) give (b: B(a)) - Modus Ponens
func (s Pi) Elim(f WithSort, a WithSort) WithSort {
	mustTermOf(f, s)
	mustTermOf(a, s.A)
	return dummyTerm(s, fmt.Sprintf("(elim %s %p)", Name(f), Name(a)))
}
