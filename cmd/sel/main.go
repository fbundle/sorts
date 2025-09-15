package main

import (
	"fmt"
	"os"

	"github.com/fbundle/sorts/el"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
)

func mustReadSource(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustRun(tokens []form.Token) {
	ctx := el.Context{}.Reset()

	var node form.Form
	var err error
	for len(tokens) > 0 {
		tokens, node, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		fmt.Println(ctx.ToString(node))
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("REPL mode not implemented")
	}
	if len(args) != 1 {
		panic("usage: el2 <filename>")
	}

	tokens := form_processor.Tokenize(mustReadSource(args[0]))
	mustRun(tokens)
}
