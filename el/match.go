package el

import (
	"github.com/fbundle/sorts/sorts"
)

func (frame Frame) match(condSort sorts.Sort, condValue Expr, comp Expr) (Frame, bool, error) {
	// Try pattern matching first - don't resolve the pattern
	if frame, matched, err := frame.patternMatch(condValue, comp); err != nil {
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

// patternMatch attempts to match a pattern against a value and bind variables
func (frame Frame) patternMatch(value Expr, pattern Expr) (Frame, bool, error) {
	switch pattern := pattern.(type) {
	case Term:
		return frame.matchVariable(value, pattern)
	case FunctionCall:
		return frame.matchConstructor(value, pattern)
	default:
		// For other types, fall back to exact matching
		return frame, String(pattern) == String(value), nil
	}
}

// matchVariable handles pattern matching against a variable pattern
func (frame Frame) matchVariable(value Expr, pattern Term) (Frame, bool, error) {
	// Always bind the variable to the value (allow rebinding for recursive functions)
	// Infer sort from the value's sort
	_, valueSort, _, err := frame.Resolve(value)
	if err != nil {
		// If we can't resolve the value, use a generic sort
		valueSort = sorts.NewAtom(1, "any", nil)
	}
	frame, err = frame.Set(pattern, valueSort, value)
	if err != nil {
		return frame, false, err
	}
	return frame, true, nil
}

// matchConstructor handles pattern matching against constructor patterns like (succ z)
func (frame Frame) matchConstructor(value Expr, pattern FunctionCall) (Frame, bool, error) {
	// Resolve the value to get its actual form
	_, _, resolvedValue, err := frame.Resolve(value)
	if err != nil {
		return frame, false, err
	}

	// Check if the resolved value is also a function call
	if valueCall, ok := resolvedValue.(FunctionCall); ok {
		// Check if constructors match
		if String(pattern.Cmd) == String(valueCall.Cmd) {
			// Constructors match, recursively match arguments
			return frame.patternMatch(valueCall.Arg, pattern.Arg)
		}
	}

	return frame, false, nil
}
