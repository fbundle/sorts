package el

import (
	sorts "github.com/fbundle/sorts/obsolete/sorts_v3"
	"github.com/fbundle/sorts/persistent/ordered_map"
)

type Value struct {
	Sort sorts.Sort
	AST  AST
}

type Frame = ordered_map.OrderedMap[Term, Value]
