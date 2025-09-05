package sorts

import "strings"

const (
	arrowToken = "->"
)

type arrow struct {
	param Sort
	body  Sort
}

func (s arrow) Level(ss SortSystem) int {
	return max(s.param.Level(ss), s.body.Level(ss))
}

func (s arrow) Name(ss SortSystem) string {
	return strings.Join([]string{
		s.param.Name(ss),
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
		// reverse cast for param
		// {any -> unit} can be cast into {int -> unit}
		// because {int} can be cast into {any}
		if !d.param.LessEqual(ss, s.param) {
			return false
		}
		// normal cast for body
		return s.body.LessEqual(ss, d.body)
	default:
		panic("unreachable")
	}
}
