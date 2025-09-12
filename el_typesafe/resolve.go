package el_typesafe

func resolvePartial(frame Frame, expr Expr) (Frame, _partialObject) {
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

func resolveTotal(frame Frame, parent _totalObject, expr Expr) (Frame, _totalObject) {
	switch expr := expr.(type) {
	case partialExpr:
		frame, o := expr.resolvePartial(frame)
		return frame, o.typeCheck(frame, parent)
	case totalExpr:
		return expr.resolveTotal(frame)
	default:
		panic("type_error")
	}
}
