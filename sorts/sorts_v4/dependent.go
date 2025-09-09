package sorts

import "github.com/fbundle/sorts/expr"

type Dependent struct {
	Param expr.Term
	Body  expr.Expr
}

func (s Dependent) sortAttr() sortAttr {
	return sortAttr{
		repr: expr.Node{expr.TermArrowDouble, s.Param, s.Body},
	}
}

func (s Dependent) BetaReduction(arg expr.Expr) expr.Expr {
	panic("not_implemented")
}
