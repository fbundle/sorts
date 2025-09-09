package main

import (
	"fmt"

	"github.com/fbundle/sorts/ast"
)

func testLexer() {
	s := "(hell=>o \"th=>is is\" a house ) hehe haha 1231 ( this( is \"another \\\" house\"))"
	tokens := ast.Tokenize(s)
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

func main() {
	testLexer()
}
