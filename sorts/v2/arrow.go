package sorts

import (
	"strings"

	"github.com/fbundle/sorts/adt"
)

const (
	arrowToken = "->"
)

func newArrow(ss SortSystem, arg Sort, body Sort) adt.Option[Sort] {
	return adt.Some[Sort](arrow{
		arg:  arg,
		body: body,
		ss:   ss,
	})
}

type arrow struct {
	arg  Sort
	body Sort
	ss   SortSystem
}

func (s arrow) Level() int {
	return max(s.arg.Level(), s.body.Level())
}

func (s arrow) Name() string {
	return strings.Join([]string{
		s.arg.Name(),
		arrowToken,
		s.body.Name(),
	}, " ")
}

func (s arrow) Parent() InhabitedSort {
	// default parent
	return s.ss.DefaultInhabited(s)
}

func (s arrow) LessEqual(dst Sort) bool {
	switch d := dst.(type) {
	case atom:
		return false // cannot compare arrow to atom
	case arrow:
		// reverse cast for arg
		// {any -> unit} can be cast into {int -> unit}
		// because {int} can be cast into {any}
		if !d.arg.LessEqual(s.arg) {
			return false
		}
		// normal cast for body
		return s.body.LessEqual(d.body)
	default:
		panic("unreachable")
	}
}
