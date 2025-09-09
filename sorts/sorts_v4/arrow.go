package sorts

import (
	"github.com/fbundle/sorts/expr"
)

// Arrow - (A -> B)
type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		repr:   expr.Node{expr.TermArrowSingle, Repr(s.A), Repr(s.B)},
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Arrow:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !LessEqual(d.A, s.A) {
					return false
				}
				return LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}
