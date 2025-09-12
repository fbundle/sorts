package el

import (
	"github.com/fbundle/sorts/sorts"
)

func matchPattern(frame Frame, condSort sorts.Sort, condValue Expr, pattern Expr) (Frame, bool, error) {
	switch pattern := pattern.(type) {
	case Term:
		frame, err := frame.Set(pattern, condSort, condValue)
		if err != nil {
			return frame, false, err
		}
		return frame, true, nil
	case FunctionCall:
		if cond, ok := condValue.(FunctionCall); ok {
			frame, cmdSort, cmdValue, err := cond.Cmd.Resolve(frame)
			if err != nil {
				return frame, false, err
			}
			frame, matched1, err := matchPattern(frame, cmdSort, cmdValue, pattern.Cmd)
			if err != nil {
				return frame, false, err
			}
			frame, argSort, argValue, err := cond.Arg.Resolve(frame)
			if err != nil {
				return frame, false, err
			}
			frame, matched2, err := matchPattern(frame, argSort, argValue, pattern.Arg)
			return frame, matched1 && matched2, err
		} else {
			return frame, false, nil
		}

	default:
		return frame, false, nil // not comparable
	}
}

func alwaysMatchExact(frame Frame, condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error) {
	return frame, true, nil
}
