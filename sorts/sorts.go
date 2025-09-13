package sorts

import (
	"errors"
	"strconv"

	"github.com/fbundle/sorts/form"
)

var typeErr = errors.New("type_error")

type Node[T any] struct {
	Value    T
	Children []Node[T]
}

type Sort[T comparable] interface {
	sortAttr() sortAttr[T]
}

type Universe[T comparable] struct {
	UniverseName func(level int) T
	InitialName  T
	TerminalName T
	lessEqualMap map[[2]T]struct{}
}

func (u *Universe[T]) AddRule(src T, dst T) {
	if u.lessEqualMap == nil {
		u.lessEqualMap = make(map[[2]T]struct{})
	}
	u.lessEqualMap[[2]T{src, dst}] = struct{}{}
}

func (u *Universe[T]) lessEqual(src T, dst T) bool {
	if src == u.InitialName || dst == u.TerminalName {
		return true
	}
	if src == dst {
		return true
	}
	if u.lessEqualMap != nil {
		if _, ok := u.lessEqualMap[[2]T{src, dst}]; ok {
			return true
		}
	}
	return false
}

func (u *Universe[T]) Universe(level int) Atom[T] {

}

// universe - ... U_{-1}, U_0, U_1, ...
func universe(level int) Sort {
	return newAtomChain(level, func(i int) form.Name {
		levelStr := strconv.Itoa(level)
		if level < 0 {
			return form.Name("U_{" + levelStr + "}")
		} else {
			return form.Name("U_" + levelStr)
		}
	})
}

//

type mustParseFunc = func(form.Form) Sort

type mustParseListFunc = func(mustParseFunc, form.List) Sort

var listParsers = map[form.Name]mustParseListFunc{
	ArrowName: mustParseArrow,
	ProdName:  mustParseProd,
	SumName:   mustParseSum,
	// TODO - fill all types
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

type sortAttr[T comparable] struct {
	repr      Node[T]                                // every Sort is identified with a Repr
	level     int                                    // universe Level
	parent    Sort[T]                                // (or Type) every Sort must have a Parent
	lessEqual func(u *Universe[T], dst Sort[T]) bool // a partial order on sorts (subtype)
}
