package el

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

func newType(level int, name string) _object {
	parent := uType(level + 1)
	return _object{
		next: Term(name),
		sort: sorts.NewTerm(name, parent.sort),
		parent: func() _object {
			return parent
		},
	}
}

// uType - U_0 U_1 ... U_n
func uType(level int) _object {
	name := "U"
	return _object{
		next: Term(fmt.Sprintf("U_%d", level)),
		sort: sorts.NewAtom(level, name, name),
		parent: func() _object {
			return uType(level + 1)
		},
	}
}

func (o _object) mustTotal(parent _object) _object {
	// TODO - this is the new type check binding
	panic("not implemented")
}
func (o _object) isTotal() bool {
	return o.sort != nil
}

type _object struct {
	next   Expr
	sort   sorts.Sort     // can be nil for partial objects
	parent func() _object // can be nil for partial objects
}
