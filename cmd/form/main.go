package main

import (
	"encoding/json"
	"fmt"

	"github.com/fbundle/sorts/form_processor2"
)

func toString(o any) string {
	b, _ := json.Marshal(o)
	return string(b)
}

func main() {
	s := "   (hello=>x haha=y)"
	t := form_processor2.NewTokenizer([]string{
		"(", ")", "=", "=>",
	})
	indt, toks := t.Tokenize(s)
	fmt.Println(indt, toString(toks))
}
