package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	input := string(data)

	tokens := form_processor.Tokenize(input)

	var e form.Form
	for len(tokens) > 0 {
		e, tokens, err = form_processor.Parse(tokens)
		if err != nil {
			fmt.Printf("Error parsing input: %v\n", err)
			fmt.Println(tokens)
			os.Exit(1)
		}

		fmt.Println(strings.Join(e.Marshal("(", ")"), " "))
	}

}
