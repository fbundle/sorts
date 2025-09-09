package main

import (
	"fmt"

	sorts "github.com/fbundle/sorts/sorts_v3"
)

func main() {
	n0 := sorts.Zero
	n1 := sorts.Succ(n0)
	n2 := sorts.Succ(n1)
	n3 := sorts.Succ(n2)
	n4 := sorts.Succ(n3)
	n5 := sorts.Succ(n4)
	for _, n := range []sorts.Sort{n0, n1, n2, n3, n4, n5} {
		fmt.Println(sorts.Name(n))
	}

	nx := sorts.Add(n2, n3)
	fmt.Println(sorts.Name(nx))
}
