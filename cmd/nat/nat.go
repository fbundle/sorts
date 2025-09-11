package main

import (
	"fmt"

	sorts2 "github.com/fbundle/sorts/obsolete/sorts_v3"
)

func main() {
	n0 := sorts2.Zero
	n1 := sorts2.Succ(n0)
	n2 := sorts2.Succ(n1)
	n3 := sorts2.Succ(n2)
	n4 := sorts2.Succ(n3)
	n5 := sorts2.Succ(n4)
	for _, n := range []sorts2.Sort{n0, n1, n2, n3, n4, n5} {
		fmt.Println(sorts2.Name(n))
	}

	nx := sorts2.Add(n2, n3)
	fmt.Println(sorts2.Name(nx))
}
