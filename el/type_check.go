package el

import (
	"github.com/fbundle/sorts/sorts"
)

func (frame Frame) typeCheckFunctionCall(cmdSort sorts.Sort, argSort sorts.Sort) (sorts.Sort, bool) {
	arrow, ok := sorts.Parent(cmdSort).(sorts.Arrow)
	if !ok {
		return nil, false
	}
	if !sorts.TermOf(argSort, arrow.A) {
		return nil, false
	}
	return arrow.B, true
}

func (frame Frame) typeCheckBinding(parentSort sorts.Sort, value Expr) bool {
	return true
}
