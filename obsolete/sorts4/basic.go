package sorts

import (
	"fmt"
)

func ListCompileArrow(H Name) ListCompileFunc {
	return func(ctx Context, list List) Sort {
		if len(list) != 3 {
			panic(fmt.Errorf("arrow must be %s domain codomain", H))
		}
		if list[0] != H {
			panic(fmt.Errorf("arrow must be %s domain codomain", H))
		}
		return Arrow{H: H, A: ctx.Compile(list[1]), B: ctx.Compile(list[2])}
	}
}

type Arrow struct {
	H Name
	A Sort
	B Sort
}

func (s Arrow) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, a.Form(s.A), a.Form(s.B)},
		level:  max(a.Level(s.A), a.Level(s.B)),
		parent: Arrow{H: s.H, A: a.Parent(s.A), B: a.Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Arrow:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				// {any -> unit} can be cast into {int -> unit}
				// because {int} can be cast into {any}
				if !a.LessEqual(d.A, s.A) {
					return false
				}
				return a.LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Elim - take (f: A -> B) (a: A) give (b: B) - Modus Ponens
func (s Arrow) Elim(sa SortAttr, f Sort, a Sort) Sort {
	must(sa).termOf(f, s)
	must(sa).termOf(a, s.A)

	// TODO - not hard code this
	return NewAtomTerm(sa, List{Name("elim"), sa.Form(f), sa.Form(a)}, s.B)
}

// Intro - take a func that maps (a: A) into (b: B)  give (f: A -> B)
func (s Arrow) Intro(sa SortAttr, name Name, f func(Sort) Sort) Sort {
	// verify
	a := NewAtomTerm(sa, Name("a"), s.A)
	b := f(a)

	must(sa).termOf(b, s.B)

	// verify ok

	return NewAtomTerm(sa, name, s)
}
func ListCompileProd(H Name) ListCompileFunc {
	return func(ctx Context, list List) Sort {
		if len(list) != 3 {
			panic(fmt.Errorf("prod must be %s A B", H))
		}
		if list[0] != H {
			panic(fmt.Errorf("prod must be %s A B", H))
		}
		return Prod{H: H, A: ctx.Compile(list[1]), B: ctx.Compile(list[2])}
	}
}

type Prod struct {
	H Name
	A Sort
	B Sort
}

func (s Prod) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, a.Form(s.A), a.Form(s.B)},
		level:  max(a.Level(s.A), a.Level(s.B)),
		parent: Prod{A: a.Parent(s.A), B: a.Parent(s.B)}, // smallest type containing A × B
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Prod:
				return a.LessEqual(s.A, d.A) && a.LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro - take (a: A) (b: B) give (a, b): A × B
func (s Prod) Intro(sa SortAttr, a Sort, b Sort) Sort {
	must(sa).termOf(a, s.A)
	must(sa).termOf(b, s.B)
	return Prod{A: a, B: b}
}

// Elim - take (t: A × B) give (a: A) and (b: B)
func (s Prod) Elim(sa SortAttr, t Sort) (left Sort, right Sort) {
	must(sa).termOf(t, s)
	if t, ok := t.(Prod); ok {
		return t.A, t.B
	}

	a := NewAtomTerm(sa, List{Name("left"), sa.Form(t)}, s.A)
	b := NewAtomTerm(sa, List{Name("right"), sa.Form(t)}, s.B)
	return a, b
}
func ListCompileSum(H Name) ListCompileFunc {
	return func(ctx Context, list List) Sort {
		if len(list) != 3 {
			panic(fmt.Errorf("sum must be %s A B", H))
		}
		if list[0] != H {
			panic(fmt.Errorf("sum must be %s A B", H))
		}
		return Sum{H: H, A: ctx.Compile(list[1]), B: ctx.Compile(list[2])}
	}
}

type Sum struct {
	H Name
	A Sort
	B Sort
}

func (s Sum) sortAttr(sa SortAttr) sortAttr {
	return sortAttr{
		form:   List{s.H, sa.Form(s.A), sa.Form(s.B)},
		level:  max(sa.Level(s.A), sa.Level(s.B)),
		parent: Sum{A: sa.Parent(s.A), B: sa.Parent(s.B)},
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Sum:
				return sa.LessEqual(s.A, d.A) && sa.LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

// Intro - IntroLeft or IntroRight
func (s Sum) Intro(sa SortAttr, a Sort, b Sort) Sort {
	if a != nil {
		// IntroLeft - take (a: A) give (x: A + B)
		must(sa).termOf(a, s.A)
		return a
	} else {
		// IntroRight - take (b: B) give (x: A + B)
		must(sa).termOf(b, s.B)
		return b
	}
}

// ByCases - take (t: A + B) (h1: A -> X) (h2: B -> X) give (x: X)
func (s Sum) ByCases(sa SortAttr, t Sort, h1 Sort, h2 Sort) Sort {
	must(sa).termOf(t, s)
	X := sa.Parent(h1).(Arrow).B
	must(sa).termOf(h1, Arrow{s.H, s.A, X})
	must(sa).termOf(h2, Arrow{s.H, s.B, X})

	return NewAtomTerm(sa, List{Name("by_cases"), sa.Form(t), sa.Form(h1), sa.Form(h2)}, X)
}
