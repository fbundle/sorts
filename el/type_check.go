package el

import (
	"github.com/fbundle/sorts/sorts"
)

func typeCheckFunctionCall(cmdSort sorts.Sort, argSort sorts.Sort) (sorts.Sort, bool) {
	arrow, ok := sorts.Parent(cmdSort).(sorts.Arrow)
	if !ok {
		return nil, false
	}
	if !sorts.TermOf(argSort, arrow.A) {
		return nil, false
	}
	return arrow.B, true
}

func typeCheckBinding(parentSort sorts.Sort, valueSort sorts.Sort) bool {
	// TODO - type check
	return true
}
