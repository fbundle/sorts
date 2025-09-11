package el

import (
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Value struct {
	Sort sorts.Sort
	Data AST
}

type Frame = ordered_map.OrderedMap[Term, Value]
