package sorts

const (
	LetCmd Name = "let"
)

type Let struct {
	Binding Binding
	Body    Code
}

func (c Let) Form() Form {
	return List{LetCmd, c.Binding.Form(), c.Body.Form()}
}

func (c Let) Eval(ctx Context) Sort {
	name, valueCode := c.Binding.Name, c.Binding.Value
	var value Sort
	if inh, ok := c.Binding.Value.(Inhabited); ok {
		// if name binding an inhabited, then rename it
		value = NewTerm(name, inh.Type.Eval(ctx))
	} else {
		value = valueCode.Eval(ctx)
	}
	value = valueCode.Eval(ctx)
	return c.Body.Eval(ctx.Set(name, value))
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
