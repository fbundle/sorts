package el_typesafe

import "github.com/fbundle/sorts/sorts"

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

func resolveTotal(frame Frame, parentSort sorts.Sort, expr Expr) (Frame, _totalObject) {
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
