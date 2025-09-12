package el

import (
	"github.com/fbundle/sorts/sorts"
)

func (frame Frame) match(condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error) {
	// Try pattern matching first - don't resolve the pattern
	if frame, matched, err := frame.matchPattern(condSort, condValue, comp); err != nil {
		return frame, false, err
	} else if matched {
		return frame, true, nil
	}

	// Fall back to exact matching - resolve the comparison for exact match
	frame, compSort, compValue, err := frame.Resolve(comp)
	if err != nil {
		return frame, false, err
	}
	if compSort == condSort {
		if String(compValue) == String(condValue) {
			return frame, true, nil
		}
	}
	return frame, false, nil
}

func (frame Frame) matchPattern(condSort sorts.Sort, condValue Expr, pattern Expr) (Frame, bool, error) {
	switch pattern := pattern.(type) {
	case Term:
		return frame.matchTerm(condSort, condValue, pattern)
	case FunctionCall:
		return frame.matchFunctionCall(condSort, condValue, pattern)
	default:
		return frame, false, nil // not comparable
	}
}

func (frame Frame) matchTerm(condSort sorts.Sort, condValue Expr, pattern Term) (Frame, bool, error) {
	frame, err := frame.Set(pattern, condSort, condValue)
	if err != nil {
		return frame, false, err
	}
	return frame, true, nil
}

func (frame Frame) matchFunctionCall(condSort sorts.Sort, condValue Expr, pattern FunctionCall) (Frame, bool, error) {
	if condValue, ok := condValue.(FunctionCall); ok {
		frame, cmdSort, cmdValue, err := frame.Resolve(condValue.Cmd)
		if err != nil {
			return frame, false, err
		}
		frame, argSort, argValue, err := frame.Resolve(condValue.Arg)
		if err != nil {
			return frame, false, err
		}

		frame, matched, err := frame.matchPattern(cmdSort, cmdValue, pattern.Cmd)
		if err != nil {
			return frame, false, err
		}
		if !matched {
			return frame, false, nil
		}
		frame, matched, err = frame.matchPattern(argSort, argValue, pattern.Arg)
		return frame, matched, err
	}
	return frame, false, nil
}
