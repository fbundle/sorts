package sorts

import (
	"errors"
)

var typeErr = errors.New("type_error")

type Node[T any] struct {
	Value    T
	Children []Node[T]
}

type term = comparable // leaf type constraint for sorts

type Sort[T term] interface {
	sortAttr() sortAttr[T]
}

type Universe[T term] struct {
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

func (u *Universe[T]) Universe(level int) Sort[T] {
	return newAtomChain[T](level, u.UniverseName)
}

func (u *Universe[T]) NewTerm(name T, parent Sort[T]) Sort[T] {
	return newAtomTerm(name, parent)
}

func (u *Universe[T]) Repr(s any) Node[T] {
	if sort, ok := s.(Sort[T]); ok {
		return sort.sortAttr().repr
	}
	if dep, ok := s.(Dependent[T]); ok {
		return dep.Repr
	}
	panic(typeErr)
}

func (u *Universe[T]) Level(s Sort[T]) int {
	return s.sortAttr().level
}
func (u *Universe[T]) Parent(s Sort[T]) Sort[T] {
	return s.sortAttr().parent
}
func (u *Universe[T]) SubTypeOf(x Sort[T], y Sort[T]) bool {
	return x.sortAttr().lessEqual(u, y)
}
func (u *Universe[T]) TermOf(x Sort[T], X Sort[T]) bool {
	return u.SubTypeOf(u.Parent(x), X)
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

type sortAttr[T term] struct {
	repr      Node[T]                                // every Sort is identified with a Repr
	level     int                                    // universe Level
	parent    Sort[T]                                // (or Type) every Sort must have a Parent
	lessEqual func(u *Universe[T], dst Sort[T]) bool // a partial order on sorts (subtype)
}
