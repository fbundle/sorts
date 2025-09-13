package sorts

import "fmt"

// Pi - (x: A) -> (y: B(x)) similar to Arrow
// this is the universal quantifier
type Pi struct {
	A Sort
	B Dependent
}

func (s Pi) sortAttr() sortAttr {
	x := makeTerm("x", s.A) // dummy term
	sBx := s.B.Apply(x)
	level := max(Level(s.A), Level(sBx))
	return sortAttr{
		name:   fmt.Sprintf("Π(x:%s)%s(x)", Name(s.A), Name(s.B)),
		level:  level,
		parent: makeAtom(level+1, "Type", "Type"),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Pi:
				// tricky: subtyping for Arrow is contravariant in domain, covariant in codomain
				if !SubTypeOf(d.A, s.A) {
					return false
				}
				y := makeTerm("y", d.A)
				dBy := d.B.Apply(y)
				return SubTypeOf(sBx, dBy)
			default:
				return false
			}
		},
	}
}

// Intro - take a func that maps (a: A) into (b: B(a))  give (f: Π(x:A)B(x))
func (s Pi) Intro(name string, arrow func(Sort) Sort) Sort {
	// verify
	a := makeTerm("a", s.A)
	b := arrow(a)
	sBa := s.B.Apply(a)
	mustTermOf(b, sBa) // TODO - think, shouldn't it have to check for every a of type A?

	return makeTerm(name, s)
}

// Elim - take (f: Π(x:A)B(x)) (a: A) give (b: B(a)) - Modus Ponens
func (s Pi) Elim(arrow Sort, a Sort) Sort {
	mustTermOf(arrow, s)
	mustTermOf(a, s.A)
	Ba := s.B.Apply(a)

	name := fmt.Sprintf("(%s %s)", Name(arrow), Name(a))
	return makeTerm(name, Ba)
}
