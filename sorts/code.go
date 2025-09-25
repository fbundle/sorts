package sorts

import "github.com/fbundle/sorts/slices_util"

type Symbol struct {
	Name Name
}

func (c Symbol) Form() Form {
	return c.Name
}

func (c Symbol) Eval(ctx Context) Sort {
	return ctx.Get(c.Name)
}

const (
	TypeCmd Name = "&"
)

type Type struct {
	Value Code
}

func (c Type) Form() Form {
	return List{TypeCmd, c.Value.Form()}
}

func (c Type) Eval(ctx Context) Sort {
	return c.Value.Eval(ctx).Parent(ctx)
}

const (
	InhabitCmd Name = "*"
)

type Inhabited struct {
	Type Code
}

func (c Inhabited) Form() Form {
	return List{InhabitCmd, c.Type.Form()}
}

func (c Inhabited) Eval(ctx Context) Sort {
	t := c.Type.Eval(ctx)
	switch t := t.(type) {
	case Pi:
		return Pi{
			Params: t.Params,
			Body: Inhabited{
				Type: t.Body,
			},
		}
	case Atom:
		return NewTerm(c.Form(), t)
	}
	panic("unreachable")
}

type Beta struct {
	Cmd  Code
	Args []Code
}

func (c Beta) Form() Form {
	var output List
	output = append(output, c.Cmd.Form())
	output = append(output, slices_util.Map(c.Args, func(code Code) Form {
		return code.Form()
	})...)
	return output
}

func (c Beta) Eval(ctx Context) Sort {
	cmd := mustType[Pi](TypeError, c.Cmd.Eval(ctx))

	if len(cmd.Params) != len(c.Args) {
		panic(TypeError)
	}
	subCtx := ctx
	for i := 0; i < len(cmd.Params); i++ {
		param := cmd.Params[i]
		paramType := param.Type.Eval(ctx)
		arg := c.Args[i].Eval(ctx)
		argType := arg.Parent(ctx)
		if !argType.LessEqual(ctx, paramType) {
			panic(TypeError)
		}
		subCtx = subCtx.Set(param.Name, arg)
	}

	return cmd.Body.Eval(subCtx)
}
