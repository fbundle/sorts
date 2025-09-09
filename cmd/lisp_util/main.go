package main

import (
	"fmt"

	"github.com/fbundle/sorts/lisp_util"
)

func testLexer() {
	s := "(hello \"this is\" a house ) hehe haha 1231 ( this( is \"another \\\" house\"))"
	tokens := lisp_util.Tokenize(s)
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

func main() {
	testLexer()
}
