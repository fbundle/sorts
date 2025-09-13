package el

import (
	"strconv"
	"strings"

	"github.com/fbundle/sorts/persistent/ordered_map"
)

type Frame struct {
	dict ordered_map.OrderedMap[Term, _object]
}

func (frame Frame) set(key Term, o _object) Frame {
	if !o.isTotal() {
		panic("set_must_total")
	}
	return Frame{dict: frame.dict.Set(key, o)}
}

func (frame Frame) get(key Term) _object {
	if o, ok := frame.dict.Get(key); ok {
		return o
	}
	return builtinValue(key)
}

func builtinValue(key Term) _object {
	keyStr := string(key)

	// integer
	intTypeName := Term("int")
	intType := newType(1, intTypeName)
	if key == intTypeName {
		return intType
	}

	if _, err := strconv.Atoi(keyStr); err == nil {
		return newTerm(key, intType)
	}

	// universe
	if strings.HasPrefix(keyStr, "U_") {
		levelStr := strings.TrimPrefix(keyStr, "U_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			panic(err)
		}
		return uType(level + 1)
	}

	// any type
	if strings.HasPrefix(keyStr, "Any_") {
		levelStr := strings.TrimPrefix(keyStr, "Any_")
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			panic(err)
		}
		return anyType(level + 1)
	}

	// unit type
	if strings.HasPrefix(keyStr, "Unit_") {
		leveStr := strings.TrimPrefix(keyStr, "Unit_")
		level, err := strconv.Atoi(leveStr)
		if err != nil {
			panic(err)
		}
		return unitType(level + 1)
	}

	// otherwise
	panic("not_found")
}
