package sorts

// every sort is identified with a name (string)

type Sort interface {
	Level(ss SortSystem) int // universe level
	Name(ss SortSystem) string
	Parent(ss SortSystem) Sort
	LessEqual(ss SortSystem, dst Sort) bool // for type casting
}

type SortSystem interface {
	Default(level int) Sort
	AddRule(src string, dst string) SortSystem
	LessEqual(src string, dst string) bool

	Atom(level int, name string, parents ...Sort) Sort
	Arrow(arg Sort, body Sort) Sort
}
