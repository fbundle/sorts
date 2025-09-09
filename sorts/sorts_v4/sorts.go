package sorts

import "github.com/fbundle/sorts/expr"

const (
	defaultTerm  = "type"
	InitialTerm  = "unit" // initial object
	TerminalTerm = "any"  // terminal object
)

type Sort interface {
	sortAttr() sortAttr
}

type sortAttr struct {
	repr      expr.Expr           // every Sort is identified with an expression
	level     int                 // universe Level
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}

func Repr(s Sort) expr.Expr {
	return s.sortAttr().repr
}

func Level(s Sort) int {
	return s.sortAttr().level
}

func Parent(s Sort) Sort {
	return s.sortAttr().parent
}

func LessEqual(x Sort, y Sort) bool {
	return x.sortAttr().lessEqual(y)
}

func Parse(e expr.Expr) Sort {
	panic("not_implemented")
}
