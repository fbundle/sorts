package nat

import (
	"strconv"

	"github.com/fbundle/sorts/sorts_v3"
)

type Nat struct {
	Zero sorts.Sort
	Succ func(sorts.Sort) sorts.Sort
}

func NewNat() Nat {
	natType := sorts.NewAtom(2, "Nat", nil)
	return Nat{
		Zero: sorts.NewAtom(1, "0", natType),
		Succ: func(s sorts.Sort) sorts.Sort {
			n, err := strconv.Atoi(sorts.Name(s))
			if err != nil {
				panic(err)
			}
			return sorts.NewAtom(1, strconv.Itoa(n+1), natType)
		},
	}
}
