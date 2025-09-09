package sorts

import (
	"fmt"

	"github.com/fbundle/sorts/expr"
)

// Pi - (x: A) -> (y: B(x)) similar to Arrow
// this is the universal quantifier
type Pi struct {
	A Sort
	B Dependent
}

func (s Pi) sortAttr() sortAttr {
	x := dummyTerm(s.A, "x")
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		repr:   expr.Node{expr.TermArrowSingle, expr.Node{expr.TermColon, Repr(x), Repr(s.A)}, expr.Node{s.B, Repr(x)}},
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
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
