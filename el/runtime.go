package el

import (
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
		frame.Get(e)
	}
}
