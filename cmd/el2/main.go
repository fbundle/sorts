package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fbundle/sorts/el2"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
)

func toString(o any) string {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustReadSource(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustRun(tokens []form.Token) {
	var ctx almost_sort_extra.Context = el2.Context{}.Reset()

	var node form.Form
	var err error
	var almostSort almost_sort_extra.AlmostSort
	for len(tokens) > 0 {
		tokens, node, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		ctx, almostSort = ctx.Compile(node)
		fmt.Println(toString(almostSort))
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
