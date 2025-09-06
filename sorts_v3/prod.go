package sorts

import "fmt"

type Prod struct {
	A WithSort
	B WithSort
}

func (s Prod) attr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("%s × %s", Name(s.A), Name(s.B)),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst WithSort) bool {
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
func (s Prod) Intro(a WithSort, b WithSort) WithSort {
	mustTermOf(a, s.A)
	mustTermOf(b, s.B)
	return dummyTerm(s, fmt.Sprintf("(%s, %s)", Name(a), Name(b)))
}

// Elim - take (t: A × B) give (a: A) and (b: B)
func (s Prod) Elim(t WithSort) (left WithSort, right WithSort) {
	mustTermOf(t, s)
	a := dummyTerm(s.A, fmt.Sprintf("(left %s %s)", Name(s), Name(t)))
	b := dummyTerm(s.B, fmt.Sprintf("(right %s %s)", Name(s), Name(t)))
	return a, b
}
