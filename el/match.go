package el

import (
	"github.com/fbundle/sorts/sorts"
)

func (frame Frame) match(condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error) {
	// Try pattern matching first - don't Resolve the pattern
	newFrame, matched, err := frame.matchPattern(condSort, condValue, comp)
	if err != nil {
		return newFrame, false, err
	}
	if matched {
		return newFrame, true, nil
	}

	// fall back to exact-match
	frame, _, compValue, err := comp.Resolve(frame)
	if err != nil {
		return frame, false, err
	}
	return frame, String(compValue) == String(condValue), nil
}

func (frame Frame) matchPattern(condSort sorts.Sort, condValue Expr, pattern Expr) (Frame, bool, error) {
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
			frame, matched1, err := frame.matchPattern(cmdSort, cmdValue, pattern.Cmd)
			if err != nil {
				return frame, false, err
			}
			frame, argSort, argValue, err := cond.Arg.Resolve(frame)
			if err != nil {
				return frame, false, err
			}
			frame, matched2, err := frame.matchPattern(argSort, argValue, pattern.Arg)
			return frame, matched1 && matched2, err
		} else {
			return frame, false, nil
		}

	default:
		return frame, false, nil // not comparable
	}
}
