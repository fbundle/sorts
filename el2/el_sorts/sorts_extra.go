package el_sorts

import (
	"fmt"
	"log"

	"github.com/fbundle/sorts/form"
)

func mustName(err error, node form.Form) form.Name {
	name, ok := node.(form.Name)
	if !ok {
		panic(err)
	}
	return name
}
func mustList(err error, node form.Form) form.List {
	list, ok := node.(form.List)
	if !ok {
		panic(err)
	}
	return list
}

func mustMatchHead(err error, Head form.Name, list form.List) {
	if len(list) == 0 || Head != list[0] {
		panic(err)
	}
}

type TypeAnnot struct {
	Name form.Name
	Type Sort
}

func ListParseTypeAnnot(Head form.Name) func(ctx Context, list form.List) TypeAnnot {
	return func(ctx Context, list form.List) TypeAnnot {
		err := fmt.Errorf("type_annot must be (%s name type)", Head)
		mustMatchHead(err, Head, list)
		if len(list) != 3 {
			panic(err)
		}
		nameForm, typeFrom := list[1], list[2]

		return TypeAnnot{
			Name: mustName(err, nameForm),
			Type: ctx.Compile(typeFrom),
		}
	}
}

type NameBinding struct {
	Name  form.Name
	Value Sort
}

func ParseNameBinding(Head form.Name) func(ctx Context, list form.List) NameBinding {
	return func(ctx Context, list form.List) NameBinding {
		err := fmt.Errorf("name_binding must be (%s name value)", Head)
		mustMatchHead(err, Head, list)
		if len(list) != 3 {
			panic(err)
		}

		nameForm, valueForm := list[1], list[2]
		binding := NameBinding{
			Name:  mustName(err, nameForm),
			Value: ctx.Compile(valueForm),
		}

		if v, ok := binding.Value.(Inhabitant); ok {
			// rename inhabitant
			binding.Value = Inhabitant{
				Atom: ctx.NewTerm(binding.Name, ctx.Parent(v.Atom)),
				Head: v.Head,
				Name: binding.Name,
			}
			log.Printf("rename inhabitant %s -> %s\n", v.Name, binding.Name)
		}
		return binding
	}
}

type MatchCase struct {
	CondType Sort
	Constr   form.Name
	Vars     []form.Name
	Value    Sort
}

func ParseMatchCase(Head form.Name, DefaultConstr form.Name, CondType Sort) func(ctx Context, list form.List) MatchCase {
	return func(ctx Context, list form.List) MatchCase {
		err := fmt.Errorf("match_lambda must be (%s (constr variable...) value)", Head)
		mustMatchHead(err, Head, list)
		if len(list) != 3 {
			panic(err)
		}

		// TODO -  allow more sophisicated pattern, like (succ (succ x))

		var constr form.Name
		var vars []form.Name
		switch pattern := list[1].(type) {
		case form.Name:
			constr = pattern
		case form.List:
			constr = mustName(err, pattern[0])
			for i := 1; i < len(pattern); i++ {
				vars = append(vars, mustName(err, pattern[i]))
			}
		default:
			panic("unreachable")
		}

		// TODO - check if CondType is inductive and can be destructed into pattern
		// maybe add a destruction function so that it can work like a lambda
		return MatchCase{
			CondType: CondType,
			Constr:   constr,
			Vars:     vars,
			Value:    ctx.Compile(list[2]),
		}
	}
}
