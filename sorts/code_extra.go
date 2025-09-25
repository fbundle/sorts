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
	value := c.Binding.Value.Eval(ctx)
	ctx = ctx.Set(c.Binding.Name, value)
	return c.Body.Eval(ctx)
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
