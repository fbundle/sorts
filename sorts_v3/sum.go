package sorts

import "fmt"

type Sum struct {
	A Sort
	B Sort
}

func (s Sum) attr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("%s + %s", Name(s.A), Name(s.B)),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sum:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// IntroLeft - take (a: A) give (x: A + B)
func (s Sum) IntroLeft(a Sort) Sort {
	mustTermOf(a, s.A)
	return dummyTerm(s, fmt.Sprintf("(intro_left %s %s)", Name(s), Name(a)))
}

// IntroRight - take (b: B) give (x: A + B)
func (s Sum) IntroRight(b Sort) Sort {
	mustTermOf(b, s.B)
	return dummyTerm(s, fmt.Sprintf("(intro_right %s %s)", Name(s), Name(b)))
}

// ByCases - take (t: A + B) (h1: A -> X) (h2: B -> X) give (x: X)
func (s Sum) ByCases(t Sort, h1 Sort, h2 Sort) Sort {
	mustTermOf(t, s)
	mustSubType(Parent(h1).(Arrow).A, s.A)
	mustSubType(Parent(h2).(Arrow).A, s.B)

	mustSubType(Parent(h1).(Arrow).B, Parent(h2).(Arrow).B)

	X := Parent(h1).(Arrow).B

	return dummyTerm(X, fmt.Sprintf("(by_cases %s %s %s)", Name(s), Name(h1), Name(h2)))
}
