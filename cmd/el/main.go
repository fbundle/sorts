package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/form_processor"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type context struct {
	dict ordered_map.OrderedMap[form.Name, sorts.Sort]
}

func (c context) Parse(f form.Form) sorts.Sort {

	switch f := f.(type) {
	case form.Name:
		if s, ok := c.dict.Get(f); ok {
			return s
		}
		panic(fmt.Errorf("name_not_found: %s", f))
	case form.List:
		if name, ok := f[0].(form.Name); ok {
			if parseFunc, ok := sorts.ListParseFuncMap[name]; ok {
				return parseFunc(c, f[1:])
			}
		}
		return sorts.DefaultParseFunc(c, f)
	}
	panic(fmt.Errorf("parse_error: %v", f))
}

func (c context) Mode() sorts.Mode {
	return sorts.ModeDebug
}

// LessEqual implements sorts.Context.
func (c context) LessEqual(s form.Form, d form.Form) bool {
	return s == d
}

// Set implements sorts.Context.
func (c context) Set(name sorts.Name, sort sorts.Sort) sorts.Context {
	return context{
		dict: c.dict.Set(name, sort),
	}
}

var _ sorts.Context = context{}

var Nat = sorts.NewTerm(form.Name("Nat"), sorts.NewChain("Parent", 2))

var Zero = sorts.NewTerm(form.Name("0"), Nat)

var Succ = sorts.Lambda{
	Param: sorts.Annot{
		Name: "x",
		Type: Nat,
	},
	Body: succBody{},
}

type succBody struct{}

// Form implements sorts.Sort.
func (l succBody) Form() form.Form {
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
	c := ctx.(context)
	arg, ok := c.dict.Get("x")
	if !ok {
		panic("x not set")
	}
	if !arg.Parent(ctx).LessEqual(ctx, Nat) {
		panic("x not of subtype of Nat")
	}
	return Nat
}

// Reduce implements sorts.Sort.
func (l succBody) Eval(ctx sorts.Context) sorts.Sort {
	c := ctx.(context)
	arg, ok := c.dict.Get("x")
	if !ok {
		panic("x not set")
	}
	if !arg.Parent(ctx).LessEqual(ctx, Nat) {
		panic("x not of subtype of Nat")
	}
	arg = arg.Eval(ctx)

	x, err := strconv.Atoi(strings.Join(slicesMap(arg.Form().Marshal(), func(tok form.Token) string {
		return tok
	}), " "))
	if err != nil {
		panic(err)
	}
	y := x + 1
	return sorts.NewTerm(form.Name(strconv.Itoa(y)), Nat)
}

func parseForm(s string) <-chan form.Form {
	ch := make(chan form.Form)
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
			ch <- f
		}
	}()
	return ch
}

var source = `
(succ 0)
(succ (succ 0))
`

func main() {
	c := context{}.
		Set("Nat", Nat).
		Set("0", Zero).
		Set("succ", Succ)

	for f := range parseForm(source) {
		s := c.Parse(f)
		_ = s.Parent(c) // type check
		s1 := s.Eval(c) // evaluate
		fmt.Println(s1.Form(), "<-", s.Form())
	}
}
func slicesMap[T1 any, T2 any](input []T1, f func(T1) T2) []T2 {
	output := make([]T2, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
