package sorts

import "fmt"

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) attr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("%s × %s", Name(s.A), Name(s.B)),
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
	mustTermOf(a, s.A)
	mustTermOf(b, s.B)
	return dummyTerm(s, fmt.Sprintf("(%s, %s)", Name(a), Name(b)))
}

// Left - take (x: A × B) give (a: A)
func (s Prod) Left(x Sort) Sort {
	mustTermOf(x, s)
	return dummyTerm(s.A, fmt.Sprintf("(left %s %s)", Name(s), Name(x)))
}

// Right - take (x: A × B) give (b: B)
func (s Prod) Right(x Sort) Sort {
	mustTermOf(x, s)
	return dummyTerm(s.B, fmt.Sprintf("(right %s %s)", Name(s), Name(x)))
}
