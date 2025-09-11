package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/obsolete/el_v2"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	input := string(data)

	tokens := form.Tokenize(input)

	for len(tokens) > 0 {
		var formExpr form.Form
		formExpr, tokens, err = form.Parse(tokens)
		if err != nil {
			fmt.Printf("Error parsing form: %v\n", err)
			fmt.Println(tokens)
			os.Exit(1)
		}

		elExpr, err := el.el.ParseForm(formExpr)
		if err != nil {
			fmt.Printf("Error parsing el_v2 expression: %v\n", err)
			fmt.Printf("Next: %s\n", strings.Join(formExpr.Marshal(), " "))
			os.Exit(1)
		}

		// Use Marshal() to convert back to form and print
		marshaledForm := elExpr.Marshal()
		fmt.Printf("%T: %s\n", elExpr, strings.Join(marshaledForm.Marshal(), " "))
	}
}
