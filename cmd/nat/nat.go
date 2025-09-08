package main

import (
	"fmt"

	"github.com/fbundle/sorts/nat"
	sorts "github.com/fbundle/sorts/sorts_v3"
)

func main() {

	n0 := nat.Zero
	n1 := sorts.Arrow{nat.Nat, nat.Nat}.Elim(nat.Succ, n0)
	fmt.Println(sorts.Name(n0), sorts.Name(n1))
}
