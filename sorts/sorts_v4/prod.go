package sorts

import "github.com/fbundle/sorts/expr"

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		repr:   expr.Node{expr.TermProd, Repr(s.A), Repr(s.B)},
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
