package el_sorts

import "github.com/fbundle/sorts/form"

type TypeAnnot struct {
	Name form.Name
	Type Sort
}

func ListParseTypeAnnot(ctx Context, args form.List) TypeAnnot {

}

type LambdaParam struct {
	Name form.Name
	Type Sort
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
