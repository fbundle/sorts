package el

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

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
		switch cmd := e.Cmd.(type) {
		case Lambda:
		case Term:
			frame, c, err := Eval(frame, cmd)
			if err != nil {
				return frame, Value{}, err
			}
			if cmd == c.Expr { // term is not assigned
				arrow, ok := c.Sort.(sorts.Arrow)
				if !ok {
					return frame, Value{}, fmt.Errorf("expected arrow: %T", c.Sort)
				}
				// like (succ (succ 0))
				return frame, Value{
					Sort: arrow.B,
					Expr: FunctionCall{cmd, e.Args},
				}, nil
			} else {
				return Eval(frame, FunctionCall{
					Cmd:  c.Expr,
					Args: e.Args,
				})
			}
		default:
			return frame, Value{}, fmt.Errorf("unknown function: %T", e.Cmd)
		}

		panic("not implemented")
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
			Expr: e.Name,
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
		panic("not implemented")
	default:
		return frame, Value{}, fmt.Errorf("unknown expression: %T", e)
	}
}
