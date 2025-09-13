package el

import (
	"fmt"

	"github.com/fbundle/sorts/sorts"
)

func newTerm(expr Expr, parent Object) Object {
	return Object{
		next: expr,
		sort: sorts.newTerm(String(expr), parent.sort),
		parent: func() Object {
			return parent
		},
	}
}

func newType(level int, name Term) Object {
	parent := uType(level + 1)
	return Object{
		next: name,
		sort: sorts.newTerm(string(name), parent.sort),
		parent: func() Object {
			return parent
		},
	}
}

func anyType(level int) Object {
	return newType(level, sorts.TerminalName)
}

func unitType(level int) Object {
	return newType(level, sorts.InitialName)
}

// uType - U_0 U_1 ... U_n
func uType(level int) Object {
	name := "U"
	return Object{
		next: Term(fmt.Sprintf("U_%d", level)),
		sort: sorts.newAtom(level, name, name),
		parent: func() Object {
			return uType(level + 1)
		},
	}
}

func (o Object) mustTotal(parent Object) Object {
	// TODO - this is the new type check binding
	panic("not implemented")
}
func (o Object) isTotal() bool {
	return o.sort != nil
}

type Object struct {
	next   Expr
	sort   sorts.Sort    // can be nil for partial objects
	parent func() Object // can be nil for partial objects
}

func (o Object) Expr() Expr {
	return o.next
}
