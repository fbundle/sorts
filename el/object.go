package el

import (
	"context"
	"strconv"

	"github.com/fbundle/sorts/expr"
	"github.com/fbundle/sorts/persistent/ordered_map"
	sorts "github.com/fbundle/sorts/sorts/sorts_v3"
)

type Name string

type Frame = ordered_map.OrderedMap[Name, Object]

type Object struct {
	Sort sorts.Sort
	Data Data
}

type Data interface {
	String() string
}
type Function struct {
	name string
	exec func(r Runtime, ctx context.Context, f Frame, e expr.Expr) (Object, error)
}

func (l Function) String() string {
	return l.name
}

type Int struct {
	val int
}

func (i Int) String() string {
	return strconv.Itoa(i.val)
}

type Str struct {
	val string
}

func (s Str) String() string {
	return s.val
}
