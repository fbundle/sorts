package sorts

import "github.com/fbundle/sorts/adt"

// Sort -
type Sort interface {
	Level() int              // universe level
	Name() string            // every sort is identified with a name (string)
	Parent() InhabitedSort   // (or Type) every sort must have a parent
	LessEqual(dst Sort) bool // for type casting
}

// InhabitedSort - represents a sort with at least one child
// (true theorems have proofs)
type InhabitedSort interface {
	Sort
	Child() Sort // (or Term)
}

// DependentSort - represent a type B(x) depends on sort x
type DependentSort interface {
	Sort
	Apply(Sort) Sort // take x, return B(x)
}
type SortSystem interface {
	Default(level int) Sort
	DefaultInhabited(child Sort) InhabitedSort
	AddRule(src string, dst string) SortSystem
	LessEqual(src string, dst string) bool

	Atom(level int, name string, parent InhabitedSort) adt.Option[Sort]
	Arrow(arg Sort, body Sort) adt.Option[Sort]
	Inhabited(sort Sort, elem Sort) adt.Option[InhabitedSort]
	Pi(arg InhabitedSort, body DependentSort) adt.Option[Sort]
	Sigma(a InhabitedSort, b DependentSort) adt.Option[Sort]
	Sum(a Sort, b Sort) adt.Option[Sort]
	Prod(a Sort, b Sort) adt.Option[Sort]
	NewDependent(level int, name string, parent InhabitedSort, apply func(Sort) Sort) adt.Option[DependentSort]
}
