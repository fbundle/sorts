package sorts

import "fmt"

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("%s × %s", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Prod:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro - take (a: A) (b: B) give (a, b): A × B
func (s Prod) Intro(a Sort, b Sort) Sort {
	MustTermOf(a, s.A)
	MustTermOf(b, s.B)
	return dummyTerm(s, fmt.Sprintf("(<%s, %s> : %s)", Name(a), Name(b), Name(s)))
}

// Elim - take (t: A × B) give (a: A) and (b: B)
func (s Prod) Elim(t Sort) (left Sort, right Sort) {
	MustTermOf(t, s)
	a := dummyTerm(s.A, fmt.Sprintf("(left %s %s)", Name(s), Name(t)))
	b := dummyTerm(s.B, fmt.Sprintf("(right %s %s)", Name(s), Name(t)))
	return a, b
}
