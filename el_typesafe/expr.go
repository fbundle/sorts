package el_typesafe

import "github.com/fbundle/sorts/sorts"

type Expr interface {
	mustExpr()
}

// partialExpr - those Expr that type cannot be totally resolved
// need type binding - like lambda
type partialExpr interface {
	Expr
	resolvePartial(frame Frame) (Frame, partialObject)
}

// totalExpr - those Expr that type can be totally resolved
// no need type binding - like term
type totalExpr interface {
	Expr
	resolveTotal(frame Frame) (Frame, totalObject)
}

type Term string

func (t Term) mustExpr() {}
func (t Term) resolveTotal(frame Frame) (Frame, totalObject) {
	panic("not implemented")
}

func resolvePartial(frame Frame, expr Expr) (Frame, partialObject) {
	switch expr := expr.(type) {
	case partialExpr:
		return expr.resolvePartial(frame)
	case totalExpr:
		frame, o := expr.resolveTotal(frame)
		return frame, o.partial()
	default:
		panic("type_error")
	}
}

func resolveTotal(frame Frame, parentSort sorts.Sort, expr Expr) (Frame, totalObject) {
	switch expr := expr.(type) {
	case partialExpr:
		frame, o := expr.resolvePartial(frame)
		return frame, o.typeCheck(frame, parentSort)
	case totalExpr:
		return expr.resolveTotal(frame)
	default:
		panic("type_error")
	}
}
