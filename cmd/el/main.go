package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fbundle/sorts/el"
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	input := string(data)

	tokens := form.Tokenize(input)

	frame := el.Frame{}

	for len(tokens) > 0 {
		var formExpr form.Form
		formExpr, tokens, err = form.Parse(tokens)
		if err != nil {
			fmt.Printf("Error parsing form: %v\n", err)
			fmt.Println(tokens)
			os.Exit(1)
		}

		elExpr, err := el.ParseForm(formExpr)
		if err != nil {
			fmt.Printf("Error parsing  expression: %v\n", err)
			fmt.Printf("Next: %s\n", strings.Join(formExpr.Marshal(), " "))
			os.Exit(1)
		}

		fmt.Println("expr", el.String(elExpr))

		var sort sorts.Sort
		var value el.Expr
		frame, sort, value, err = elExpr.Resolve(frame)
		if err != nil {
			fmt.Printf("Error evaluating expression: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("sort\t", sorts.Name(sort))
		fmt.Println("value\t", el.String(value))
	}
}
