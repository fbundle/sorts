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

type Universe[T term] interface {
	U(level int) Sort[T]
	NewTerm(name T, parent Sort[T]) Sort[T]

	Repr(s any) Node[T]
	Level(s Sort[T]) int
	Parent(s Sort[T]) Sort[T]
	SubTypeOf(x Sort[T], y Sort[T]) bool
	TermOf(x Sort[T], X Sort[T]) bool

	AddRule(src T, dst T)
	lessEqual(src T, dst T) bool
}

func NewUniverse[T term](universeName func(level int) T, initialName T, terminalName T) Universe[T] {
	if initialName == terminalName {
		panic(typeErr)
	}
	return &universe[T]{
		universeName: universeName,
		initialName:  initialName,
		terminalName: terminalName,
	}
}

type universe[T term] struct {
	universeName func(level int) T
	initialName  T
	terminalName T
	lessEqualMap map[[2]T]struct{}
}

func (u *universe[T]) AddRule(src T, dst T) {
	if u.lessEqualMap == nil {
		u.lessEqualMap = make(map[[2]T]struct{})
	}
	u.lessEqualMap[[2]T{src, dst}] = struct{}{}
}

func (u *universe[T]) U(level int) Sort[T] {
	return newAtomChain[T](level, u.universeName)
}

func (u *universe[T]) NewTerm(name T, parent Sort[T]) Sort[T] {
	return newAtomTerm(u, name, parent)
}

func (u *universe[T]) Repr(s any) Node[T] {
	if sort, ok := s.(Sort[T]); ok {
		return sort.sortAttr().repr
	}
	if dep, ok := s.(Dependent[T]); ok {
		return dep.Repr
	}
	panic(typeErr)
}

func (u *universe[T]) Level(s Sort[T]) int {
	return s.sortAttr().level
}
func (u *universe[T]) Parent(s Sort[T]) Sort[T] {
	return s.sortAttr().parent
}
func (u *universe[T]) SubTypeOf(x Sort[T], y Sort[T]) bool {
	return x.sortAttr().lessEqual(u, y)
}
func (u *universe[T]) TermOf(x Sort[T], X Sort[T]) bool {
	return u.SubTypeOf(u.Parent(x), X)
}

// private

func (u *universe[T]) lessEqual(src T, dst T) bool {
	if src == u.initialName || dst == u.terminalName {
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
	repr      Node[T]                               // every Sort is identified with a Repr
	level     int                                   // universe Level
	parent    Sort[T]                               // (or Type) every Sort must have a Parent
	lessEqual func(u Universe[T], dst Sort[T]) bool // a partial order on sorts (subtype)
}
