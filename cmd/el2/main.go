package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fbundle/sorts/el2"
	"github.com/fbundle/sorts/el2/parser"
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

func mustRun(r el2.Runtime, p el2_parser.Parser, tokens []form.Token) {
	var node form.Form
	var err error
	for len(tokens) > 0 {
		tokens, node, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		s := p.Parse(node)
		fmt.Println(toString(s))
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

	r := el2.NewRuntime()
	p := el2.NewParser(r)
	tokens := form_processor.Tokenize(mustReadSource(args[0]))
	mustRun(r, p, tokens)
}
