package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
	"github.com/fbundle/sorts/slices_util"
	"github.com/fbundle/sorts/sorts"
	"github.com/fbundle/sorts/sorts_context"
	"github.com/fbundle/sorts/sorts_parser"
)

// el - basic EL with integers, succ, and addition
func el() (sorts.Context, func(form.Form) sorts.Code) {
	ctx := sorts_context.Context{
		Univ: sorts_context.Univ{
			InitialTypeName:  "Unit",
			TerminalTypeName: "Any",
			DefaultTypeName:  "Type",
		},
	}
	Type2 := ctx.Default(2)
	Nat := sorts.MakeBuiltinSort(
		form.Name("Nat"),
		Type2,
		nil, nil,
	)
	ctx = ctx.Set("Nat", Nat).(sorts_context.Context)
	ctx = ctx.Set("succ", sorts.MakeBuiltinSort(
		form.Name("succ"),
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
	)).(sorts_context.Context)
	ctx = ctx.Set("add", sorts.MakeBuiltinSort(
		form.Name("add"),
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
	)).(sorts_context.Context)

	ctx = ctx.
		WithBuiltin(func(name form.Name) (sorts.Sort, bool) {
			// add integer
			_, err := strconv.Atoi(string(name))
			if err != nil {
				return nil, false
			}
			return sorts.MakeBuiltinSort(
				name,
				Nat,
				nil, nil,
			), true
		})
	return ctx, sorts_parser.Parser{}.Init().Parse
}

func main() {
	ctx, parse := el()

	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	source := string(b)

	tokens := form_processor.Tokenize(source)
	var f form.Form
	for len(tokens) > 0 {
		tokens, f, err = form_processor.Parse(tokens)
		if err != nil {
			panic(err)
		}
		code := parse(f)
		fmt.Println("evaluating", code.Form())
		sort := code.Eval(ctx) // evaluate
		fmt.Println("\t value:", sort.Form())
	}
}
