package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fbundle/sorts/el"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
)

func mustReadSource() string {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustRun(tokens []form.Token) {
	ctx := el.Context{}.Init()

	var node form.Form
	var err error
	for len(tokens) > 0 {
		tokens, node, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		sort := ctx.Compile(node)
		sort = sort.Reduce(ctx)
		fmt.Println(sort.Form())
	}
}

func main() {
	tokens := form_processor.Tokenize(mustReadSource())
	mustRun(tokens)
}
