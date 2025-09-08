package nat

import (
	"strconv"

	"github.com/fbundle/sorts/sorts_v3"
)

var Nat = sorts.NewAtom(0, "Nat", nil)

var Zero = sorts.NewAtom(-1, "0", Nat)

var succSort = sorts.NewAtom(-1, "succ", sorts.Arrow{Nat, Nat})

var Succ = sorts.Inhabited{
	Sort: succSort,
	Child: sorts.Arrow{Nat, Nat}.Intro(func(s sorts.Sort) sorts.Sort {
		sorts.MustTermOf(s, Nat)
		n, err := strconv.Atoi(sorts.Name(s))
		if err != nil {
			panic(err)
		}
		return sorts.NewAtom(1, strconv.Itoa(n+1), Nat)
	}),
}
