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
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		level:  level,
		name:   fmt.Sprintf("Σ(x:%s)%s(x)", Name(s.A), Name(s.B)),
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sigma:
				y := dummyTerm(d.A, "y")
				dBy := d.B.Apply(y)
				return LessEqual(s.A, d.A) && LessEqual(sBx, dBy)
			default:
				return false
			}
		},
	}
}

// Intro - take (a: A) (b: B(a)) give (t: )
func (s Sigma) Intro(a Sort, b Sort) Sort {

}
