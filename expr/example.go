package expr

import (
	"fmt"
)

// Example demonstrates tokenization, parsing (prefix and infix), strings, nested calls, and marshaling.
func Example() {
	// Tokenize with comments and commas inside infix blocks
	src1 := "{a, b, c} # a list"
	toks1 := Tokenize(src1)
	fmt.Println(toks1)

	// Parse a simple prefix form
	src2 := "(f x y)"
	expr2, rest2, err := Parse(Tokenize(src2))
	if err != nil || len(rest2) != 0 {
		panic("unexpected parse result for src2")
	}
	fmt.Println(expr2.Marshal())

	// Parse an infix sum (left associative): {1 + 2 + 3}
	src3 := "{1 + 2 + 3}"
	expr3, rest3, err := Parse(Tokenize(src3))
	if err != nil || len(rest3) != 0 {
		panic("unexpected parse result for src3")
	}
	fmt.Println(expr3.Marshal())

	// Parse a right-associative comma list: {a, b, c}
	src4 := "{a, b, c}"
	expr4, rest4, err := Parse(Tokenize(src4))
	if err != nil || len(rest4) != 0 {
		panic("unexpected parse result for src4")
	}
	fmt.Println(expr4.Marshal())

	// Parse strings and nested functions
	src5 := "(print \"hello, world\" (upper (concat \"a\" \"b\")))"
	expr5, rest5, err := Parse(Tokenize(src5))
	if err != nil || len(rest5) != 0 {
		panic("unexpected parse result for src5")
	}
	fmt.Println(expr5.Marshal())

	// Parse a right-associative form with lambda: {x => y => (add x y)}
	src6 := "{x => y => (add x y)}"
	expr6, rest6, err := Parse(Tokenize(src6))
	if err != nil || len(rest6) != 0 {
		panic("unexpected parse result for src6")
	}
	fmt.Println(expr6.Marshal())
}
