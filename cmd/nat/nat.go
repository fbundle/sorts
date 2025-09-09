package main

import (
	"fmt"

	"github.com/fbundle/sorts/nat"
	sorts "github.com/fbundle/sorts/sorts_v3"
)

func main() {
	n0 := nat.Zero
	n1 := sorts.Arrow{nat.Nat, nat.Nat}.Elim(nat.Succ, n0)
	n2 := sorts.Arrow{nat.Nat, nat.Nat}.Elim(nat.Succ, n1)
	n3 := sorts.Arrow{nat.Nat, nat.Nat}.Elim(nat.Succ, n2)
	for _, n := range []sorts.Sort{n0, n1, n2, n3} {
		fmt.Println(sorts.Name(n))
	}
}
