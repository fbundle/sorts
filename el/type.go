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

func (frame Frame) typeCheckBinding(parentSort sorts.Sort, name Term, value Expr) bool {
	// TODO - for functionCall - add dummy param into frame then check the body
	// TODO - for match - match all cases then check the value
	return true
}
