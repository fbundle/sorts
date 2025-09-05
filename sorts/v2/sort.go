package sorts

import "github.com/fbundle/sorts/adt"

// every sort is identified with a name (string)

type Sort interface {
	Level(ss SortSystem) int // universe level
	Name(ss SortSystem) string
	Parent(ss SortSystem) Sort
	LessEqual(ss SortSystem, dst Sort) bool // for type casting
}

type SortSystem interface {
	DefaultName() string
	AddRule(src string, dst string) SortSystem
	LessEqual(src string, dst string) bool

	Atom(level int, name string, parent Sort) adt.Option[Sort]
	Arrow(param Sort, body Sort) adt.Option[Sort]
}
