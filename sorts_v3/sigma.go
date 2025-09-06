package sorts

import "fmt"

// Sigma - (x: A, y: B(x)) , similar to Prod
// this is the existential quantifier
type Sigma struct {
	A Inhabited
	B Dependent
}

func (s Sigma) attr() sortAttr {
	x := dummyTerm(s.A, "x")
	level := max(Level(s.A), Level(s.B(x)))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("Σ%s:%s. %s", Name(x), Name(s.A), Name(s.B(x))),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sigma:
				y := dummyTerm(d.A, "y")
				return LessEqual(s.A, d.A) && LessEqual(s.B(x), d.B(y))
			default:
				return false
			}
		},
	}
}
