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

func Eval(frame Frame, expr Expr) (Value, error) {
	switch e := expr.(type) {
	case Term:
		value, ok := frame.Get(e)
		if !ok {
			return Value{}, fmt.Errorf("undefined variable: %s", e)
		}
		return value, nil
	case FunctionCall:
		panic("not implemented")
	case Lambda:
		panic("not implemented")
	case Define:
		panic("not implemented")
	case Assign:
		panic("not implemented")
	case Chain:
		panic("not implemented")
	case Match:
		panic("not implemented")
	default:
		return Value{}, fmt.Errorf("unknown expression: %T", e)
	}
}
