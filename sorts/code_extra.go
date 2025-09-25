package sorts

import (
	"github.com/fbundle/sorts/slices_util"
)

const (
	LetCmd Name = "let"
)

func init() {
	addListParseFunc(LetCmd, func(parse func(form Form) Code, list List) Code {
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
		bindings := slices_util.Map(list[:len(list)-1], func(form Form) Binding {
			return compileBinding(parse, mustType[List](err, form)[1:])
		})
		output := parse(list[len(list)-1])
		slices_util.ForEach(slices_util.Reverse(bindings), func(binding Binding) {
			output = Let{
				Binding: binding,
				Body:    output,
			}
		})
		return output
	})
}

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
