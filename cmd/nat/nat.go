package main

import (
	"fmt"

	"github.com/fbundle/sorts/nat"
	sorts "github.com/fbundle/sorts/sorts_v3"
)

func main() {
	var n sorts.Sort = nat.Zero
	fmt.Println(sorts.Name(n))
	for i := 0; i < 5; i++ {
		n = nat.Succ(n)
		fmt.Println(sorts.Name(n))
	}
}
