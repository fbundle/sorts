package sorts

import "github.com/fbundle/sorts/slices_util"

const (
	LetCmd Name = "let"
)

type Let struct {
	Bindings []Binding
	Body     Code
}

func (c Let) Form() Form {
	var output List
	output = append(output, LetCmd)
	output = append(output, slices_util.Map(c.Bindings, func(binding Binding) Form {
		return binding.Form()
	})...)
	output = append(output, c.Body.Form())
	return output
}

func (c Let) Eval(ctx Context) Sort {
	subCtx := ctx
	slices_util.ForEach(c.Bindings, func(binding Binding) {
		name, valueCode := binding.Name, binding.Value
		value := valueCode.Eval(subCtx)
		subCtx = subCtx.Set(name, value)
	})
	return c.Body.Eval(subCtx)
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
