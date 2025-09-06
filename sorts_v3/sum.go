package sorts

import "fmt"

type Sum struct {
	A WithSort
	B WithSort
}

func (s Sum) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("%s + %s", Name(s.A), Name(s.B)),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst WithSort) bool {
			switch d := dst.(type) {
			case Sum:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro -
func (s Sum) Intro(a WithSort, b WithSort) WithSort {
	if a != nil {
		return s.IntroLeft(a)
	} else {
		return s.IntroRight(b)
	}
}

// IntroLeft - take (a: A) give (x: A + B)
func (s Sum) IntroLeft(a WithSort) WithSort {
	mustTermOf(a, s.A)
	return dummyTerm(s, fmt.Sprintf("(intro_left %s %s)", Name(s), Name(a)))
}

// IntroRight - take (b: B) give (x: A + B)
func (s Sum) IntroRight(b WithSort) WithSort {
	mustTermOf(b, s.B)
	return dummyTerm(s, fmt.Sprintf("(intro_right %s %s)", Name(s), Name(b)))
}

// ByCases - take (t: A + B) (h1: A -> X) (h2: B -> X) give (x: X)
func (s Sum) ByCases(t WithSort, h1 WithSort, h2 WithSort) WithSort {
	mustTermOf(t, s)
	mustSubType(Parent(h1).(Arrow).A, s.A)
	mustSubType(Parent(h2).(Arrow).A, s.B)

	mustSubType(Parent(h1).(Arrow).B, Parent(h2).(Arrow).B)

	X := Parent(h1).(Arrow).B

	return dummyTerm(X, fmt.Sprintf("(by_cases %s %s %s)", Name(s), Name(h1), Name(h2)))
}
