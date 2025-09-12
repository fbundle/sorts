package el

import (
	"errors"
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

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

// reverseMatchPattern - similar to matchPattern, but always matches to get the updated frame
func reverseMatchPattern(frame Frame, condSort sorts.Sort, pattern Expr) (Frame, error) {
	fmt.Println("match", pattern, "of type", sorts.Name(sorts.Parent(condSort)))
	switch pattern := pattern.(type) {
	case Exact:
		frame, _, _, err := pattern.Expr.Resolve(frame)
		return frame, err
	case Term:
		return frame.Set(pattern, condSort, pattern)
	case FunctionCall:
		frame, cmdSort, _, err := pattern.Cmd.Resolve(frame)
		if err != nil {
			return frame, err
		}
		cmdArrow, ok := sorts.Parent(cmdSort).(sorts.Arrow)
		if !ok {
			return frame, errors.New("expected function")
		}
		if cmdArrow.B != sorts.Parent(condSort) {
			return frame, errors.New("wrong output type")
		}
		return reverseMatchPattern(frame, sorts.NewTerm(cmdArrow.A, String(pattern.Arg)), pattern.Arg)
	default:
		return frame, nil // not comparable
	}
}
