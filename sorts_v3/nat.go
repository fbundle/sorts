package sorts

import (
	"fmt"
	"strconv"
)

const (
	NatLevel = 0
)

var Nat = NewAtom(NatLevel+1, "Nat", nil)

var Zero = NewAtom(NatLevel, "0", Nat)

var NatToNat = Arrow{A: Nat, B: Nat}
var NatToNatToNat = Arrow{A: Nat, B: NatToNat}

func Succ(n Sort) Sort {
	return NatToNat.Elim(SuccArrow, n)
}

var SuccArrow = NatToNat.Intro("succ", func(a Sort) Sort {
	MustTermOf(a, Nat)
	aVal, err := strconv.Atoi(Name(a))
	if err != nil {
		return NewAtom(NatLevel, err.Error(), Nat)
	}
	return NewAtom(NatLevel, strconv.Itoa(aVal+1), Nat)
})

func Add(n1 Sort, n2 Sort) Sort {
	addN1 := NatToNatToNat.Elim(AddArrow, n1)
	return NatToNat.Elim(addN1, n2)
}

var AddArrow = NatToNatToNat.Intro("add", func(a Sort) Sort {
	MustTermOf(a, Nat)
	return NatToNat.Intro("add_a", func(b Sort) Sort {
		MustTermOf(b, Nat)
		return NewAtom(NatLevel, "add_a_b", Nat)
	})
})

type Equal struct {
	A Sort
	B Sort
}

func (s Equal) sortAttr() sortAttr {
	level := max(Level(s.A), Level(s.B))
	return sortAttr{
		name:   fmt.Sprintf("Eq(%s, %s)", Name(s.A), Name(s.B)),
		level:  level,
		parent: defaultSort(nil, level+1),
		lessEqual: func(dst Sort) bool {
			switch d := dst.(type) {
			case Equal:
				return LessEqual(s.A, d.A) && LessEqual(s.B, d.B)
			default:
				return false
			}
		},
	}
}

func (s Equal) Refl(x Sort) Sort {
	MustTermOf(x, s.A)
	MustTermOf(x, s.B)
	return dummyTerm(s, fmt.Sprintf("%s = %s", Name(x), Name(x)))
}
