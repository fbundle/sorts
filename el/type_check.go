package el

import (
	"github.com/fbundle/sorts/sorts"
)

func typeCheckFunctionCall(cmdSort sorts.Sort, argSort sorts.Sort) sorts.Sort {
	arrow, ok := sorts.Parent(cmdSort).(sorts.Arrow)
	if !ok {
		return nil
	}
	if !sorts.TermOf(argSort, arrow.A) {
		return nil
	}
	return arrow.B
}
