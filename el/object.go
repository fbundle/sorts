package el

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

func newTerm(expr Expr, parent _object) _object {
	return _object{
		next: expr,
		sort: sorts.NewTerm(String(expr), parent.sort),
		parent: func() _object {
			return parent
		},
	}
}

func newType(level int, name Term) _object {
	parent := uType(level + 1)
	return _object{
		next: name,
		sort: sorts.NewTerm(string(name), parent.sort),
		parent: func() _object {
			return parent
		},
	}
}

func anyType(level int) _object {
	return newType(level, sorts.TerminalName)
}

func unitType(level int) _object {
	return newType(level, sorts.InitialName)
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
