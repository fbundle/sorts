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

type mode string

const (
	modeComp mode = "comp"
	modeEval mode = "eval"
)

var currentMode = modeComp

func toNat(f form.Form) int {
	switch f := f.(type) {
	case form.Name:
		if f != "0" {
			panic("internal")
		}
		return 0
	case form.List:
		if len(f) != 2 {
			panic("internal")
		}
		if f[0] != form.Name("succ") {
			panic("internal")
		}
		return 1 + toNat(f[1])
	default:
		panic("internal")
	}
}

func fromNat(n int) form.Form {
	var output form.Form = form.Name("0")
	for i := 0; i < n; i++ {
		output = form.List{form.Name("succ"), output}
	}
	return output
}

// el - basic EL with naturals and addition
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
			switch currentMode {
			case modeComp:
				return form.List{form.Name("succ"), args[0]}
			case modeEval:
				x, err := strconv.Atoi(form.String(args[0]))
				if err != nil {
					panic(err)
				}
				y := x + 1
				ret := form.Name(strconv.Itoa(y))
				return ret
			}
			panic("internal")
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
			switch currentMode {
			case modeComp:
				values := slices_util.Map(args, toNat)
				output := values[0] + values[1]
				return fromNat(output)
			case modeEval:
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
			}
			panic("internal")
		},
	)).(sorts_context.Context)

	ctx = ctx.
		WithBuiltin(func(name form.Name) (sorts.Sort, bool) {
			// add integer
			v, err := strconv.Atoi(string(name))
			if err != nil {
				return nil, false
			}
			switch currentMode {
			case modeComp:
				return sorts.MakeBuiltinSort(
					fromNat(v),
					Nat,
					nil, nil,
				), true
			case modeEval:
				return sorts.MakeBuiltinSort(
					name,
					Nat,
					nil, nil,
				), true
			}
			panic("internal")
		})
	return ctx, sorts_parser.Parser{}.Init().Parse
}

var source = `
(succ 0)
(succ (succ 0))

(let
	{x := 0}
	x
)

(add 3 5)

(let
	{Nat := (* Any_2)}
	{0 := (* Nat)}
	{succ := (* {{_: Nat} => Nat})} # or syntactic sugar {succ := (* {Nat -> Nat})}

	(succ (succ 0))
)

`

func main() {
	ctx, parse := el()

	tokens := form_processor.Tokenize(source)
	var f form.Form
	var err error
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
