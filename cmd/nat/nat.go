package main

import (
	"fmt"

	"github.com/fbundle/sorts/nat"
	sorts "github.com/fbundle/sorts/sorts_v3"
)

func main() {
	n0 := nat.Zero
	n1 := nat.Succ(n0)
	n2 := nat.Succ(n1)
	n3 := nat.Succ(n2)
	n4 := nat.Succ(n3)
	n5 := nat.Succ(n4)
	for _, n := range []sorts.Sort{n0, n1, n2, n3, n4, n5} {
		fmt.Println(sorts.Name(n))
	}

	nx := nat.Add(n2, n3)
	fmt.Println(sorts.Name(nx))
}
