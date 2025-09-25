package main

import (
	"fmt"
	"strconv"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
	"github.com/fbundle/sorts/slices_util"
	"github.com/fbundle/sorts/sorts"
	"github.com/fbundle/sorts/sorts_context"
	"github.com/fbundle/sorts/sorts_parser"
)

var ctx sorts.Context = sorts_context.Context{
	Univ: sorts_context.Univ{
		InitialTypeName:  "Unit",
		TerminalTypeName: "Any",
		DefaultTypeName:  "Type",
	},
}

var parser = sorts_parser.Parser{}.Init()

func parseForm(s string) <-chan sorts.Code {
	ch := make(chan sorts.Code)
	go func() {
		defer close(ch)

		tokens := form_processor.Tokenize(s)
		var f form.Form
		var err error
		for len(tokens) > 0 {
			tokens, f, err = form_processor.Parse(tokens)
			if err != nil {
				panic(err)
			}
			ch <- parser.Parse(f)
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
	{1 := (succ 0)}
	{2 := (succ 1)}
	(add 1 2)
)

(let
	{Nat := (* Any_2)}
	{0 := (* Nat)}
	{succ := (* {{_: Nat} => Nat})} # or syntactic sugar {succ := (* {Nat -> Nat})}

	(succ (succ 0))
)

`

func main() {
	Any2 := parser.Parse(form.Name("Any_2")).Eval(ctx)

	Nat := sorts.MakeBuiltinSort(
		"Nat",
		Any2,
		nil,
		func(args []form.Form) form.Form {
			if len(args) != 0 {
				panic("internal")
			}
			return form.Name("Nat")
		},
	)
	Zero := sorts.MakeBuiltinSort(
		"0",
		Nat,
		nil,
		func(args []form.Form) form.Form {
			if len(args) != 0 {
				panic("internal")
			}
			return form.Name("0")
		},
	)
	Succ := sorts.MakeBuiltinSort(
		"succ",
		Nat,
		[]sorts.Sort{Nat},
		func(args []form.Form) form.Form {
			if len(args) != 1 {
				panic("internal")
			}
			x, err := strconv.Atoi(form.String(args[0]))
			if err != nil {
				panic(err)
			}
			y := x + 1
			ret := form.Name(strconv.Itoa(y))
			return ret
		},
	)

	Add := sorts.MakeBuiltinSort(
		"add",
		Nat,
		[]sorts.Sort{Nat, Nat},
		func(args []form.Form) form.Form {
			if len(args) != 2 {
				panic("internal")
			}
			values := slices_util.Map(args, func(f form.Form) int {
				v, err := strconv.Atoi(form.String(f))
				if err != nil {
					panic(err)
				}
				return v
			})
			output := values[0] + values[1]
			ret := form.Name(strconv.Itoa(output))
			return ret
		},
	)

	ctx = ctx.
		Set("Nat", Nat).
		Set("0", Zero).
		Set("succ", Succ).
		Set("add", Add)

	for code := range parseForm(source) {
		fmt.Println("evaluating", code.Form())
		sort := code.Eval(ctx) // evaluate
		fmt.Println("\t value:", sort.Form())
	}
}
