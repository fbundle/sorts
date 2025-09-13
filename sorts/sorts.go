package sorts

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

var typeErr = errors.New("type_error")

type MustParseFunc = func(form.Form) Sort

type mustParseListFunc = func(MustParseFunc, form.List) Sort

var listParsers = map[form.Name]mustParseListFunc{
	ArrowName: mustParseArrow,
	ProdName:  mustParseProd,
	SumName:   mustParseSum,
}

func Repr(s any) form.Form {
	if s == nil {
		return form.List{}
	}
	if s, ok := s.(Sort); ok {
		return s.sortAttr().repr
	}
	if s, ok := s.(Dependent); ok {
		return s.Repr
	}
	panic(typeErr)
}

func Level(s Sort) int {
	return s.sortAttr().level
}

func Parent(s Sort) Sort {
	return s.sortAttr().parent
}

func SubTypeOf(x Sort, y Sort) bool {
	if Level(x) == Level(y) && (Repr(x) == InitialName || Repr(y) == TerminalName) {
		return true
	}

	return x.sortAttr().lessEqual(y)
}
func TermOf(x Sort, X Sort) bool {
	X1 := Parent(x)
	return SubTypeOf(X1, X)
}

type Sort interface {
	sortAttr() sortAttr
}

type sortAttr struct {
	repr      form.Form           // every Sort is identified with a Repr
	level     int                 // universe Level
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}
