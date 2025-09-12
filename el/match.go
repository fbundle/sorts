package el

import (
	"github.com/fbundle/sorts/sorts"
)

type matchFunc func(frame Frame, condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error)

func chainMatchFunc(matchFuncs ...matchFunc) matchFunc {
	return func(frame Frame, condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error) {
		for _, matchFunc := range matchFuncs {
			newFrame, matched, err := matchFunc(frame, condSort, condValue, comp)
			if err != nil {
				return newFrame, false, err
			}
			if matched {
				return newFrame, true, nil
			}
		}
		return frame, false, nil
	}
}

var match = chainMatchFunc(matchPattern)

func matchPattern(frame Frame, condSort sorts.Sort, condValue Expr, pattern Expr) (Frame, bool, error) {
	switch pattern := pattern.(type) {
	case Exact:
		frame, _, compValue, err := pattern.Expr.Resolve(frame)
		if err != nil {
			return frame, false, err
		}
		return frame, String(compValue) == String(condValue), nil
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
			frame, matched, err := matchPattern(frame, cmdSort, cmdValue, pattern.Cmd)
			if err != nil {
				return frame, false, err
			}
			if !matched {
				return frame, false, nil
			}
			frame, argSort, argValue, err := cond.Arg.Resolve(frame)
			if err != nil {
				return frame, false, err
			}
			return matchPattern(frame, argSort, argValue, pattern.Arg)
		} else {
			return frame, false, nil
		}

	default:
		return frame, false, nil // not comparable
	}
}

// alwaysMatchPattern - similar to matchPattern, but always matches to get the updated frame
func alwaysMatchPattern(frame Frame, condSort sorts.Sort, condValue Expr, pattern Expr) (Frame, error) {
	switch pattern := pattern.(type) {
	case Exact:
		frame, _, _, err := pattern.Expr.Resolve(frame)
		return frame, err
	case Term:
		return frame.Set(pattern, condSort, condValue)
	case FunctionCall:
		if cond, ok := condValue.(FunctionCall); ok {
			frame, cmdSort, cmdValue, err := cond.Cmd.Resolve(frame)
			if err != nil {
				return frame, err
			}
			frame, err = alwaysMatchPattern(frame, cmdSort, cmdValue, pattern.Cmd)
			if err != nil {
				return frame, err
			}
			frame, argSort, argValue, err := cond.Arg.Resolve(frame)
			if err != nil {
				return frame, err
			}
			return alwaysMatchPattern(frame, argSort, argValue, pattern.Arg)
		} else {
			return frame, nil
		}

	default:
		return frame, nil // not comparable
	}
}
