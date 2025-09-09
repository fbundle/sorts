package sorts

import "github.com/fbundle/sorts/expr"

type Sum struct {
	A Sort
	B Sort
}

func (s Sum) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		repr:   expr.Node{expr.TermSum, Repr(s.A), Repr(s.B)},
		level:  level,
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
