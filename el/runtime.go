package el

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

// TODO - remove all Next

// TODO check example.el_v2 for new syntax - no longer using define and assign

// TODO - support multi argument function call

type Value struct {
	Sort sorts.Sort
	Next Expr
}

type Frame struct {
	dict ordered_map.OrderedMap[Term, Value]
}

func (frame Frame) Set(key Term, sort sorts.Sort, next Expr) Frame {
	return Frame{dict: frame.dict.Set(key, Value{Sort: sort, Next: next})}
}

func (frame Frame) Get(key Term) (sort sorts.Sort, next Expr, ok bool) {
	if value, ok := frame.dict.Get(key); ok {
		return value.Sort, value.Next, true
	}
	keyStr := string(key)
	if strings.HasPrefix(keyStr, "U_") {
		levelStr := strings.TrimPrefix(keyStr, "U_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return nil, nil, false
		}
		// U_0 is at universe level 1
		// U_1 is at universe level 2

		return sorts.NewAtom(level+1, string(key), nil), key, true
	} else {
		return nil, nil, false
	}
}

func (frame Frame) resolveTerm(term Term) (Frame, sorts.Sort, Expr, error) {
	sort, next, ok := frame.Get(term)
	if !ok {
		return frame, nil, nil, fmt.Errorf("variable not found: %s", term)
	}
	if term == next {
		// term is not assigned hence cannot be simplified using Resolve
		return frame, sort, next, nil
	}
	// recursive
	return frame.Resolve(next)
}

func (frame Frame) resolveFunctionCall(expr FunctionCall) (Frame, sorts.Sort, Expr, error) {
	frame, argSort, argValue, err := frame.Resolve(expr.Arg)
	if err != nil {
		return frame, nil, nil, err
	}
	frame, cmdSort, cmdValue, err := frame.Resolve(expr.Cmd)
	if err != nil {
		return frame, nil, nil, err
	}
	switch cmd := cmdValue.(type) {
	case Lambda:
		return frame.Set(cmd.Param, argSort, argValue).Resolve(cmd.Body)
	case Term:
		// type check
		parentArrow, ok := sorts.Parent(cmdSort).(sorts.Arrow)
		if !ok {
			return frame, nil, nil, fmt.Errorf("cmd is not a function %s", cmd)
		}
		if !sorts.TermOf(argSort, parentArrow.A) {
			return frame, nil, nil, fmt.Errorf("expected argument of type %s, got %s", sorts.Name(parentArrow.A), sorts.Name(argSort))
		}
		return frame, sorts.NewAtomTerm(parentArrow.B, fmt.Sprintf("(%s %s)", String(cmd), String(argValue))), FunctionCall{cmd, argValue}, nil
	default:
		return frame, nil, nil, fmt.Errorf("unknown function: %T", expr.Cmd)
	}

}

func (frame Frame) Resolve(expr Expr) (Frame, sorts.Sort, Expr, error) {
	switch expr := expr.(type) {
	case Term:
		return frame.resolveTerm(expr)
	case FunctionCall:
		return frame.resolveFunctionCall(expr)
	case Lambda:
		return frame, nil, expr, nil
	case Let:
		var err error
		var parentSort sorts.Sort
		for _, binding := range expr.Bindings {
			frame, parentSort, _, err = frame.Resolve(binding.Type)
			if err != nil {
				return frame, nil, nil, err
			}
			var value Expr
			if valTerm, ok := binding.Value.(Term); ok && valTerm == "undef" {
				value = binding.Name
			} else {
				value = binding.Value
			}
			frame = frame.Set(binding.Name, sorts.NewAtomTerm(parentSort, string(binding.Name)), value)
		}
		return frame.Resolve(expr.Final)
	case Match:
		// match should not be hard, just compare Next
		panic("not implemented")
	case Arrow:
		frame, aSort, _, err := frame.Resolve(expr.A)
		if err != nil {
			return frame, nil, nil, err
		}
		frame, bSort, _, err := frame.Resolve(expr.B)
		if err != nil {
			return frame, nil, nil, err
		}
		return frame, sorts.Arrow{
			A: aSort,
			B: bSort,
		}, expr, nil
	default:
		return frame, nil, nil, fmt.Errorf("unknown expression: %T", expr)
	}
}
