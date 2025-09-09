package nat

import (
	"strconv"

	sorts "github.com/fbundle/sorts/sorts_v3"
)

var Nat = sorts.NewAtom(0, "Nat", nil)

var Zero = sorts.NewAtom(-1, "0", Nat)

var NatToNat = sorts.Arrow{A: Nat, B: Nat}
var NatToNatToNat = sorts.Arrow{A: Nat, B: NatToNat}

func Succ(n sorts.Sort) sorts.Sort {
	return NatToNat.Elim(SuccArrow, n)
}

var SuccArrow = NatToNat.Intro("succ", func(a sorts.Sort) sorts.Sort {
	sorts.MustTermOf(a, Nat)
	aVal, err := strconv.Atoi(sorts.Name(a))
	if err != nil {
		return sorts.NewAtom(-1, err.Error(), Nat)
	}
	return sorts.NewAtom(-1, strconv.Itoa(aVal+1), Nat)
})

func Add(n1 sorts.Sort, n2 sorts.Sort) sorts.Sort {
	addN1 := NatToNatToNat.Elim(AddArrow, n1)
	return NatToNat.Elim(addN1, n2)
}

var AddArrow = NatToNatToNat.Intro("add", func(a sorts.Sort) sorts.Sort {
	sorts.MustTermOf(a, Nat)
	return NatToNat.Intro("add_a", func(b sorts.Sort) sorts.Sort {
		sorts.MustTermOf(b, Nat)
		return sorts.NewAtom(-1, "add_a_b", Nat)
	})
})
