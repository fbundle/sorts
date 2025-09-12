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
