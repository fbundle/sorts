package el

import (
	"errors"

	"github.com/fbundle/sorts/sorts"
)

var typeErr = errors.New("type_error")

func mustTypeCheckFunctionCall(cmd Object, arg Object) Object {
	arrow, ok := cmd.parent().sort.(sorts.Arrow)
	if !ok {
		panic(typeErr)
	}
	if ok := sorts.TermOf(arg.sort, arrow.A); !ok {
		panic(typeErr)
	}

	bParent := newTerm()
	return newTerm()

}

// typeCheckFunctionCall - check if the function call is valid - (cmdSort argSort)
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

// typeCheckBinding - check if the binding is valid - (name: parentSort = expr)
func (frame Frame) typeCheckBinding(parentSort sorts.Sort, name Term, expr Expr) bool {
	if expr == Undef {
		return true
	}

	if lambda, ok := expr.(Lambda); ok {
		callFrame := frame
		// add function into frame
		callFrame, err := callFrame.set(name, sorts.NewTerm(parentSort, string(name)), name)
		if err != nil {
			return false
		}
		// add param into frame
		parentArrow, ok := parentSort.(sorts.Arrow)
		if !ok {
			return false
		}
		argValue := lambda.Param
		argSort := sorts.NewTerm(parentArrow.A, string(lambda.Param))
		callFrame, err = callFrame.set(lambda.Param, argSort, argValue)
		if err != nil {
			return false
		}
		// call
		return callFrame.typeCheckBinding(parentArrow.B, "", lambda.Body)
	}
	if match, ok := expr.(Match); ok {
		frame, condSort, _, err := match.Cond.Resolve(frame)
		if err != nil {
			return false
		}

		for _, c := range match.Cases {
			matchedFrame, err := reverseMatchPattern(frame, condSort, c.Comp)
			if err != nil {
				return false
			}
			if !matchedFrame.typeCheckBinding(parentSort, "", c.Value) {
				return false
			}
		}

		return frame.typeCheckBinding(parentSort, "", match.Final)
	}

	_, sort, _, err := expr.Resolve(frame)
	if err != nil {
		return false
	}
	if !sorts.TermOf(sort, parentSort) {
		return false
	}

	// TODO - for functionCall - add dummy param into frame then check the body
	// TODO - for match - match all cases then check the expr
	return true
}
