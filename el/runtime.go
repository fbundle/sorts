package el

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

// TODO - remove all Expr

// TODO check example.el for new syntax

// TODO - support multi argument function call

// TODO - support type inference

type Value struct {
	Sort sorts.Sort
	Expr Expr
}

type Frame struct {
	ordered_map.OrderedMap[Term, Value]
}

func (frame Frame) Get(key Term) (Value, bool) {
	if value, ok := frame.OrderedMap.Get(key); ok {
		return value, true
	}
	keyStr := string(key)
	if strings.HasPrefix(keyStr, "U_") {
		levelStr := strings.TrimPrefix(keyStr, "U_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			return Value{}, false
		}
		// U_0 is at universe level 1
		// U_1 is at universe level 2
		return Value{
			Sort: sorts.NewAtom(level+1, string(key), nil),
			Expr: key,
		}, true
	} else {
		return Value{}, false
	}

}

func Eval(frame Frame, expr Expr) (Frame, Value, error) {
	switch e := expr.(type) {
	case Term:
		value, ok := frame.Get(e)
		if !ok {
			return frame, Value{}, fmt.Errorf("undefined variable: %s", e)
		}
		return frame, value, nil
	case FunctionCall:
		frame, avalue, err := Eval(frame, e.Arg)
		if err != nil {
			return frame, Value{}, err
		}

		switch cmd := e.Cmd.(type) {
		case Lambda:
			frame, arg, err := Eval(frame, avalue.Expr)
			if err != nil {
				return frame, Value{}, err
			}
			callFrame := Frame{frame.Set(cmd.Param, arg)}
			return Eval(callFrame, cmd.Body)

		case Term:
			frame, c, err := Eval(frame, cmd)
			if err != nil {
				return frame, Value{}, err
			}
			if cmd == c.Expr {
				// term is not assigned hence cannot be simplified using Eval
				// like (succ (succ 0))
				// return expr form
				parentArrow, ok := sorts.Parent(c.Sort).(sorts.Arrow)
				if !ok {
					return frame, Value{}, fmt.Errorf("expected arrow: %T", sorts.Parent(c.Sort))
				}
				argToks := avalue.Expr.Marshal().Marshal()
				argStr := strings.Join(argToks, " ")
				return frame, Value{
					Sort: sorts.NewAtom(sorts.Level(parentArrow.B)-1, fmt.Sprintf("(%s %s)", string(cmd), argStr), parentArrow.B),
					Expr: FunctionCall{cmd, avalue.Expr},
				}, nil
			} else {
				return Eval(frame, FunctionCall{
					Cmd: c.Expr,
					Arg: avalue.Expr,
				})
			}
		default:
			return frame, Value{}, fmt.Errorf("unknown function: %T", e.Cmd)
		}
	case Lambda:
		return frame, Value{
			Sort: nil,
			Expr: e,
		}, nil
	case Define:
		frame, parent, err := Eval(frame, e.Type)
		if err != nil {
			return frame, Value{}, err
		}
		value := Value{
			Sort: sorts.NewAtom(sorts.Level(parent.Sort)-1, string(e.Name), parent.Sort),
			Expr: e.Name, // Expr == Term means not defined yet
		}

		frame = Frame{frame.Set(e.Name, value)}

		return frame, value, nil
	case Assign:
		value, ok := frame.Get(e.Name)
		if !ok {
			return frame, Value{}, fmt.Errorf("undefined variable: %s", e.Name)
		}
		frame, rvalue, err := Eval(frame, e.Value)
		if err != nil {
			return frame, Value{}, err
		}
		value.Expr = rvalue.Expr
		frame = Frame{frame.Set(e.Name, value)}
		return frame, value, nil
	case Chain:
		var err error
		for _, expr := range e.Init {
			frame, _, err = Eval(frame, expr)
			if err != nil {
				return frame, Value{}, err
			}
		}
		return Eval(frame, e.Tail)
	case Match:
		// match should not be hard, just compare Expr
		panic("not implemented")
	case Arrow:
		frame, avalue, err := Eval(frame, e.A)
		if err != nil {
			return frame, Value{}, err
		}
		frame, bvalue, err := Eval(frame, e.B)
		if err != nil {
			return frame, Value{}, err
		}
		return frame, Value{
			Sort: sorts.Arrow{
				A: avalue.Sort,
				B: bvalue.Sort,
			},
			Expr: e,
		}, nil
	default:
		return frame, Value{}, fmt.Errorf("unknown expression: %T", e)
	}
}
