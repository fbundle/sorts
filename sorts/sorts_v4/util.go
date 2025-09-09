package sorts

import "github.com/fbundle/sorts/expr"

type Runtime interface{} // interpreter

type Dependent struct {
	Runtime Runtime
	Param   expr.Term
	Body    expr.Expr
}

func (d Dependent) Apply(Arg Sort) Sort {
	// TODO - use runtime
	panic("not_implemented")
}

func (d Dependent) Repr() expr.Expr {
	return expr.Node{expr.TermArrowDouble, d.Param, d.Body}
}

// Inhabited - represents a Sort with at least one child
// (true theorems have proofs)
type Inhabited struct {
	Sort  Sort // underlying sort
	Child Sort
}

func (s Inhabited) sortAttr() sortAttr {
	return s.Sort.sortAttr()
}
