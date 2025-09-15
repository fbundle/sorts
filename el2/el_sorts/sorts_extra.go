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

type MatchLambda struct {
	Pattern any // Union[form.Form, Sort] - pattern matching vs exact matching
	Value   Sort
}

func ParseMatchLambda(Head form.Name, CondType Sort) func(ctx Context, list form.List) MatchLambda {
	return func(ctx Context, list form.List) MatchLambda {
		err := fmt.Errorf("match_lambda must be (%s (constr var) value)", Head)
		mustMatchHead(err, Head, list)
		if len(list) != 3 {
			panic(err)
		}

	}
}

func (mc MatchLambda) form(ctx Context) (form.Form, form.Form) {
	switch pattern := mc.Pattern.(type) {
	case form.Form: // pattern matching
		return pattern, ctx.Form(mc.Value)
	case Sort: // exact matching
		return ctx.Form(pattern), ctx.Form(mc.Value)
	default:
		panic(TypeErr)
	}
}
