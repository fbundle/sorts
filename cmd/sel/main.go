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

	var letBindings [][3]form.Form

	var node1, node2, node3 form.Form
	var err error
	for len(tokens) > 0 {
		tokens, node1, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		name, ok := node1.(form.Name)
		if !ok {
			panic("node1 must be name")
		}
		if name == "@inspect" {
			tokens, node2, err = form_processor.Parse(tokens)
			if err != nil {
				panic(err)
			}
			letBindings = append(letBindings, [3]form.Form{
				form.Name("_"), form.Name(":="), form.List{
					form.Name("inspect"), node2,
				},
			})
			continue
		}
		tokens, node2, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		tokens, node3, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		letBindings = append(letBindings, [3]form.Form{
			node1, node2, node3,
		})
	}

	let := form.List{form.Name("let")}
	for _, binding := range letBindings {
		let = append(let, form.List(binding[:]))
	}
	let = append(let, form.Name("Unit_0"))

	sort := ctx.Compile(let)
	fmt.Println(ctx.ToString(sort))
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
