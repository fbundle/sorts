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

func (frame Frame) Set(key Term, sort sorts.Sort, next Expr) (Frame, error) {
	if sort == nil || next == nil {
		return frame, fmt.Errorf("cannot set nil sort or next: %v, %v, %v", key, sort, next)
	}
	return Frame{dict: frame.dict.Set(key, Value{Sort: sort, Next: next})}, nil
}

func (frame Frame) Get(key Term) (sort sorts.Sort, next Expr, err error) {
	notFoundErr := fmt.Errorf("variable not found: %s", key)
	if value, ok := frame.dict.Get(key); ok {
		return value.Sort, value.Next, nil
	}
	keyStr := string(key)
	if strings.HasPrefix(keyStr, "U_") {
		levelStr := strings.TrimPrefix(keyStr, "U_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return nil, nil, notFoundErr
		}
		// U_0 is at universe level 1
		// U_1 is at universe level 2
		return sorts.NewAtom(level+1, string(key), nil), key, nil
	} else if strings.HasPrefix(keyStr, "Any_") {
		levelStr := strings.TrimPrefix(keyStr, "Any_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return nil, nil, notFoundErr
		}
		// Any_0 is at universe level 1
		// Any_1 is at universe level 2
		return sorts.NewAtom(level+1, sorts.TerminalName, nil), key, nil
	} else if strings.HasPrefix(keyStr, "Unit_") {
		levelStr := strings.TrimPrefix(keyStr, "Unit_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return nil, nil, notFoundErr
		}
		// Unit_0 is at universe level 1
		// Unit_1 is at universe level 2
		return sorts.NewAtom(level+1, sorts.InitialName, nil), key, nil
	} else {
		return nil, nil, notFoundErr
	}
}

func (frame Frame) resolveTerm(term Term) (Frame, sorts.Sort, Expr, error) {
	sort, next, err := frame.Get(term)
	if err != nil {
		return frame, nil, nil, err
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
		frame, err := frame.Set(cmd.Param, argSort, argValue)
		if err != nil {
			return frame, nil, nil, err
		}
		return frame.Resolve(cmd.Body)
	case Term:
		if B, ok := frame.typeCheckFunctionCall(cmdSort, argSort); ok {
			return frame, sorts.NewTerm(B, fmt.Sprintf("(%s %s)", String(cmd), String(argValue))), FunctionCall{cmd, argValue}, nil
		}
		return frame, nil, nil, fmt.Errorf("type_error: cmd %s, arg %s", sorts.Name(cmdSort), sorts.Name(argSort))
	default:
		return frame, nil, nil, fmt.Errorf("unknown function: %T", cmd)
	}
}

func (frame Frame) resolveLet(expr Let) (Frame, sorts.Sort, Expr, error) {
	var err error
	var parentSort sorts.Sort
	var valueSort sorts.Sort
	var value Expr
	for _, binding := range expr.Bindings {
		frame, parentSort, _, err = frame.Resolve(binding.Type)
		if err != nil {
			return frame, nil, nil, err
		}

		if valTerm, ok := binding.Value.(Term); ok && valTerm == "undef" {
			value = binding.Name
		} else {
			value = binding.Value
			if !frame.typeCheckBinding(parentSort, binding.Name, value) {
				return frame, nil, nil, fmt.Errorf("type_error: type %s, value %s", sorts.Name(parentSort), sorts.Name(valueSort))
			}
		}
		frame, err = frame.Set(binding.Name, sorts.NewTerm(parentSort, string(binding.Name)), value)
		if err != nil {
			return frame, nil, nil, err
		}
	}
	return frame.Resolve(expr.Final)
}

func (frame Frame) resolveMatch(expr Match) (Frame, sorts.Sort, Expr, error) {
	frame, condSort, condValue, err := frame.Resolve(expr.Cond)
	if err != nil {
		return frame, nil, nil, err
	}

	var matched bool
	for _, c := range expr.Cases {
		frame, matched, err = frame.match(condSort, condValue, c.Comp)
		if matched {
			return frame.Resolve(c.Value)
		}
	}
	return frame.Resolve(expr.Final)
}
func (frame Frame) resolveArrow(expr Arrow) (Frame, sorts.Sort, Expr, error) {
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
}

func (frame Frame) resolveSum(expr Sum) (Frame, sorts.Sort, Expr, error) {
	frame, aSort, _, err := frame.Resolve(expr.A)
	if err != nil {
		return frame, nil, nil, err
	}
	frame, bSort, _, err := frame.Resolve(expr.B)
	if err != nil {
		return frame, nil, nil, err
	}
	return frame, sorts.Sum{
		A: aSort,
		B: bSort,
	}, expr, nil
}

func (frame Frame) resolveProd(expr Prod) (Frame, sorts.Sort, Expr, error) {
	frame, aSort, _, err := frame.Resolve(expr.A)
	if err != nil {
		return frame, nil, nil, err
	}
	frame, bSort, _, err := frame.Resolve(expr.B)
	if err != nil {
		return frame, nil, nil, err
	}
	return frame, sorts.Prod{
		A: aSort,
		B: bSort,
	}, expr, nil
}

func (frame Frame) Resolve(expr Expr) (Frame, sorts.Sort, Expr, error) {
	switch expr := expr.(type) {
	case Term:
		return frame.resolveTerm(expr)
	case FunctionCall:
		return frame.resolveFunctionCall(expr)
	case Let:
		return frame.resolveLet(expr)
	case Match:
		return frame.resolveMatch(expr)
	case Arrow:
		return frame.resolveArrow(expr)
	case Sum:
		return frame.resolveSum(expr)
	case Prod:
		return frame.resolveProd(expr)
	default:
		// for other types, just return the original expression without sort
		// hence these cannot be bound to variables
		return frame, nil, expr, nil
	}
}
