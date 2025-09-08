package nat

import "github.com/fbundle/sorts/sorts_v3"

type Nat struct {
	Zero sorts.Sort
	Succ func(sorts.Sort) sorts.Sort
}

func NewNat() Nat {
	natType := sorts.NewAtom()
	return Nat{
		Zero: sorts.
	}
}
