package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/slices_util"
	"github.com/fbundle/sorts/sorts"
	"github.com/fbundle/sorts/sorts_context"
	"github.com/fbundle/sorts/sorts_parser"
)

var ctx sorts.Context = sorts_context.Context{
	Frame: ordered_map.OrderedMap[form.Name, sorts.Sort]{},
	Univ: sorts_context.Univ{
		InitialTypeName:  "Unit",
		TerminalTypeName: "Any",
		DefaultTypeName:  "Type",
	},
}.Init()

type sortCode struct {
	sorts.Sort
}

func (s sortCode) Eval(ctx sorts.Context) sorts.Sort {
	return s.Sort
}

var Nat = sorts.NewTerm(form.Name("Nat"), sorts.NewChain("Parent", 2))

var Zero = sorts.NewTerm(form.Name("0"), Nat)

var Succ = sorts.Pi{
	Param: sorts.Annot{
		Name: "x",
		Type: sortCode{Nat},
	},
	Body: succBody{},
}

type succBody struct{}

func (l succBody) Form() form.Form {
	return form.List{form.Name("succ"), form.Name("x")}
}

func (l succBody) Eval(ctx sorts.Context) sorts.Sort {
	arg := ctx.Get("x")
	if !arg.Parent(ctx).LessEqual(ctx, Nat) {
		panic("x not of subtype of Nat")
	}

	x, err := strconv.Atoi(strings.Join(slices_util.Map(arg.Form().Marshal(), func(tok form.Token) string {
		return tok
	}), " "))
	if err != nil {
		panic(err)
	}
	y := x + 1
	return sorts.NewTerm(form.Name(strconv.Itoa(y)), Nat)
}

func parseForm(s string) <-chan sorts.Code {
	ch := make(chan sorts.Code)
	go func() {
		defer close(ch)
		parse := sorts_parser.Parser{}.Init().Parse
		tokens := form_processor.Tokenize(s)
		var f form.Form
		var err error
		for len(tokens) > 0 {
			tokens, f, err = form_processor.Parse(tokens)
			if err != nil {
				panic(err)
			}
			ch <- parse(f)
		}
	}()
	return ch
}

var source = `
(succ 0)
(succ (succ 0))

(let
	{x := 0}
	x
)

(let
	{Nat := (* Any_2)}
	{0 := (* Nat)}
	{succ := (* {{_: Nat} => Nat})} # or syntactic sugar {succ := (* {Nat -> Nat})}

	(succ (succ 0))
)

`

func main() {
	ctx = ctx.
		Set("Nat", Nat).
		Set("0", Zero).
		Set("succ", Succ)

	for code := range parseForm(source) {
		sort := code.Eval(ctx) // evaluate
		fmt.Println(sort.Form(), "<-", code.Form())
	}
}
