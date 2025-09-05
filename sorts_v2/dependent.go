package sorts

import (
	"fmt"

	"github.com/fbundle/lab_public/lab/go_util/pkg/adt"
)

func newDependent(ss SortSystem, level int, name string, parent InhabitedSort, apply func(Sort) Sort) adt.Option[DependentSort] {
	if parent != nil && parent.Sort().Level() < level {
		// TODO - chat GPT said parent of dependent can be the same level as level
		return adt.None[DependentSort]()
	}
	return adt.Some[DependentSort](dependent{
		level:  level,
		name:   name,
		parent: parent,
		apply:  apply,
		ss:     ss,
	})
}

type dependent struct {
	level  int
	name   string
	parent InhabitedSort
	apply  func(Sort) Sort
	ss     SortSystem
}

func (s dependent) Sort() Sort {
	return s
}

func (s dependent) Level() int {
	return s.level
}

func (s dependent) Name() string {
	return s.name
}

func (s dependent) Parent() InhabitedSort {
	if s.parent != nil {
		return s.parent
	}
	// default parent
	return s.ss.DefaultInhabited(s)
}

func (s dependent) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s dependent) Apply(sort Sort) Sort {
	// TODO - check
	out := s.apply(sort)
	if out.Level() != s.level {
		panic("type_error")
	}
	if out.Name() != fmt.Sprintf("%s(%s)", s.name, sort.Name()) {
		panic("type_error")
	}
	return out
}
