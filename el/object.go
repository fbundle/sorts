package el

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

// uType - U_0 U_1 ... U_n
func uType(level int) _object {
	name := fmt.Sprintf("U_%d", level)
	sort := sorts.NewAtom(level, name, nil)
}

type _object struct {
	sort   sorts.Sort // can be nil for partial objects
	next   Expr
	parent func() _object // can be nil for partial objects
}

func (o _object) MustTotal(parent _object) _object {
	// TODO - this is the new type check binding
	panic("not implemented")
}
func (o _object) Total() bool {
	return o.sort != nil
}
func (o _object) Parent() _object {
	return o.parent()
}
