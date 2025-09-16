package main

import (
	"encoding/json"
	"fmt"

	"github.com/fbundle/sorts/form2_processor"
)

var source string = `
# python-like syntax using transpiler
# start of block patterns: let, match, lambda, and (

let
   Nil = inh Any_2
   nil = inh Nil

   Bool = inh Any_2
   True = inh Bool
   False = inh Bool

   Nat = inh Any_2
   n0 = inh Nat
   succ = inh (Nat -> Nat)

   n1 = succ n0
   n2 = succ n1
   n3 = succ n2
   n4 = succ n3
   x = n1 ⊕ n2 ⊕ n3
   x = n1 ⊗ n2 ⊗ n3 ⊗ n4

   is_pos = lambda (x: Nat)
      match x with
         | succ z    => True
         | n0        => False

   must_pos = lambda (x: Nat)
      match x with
         | succ z    => x
         | n0        => nil


   print is_pos                     # resolved type as       Nat -> Bool
   print must_pos                   # resolved type as       Nat -> (Nat ⊕ Nil)
                                    # better to resolve as   Π_{x: Nat} B(x) where B(x) = (type (must_pos x))

   type Unit_0
`

func toString(o any) string {
	b, _ := json.Marshal(o)
	return string(b)
}

func main() {

	t := form2_processor.Tokenizer{
		LineCommentBegin: "#",
		SplitTokens:      []string{"+", "*", "$", "⊕", "⊗", "Π", "Σ", "=>", "->", ":", ",", "=", ":="},
	}
	p := form2_processor.Parser{
		OpenBlockTokens: []string{"let", "match", "lambda"},
		CloseBlockToken: "end",
		NewLineToken:    "__newline__",
	}

	lines := t.Tokenize(source)
	toks := p.Parse(lines)
	fmt.Println(toks)
}
