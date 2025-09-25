package sorts

const (
	LetCmd Name = "let"
)

func init() {
	ListParseFuncMap[LetCmd] = func(ctx Context, list List) Code {
		err := compileErr(list, []string{
			string(LetCmd),
			makeForm(BindingCmd, "name1", "value1"),
			"...",
			makeForm(BindingCmd, "nameN", "valueN"),
			"body",
		}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}
		
	}
}

type Let struct {
	Binding Binding
	Body    Code
}

func (c Let) Form() Form {
	return List{LetCmd, c.Binding.Form(), c.Body.Form()}
}

func (c Let) Eval(ctx Context) Sort {
	return c.Body.Eval(ctx.Set(c.Binding.Name, c.Binding.Value.Eval(ctx)))
}

const (
	InductiveCmd Name = "inductive"
)

const (
	MatchCmd Name = "match"
)
