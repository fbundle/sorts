package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fbundle/sorts/expr"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	input := string(data)

	tokens := expr.Tokenize(input)

	var e expr.Form
	for len(tokens) > 0 {
		e, tokens, err = expr.Parse(tokens)
		if err != nil {
			fmt.Printf("Error parsing input: %v\n", err)
			fmt.Println(tokens)
			os.Exit(1)
		}

		fmt.Println(strings.Join(e.Marshal(), " "))
	}

}
