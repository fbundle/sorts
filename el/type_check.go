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
	if lambda, ok := value.(Lambda); ok {
		// suppose name is already of the correct type - for recursive function
		frame, err := frame.Set(name, sorts.NewTerm(parentSort, string(name)), name)
		if err != nil {
			return false
		}
		// parentSort must be arrow
		arrow, ok := parentSort.(sorts.Arrow)
		if !ok {
			return false
		}
		paramName := lambda.Param
		paramSort := sorts.NewTerm(arrow.A, string(paramName)) // dummy param
		frame, err = frame.Set(paramName, paramSort, paramName)

		return frame.typeCheckBinding(arrow.B, "dummy", lambda.Body)
	} else {
		_, sort, _, err := frame.Resolve(value)
		if err != nil {
			return false
		}
		if !sorts.TermOf(sort, parentSort) {
			return false
		}
		return true
	}
}
