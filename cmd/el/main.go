package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type context struct {
	dict ordered_map.OrderedMap[form.Name, sorts.Sort]
}

// LessEqual implements sorts.Context.
func (c context) LessEqual(s sorts.Form, d sorts.Form) bool {
	return s == d
}

// Set implements sorts.Context.
func (c context) Set(name sorts.Name, sort sorts.Sort) sorts.Context {
	return context{
		dict: c.dict.Set(name, sort),
	}
}

var _ sorts.Context = context{}

var Nat = sorts.NewTerm(form.Name("Nat"), sorts.NewChain("Type", 2))

var Zero = sorts.NewTerm(form.Name("0"), Nat)

var Succ = sorts.Pi{
	Param: sorts.Annot{
		Name: "x",
		Type: Nat,
	},
	Body: succBody{},
}

type succBody struct{}

// Compile implements sorts.Sort.
func (l succBody) Compile(ctx sorts.Context) sorts.Sort {
	c := ctx.(context)
	arg, ok := c.dict.Get("x")
	if !ok {
		panic("x not set")
	}
	if arg.Parent(ctx).LessEqual(ctx, Nat) {
		panic("x not of subtype of Nat")
	}
	return sorts.NewTerm(form.List{form.Name("Succ"), arg.Form()}, Nat)
}

// Form implements sorts.Sort.
func (l succBody) Form() sorts.Form {
	return form.List{form.Name("Succ"), form.Name("x")}
}

// LessEqual implements sorts.Sort.
func (l succBody) LessEqual(ctx sorts.Context, d sorts.Sort) bool {
	panic("unimplemented")
}

// Level implements sorts.Sort.
func (l succBody) Level(ctx sorts.Context) int {
	panic("unimplemented")
}

// Parent implements sorts.Sort.
func (l succBody) Parent(ctx sorts.Context) sorts.Sort {
	panic("unimplemented")
}

// Reduce implements sorts.Sort.
func (l succBody) Reduce(ctx sorts.Context) sorts.Sort {
	c := ctx.(context)
	arg, ok := c.dict.Get("x")
	if !ok {
		panic("x not set")
	}
	if arg.Parent(ctx).LessEqual(ctx, Nat) {
		panic("x not of subtype of Nat")
	}

	x, err := strconv.Atoi(strings.Join(slicesMap(arg.Form().Marshal(), func(tok form.Token) string {
		return string(tok)
	}), " "))
	if err != nil {
		panic(err)
	}
	y := x + 1
	return sorts.NewTerm(form.Name(strconv.Itoa(y)), Nat)
}

func main() {
	c := context{}

	fmt.Println(Nat.Form())
	fmt.Println(Zero.Form())
	fmt.Println(Succ.Form())
	One := sorts.Beta{
		Cmd: Succ,
		Arg: Zero,
	}
	Two := sorts.Beta{
		Cmd: Succ,
		Arg: One,
	}
	fmt.Println(One.Compile(c).Form())
	fmt.Println(Two.Compile(c).Form())

	fmt.Println(One.Reduce(c).Form())
	fmt.Println(Two.Reduce(c).Form())
}
func slicesMap[T1 any, T2 any](input []T1, f func(T1) T2) []T2 {
	output := make([]T2, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
