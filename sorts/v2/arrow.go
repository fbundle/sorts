package sorts

import "strings"

const (
	arrowToken = "->"
)

type arrow struct {
	arg  Sort
	body Sort
}

func (s arrow) Level(ss SortSystem) int {
	return max(s.arg.Level(ss), s.body.Level(ss))
}

func (s arrow) Name(ss SortSystem) string {
	return strings.Join([]string{
		s.arg.Name(ss),
		arrowToken,
		s.body.Name(ss),
	}, " ")
}

func (s arrow) Parent(ss SortSystem) Sort {
	return ss.Default(s.Level(ss) + 1)
}

func (s arrow) LessEqual(ss SortSystem, dst Sort) bool {
	switch d := dst.(type) {
	case atom:
		return false // cannot compare arrow to atom
	case arrow:
		// reverse cast for arg
		// {any -> unit} can be cast into {int -> unit}
		// because {int} can be cast into {any}
		if !d.arg.LessEqual(ss, s.arg) {
			return false
		}
		// normal cast for body
		return s.body.LessEqual(ss, d.body)
	default:
		panic("unreachable")
	}
}
