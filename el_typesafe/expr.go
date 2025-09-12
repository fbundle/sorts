package el_typesafe

type Expr interface {
	mustExpr()
}

// partialExpr - those Expr that type cannot be totally resolved
// need type binding - like lambda
type partialExpr interface {
	Expr
	resolvePartial(frame Frame) (Frame, _partialObject)
}

// totalExpr - those Expr that type can be totally resolved
// no need type binding - like term
type totalExpr interface {
	Expr
	resolveTotal(frame Frame) (Frame, _totalObject)
}

type Term string

func (t Term) mustExpr() {}
func (t Term) resolveTotal(frame Frame) (Frame, _totalObject) {
	o := frame.get(t)
	if o.next == t { // term undef - loopback
		return frame, o
	}
	return o.next.resolveTotal(frame)
}

// FunctionCall - (cmd arg1 arg2 ...)
type FunctionCall struct {
	Cmd Expr
	Arg Expr
}

func (f FunctionCall) mustExpr() {}
func (f FunctionCall) resolveTotal(frame Frame) (Frame, _totalObject) {
	panic("not implemented")
}
