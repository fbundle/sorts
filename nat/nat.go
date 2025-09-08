package nat

import (
	"strconv"

	"github.com/fbundle/sorts/sorts_v3"
)

var Nat = sorts.NewAtom(0, "Nat", nil)

var Zero = sorts.NewAtom(-1, "0", Nat)

var Succ = sorts.Arrow{Nat, Nat}.Intro(func(a sorts.Sort) sorts.Sort {
	sorts.MustTermOf(a, Nat)
	aVal, err := strconv.Atoi(sorts.Name(a))
	if err != nil {
		return sorts.NewAtom(-1, err.Error(), Nat)
	}
	return sorts.NewAtom(-1, strconv.Itoa(aVal+1), Nat)
})

var Add = sorts.Arrow{Nat, sorts.Arrow{Nat, Nat}}.Intro(func(a sorts.Sort) sorts.Sort {
	sorts.MustTermOf(a, Nat)
	aVal, err := strconv.Atoi(sorts.Name(a))
	if err != nil {
		return sorts.NewAtom(-1, err.Error(), sorts.Arrow{Nat, Nat})
	}

	return sorts.Arrow{Nat, Nat}.Intro(func(b sorts.Sort) sorts.Sort {
		sorts.MustTermOf(b, Nat)
		bVal, err := strconv.Atoi(sorts.Name(b))
		if err != nil {
			return sorts.NewAtom(-1, err.Error(), Nat)
		}
		return sorts.NewAtom(-1, strconv.Itoa(aVal+bVal), Nat)
	})
})
