package sorts

import "github.com/fbundle/sorts/adt"

// every sort is identified with a name (string)

type Sort interface {
	Level() int // universe level
	Name() string
	Parent() InhabitedSort
	LessEqual(dst Sort) bool // for type casting
}

type InhabitedSort interface {
	Sort
	Child() Sort // give a sort of one level down
}

type SortSystem interface {
	DefaultInhabited(child Sort) InhabitedSort
	AddRule(src string, dst string) SortSystem
	LessEqual(src string, dst string) bool

	Atom(level int, name string, parent InhabitedSort) adt.Option[Sort]
	Arrow(arg Sort, body Sort) adt.Option[Sort]
	Inhabited(sort Sort, elem Sort) adt.Option[InhabitedSort]
}
