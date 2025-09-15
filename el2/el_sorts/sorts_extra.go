package el_sorts

import (
	"fmt"

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
		err := fmt.Errorf("type annot must be (%s name type)", Head)
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

type LetBinding struct {
	Name  form.Name
	Value Sort
}
type MatchCase struct {
	Pattern any // Union[form.Form, Sort] - pattern matching vs exact matching
	Value   Sort
}

func (mc MatchCase) form(ctx Context) (form.Form, form.Form) {
	switch pattern := mc.Pattern.(type) {
	case form.Form: // pattern matching
		return pattern, ctx.Form(mc.Value)
	case Sort: // exact matching
		return ctx.Form(pattern), ctx.Form(mc.Value)
	default:
		panic(TypeErr)
	}
}
