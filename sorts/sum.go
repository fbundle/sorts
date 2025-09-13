package sorts

import "fmt"

type Sum struct {
	A Sort
	B Sort
}

func (s Sum) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("%s + %s", Name(s.A), Name(s.B)),
		level:  level,
		parent: Sum{A: Parent(s.A), B: Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sum:
				return SubTypeOf(s.A, d.A) && SubTypeOf(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro - IntroLeft or IntroRight
func (s Sum) Intro(a Sort, b Sort) Sort {
	if a != nil {
		// IntroLeft - take (a: A) give (x: A + B)
		mustTermOf(a, s.A)
		return a
	} else {
		// IntroRight - take (b: B) give (x: A + B)
		mustTermOf(b, s.B)
		return b
	}
}

// ByCases - take (t: A + B) (h1: A -> X) (h2: B -> X) give (x: X)
func (s Sum) ByCases(t Sort, h1 Sort, h2 Sort) Sort {
	mustTermOf(t, s)
	X := Parent(h1).(Arrow).B
	mustTermOf(h1, Arrow{s.A, X})
	mustTermOf(h2, Arrow{s.B, X})

	name := fmt.Sprintf("(by_cases %s %s %s)", Name(t), Name(h1), Name(h2))
	return NewTerm(name, X)
}
