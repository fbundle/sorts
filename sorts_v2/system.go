package sorts

import (
	"strconv"

	"github.com/fbundle/sorts/adt"
)

type SortSystemOption func(*sortSystem)

func WithInitialName(name string) SortSystemOption {
	return func(ss *sortSystem) {
		ss.initialName = name
	}
}

func WithTerminalName(name string) SortSystemOption {
	return func(ss *sortSystem) {
		ss.terminalName = name
	}
}

func NewSortSystem(defaultName string, opts ...SortSystemOption) SortSystem {
	ss := &sortSystem{
		initialName:  "",
		terminalName: "",
		defaultName:  defaultName,
		lessEqualMap: make(map[rule]struct{}),
	}
	for _, opt := range opts {
		opt(ss)
	}
	return ss.validate()
}

type rule struct {
	src string
	dst string
}
type sortSystem struct {
	initialName  string // empty = noset
	terminalName string // empty = noset
	defaultName  string // must be non-empty, used when no parent is set
	lessEqualMap map[rule]struct{}
}

func (ss *sortSystem) LessEqual(src string, dst string) bool {
	if ss.isInitial(src) {
		return true
	}
	if ss.isTerminal(dst) {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := ss.lessEqualMap[rule{src: src, dst: dst}]; ok {
		return true
	}
	return false
}
func (ss *sortSystem) DefaultInhabited(child Sort) InhabitedSort {
	level := child.Level() + 1
	name := ss.defaultName + "_" + strconv.Itoa(level)
	sort := newAtom(ss, level, name, nil).MustUnwrap()
	return newInhabited(ss, sort, child).MustUnwrap()
}

func (ss *sortSystem) Default(level int) Sort {
	name := ss.defaultName + "_" + strconv.Itoa(level)
	return newAtom(ss, level, name, nil).MustUnwrap()
}

func (ss *sortSystem) Atom(level int, name string, parent Sort) adt.Option[Sort] {
	return newAtom(ss, level, name, parent)
}

func (ss *sortSystem) Arrow(arg Sort, body Sort) adt.Option[Sort] {
	return newArrow(ss, arg, body)
}

func (ss *sortSystem) Inhabited(sort Sort, child Sort) adt.Option[InhabitedSort] {
	return newInhabited(ss, sort, child)
}

func (ss *sortSystem) Pi(arg InhabitedSort, body DependentSort) adt.Option[Sort] {
	return newPi(ss, arg, body)
}

func (ss *sortSystem) Sigma(a InhabitedSort, b DependentSort) adt.Option[Sort] {
	return newSigma(ss, a, b)
}

func (ss *sortSystem) Sum(a Sort, b Sort) adt.Option[Sort] {
	return newSum(ss, a, b)
}
func (ss *sortSystem) Prod(a Sort, b Sort) adt.Option[Sort] {
	return newProd(ss, a, b)
}
func (ss *sortSystem) NewDependent(level int, name string, parent InhabitedSort, apply func(Sort) Sort) adt.Option[DependentSort] {
	return newDependent(ss, level, name, parent, apply)
}

func (ss *sortSystem) AddRule(src string, dst string) SortSystem {
	ss.lessEqualMap[rule{src: src, dst: dst}] = struct{}{}
	return ss
}

func (ss *sortSystem) isInitial(name string) bool {
	return len(ss.initialName) > 0 && ss.initialName == name
}
func (ss *sortSystem) isTerminal(name string) bool {
	return len(ss.terminalName) > 0 && ss.terminalName == name
}

func (ss *sortSystem) validate() *sortSystem {
	if len(ss.defaultName) == 0 {
		panic("validate_error: default name must not be empty")
	}
	if len(ss.initialName) > 0 && len(ss.terminalName) > 0 && ss.initialName == ss.terminalName {
		panic("validate_error: initial and terminal must be distinct")
	}
	return ss
}
