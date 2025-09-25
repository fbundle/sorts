package sorts

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
			Param: t.Param,
			Body: Inhabited{
				Type: t.Body,
			},
		}
	default:
		return NewTerm(c.Form(), t)
	}
}

type Beta struct {
	Cmd Code
	Arg Code
}

func (c Beta) Form() Form {
	return List{c.Cmd.Form(), c.Arg.Form()}
}

func (c Beta) Eval(ctx Context) Sort {
	cmd := mustType[Pi](TypeError, c.Cmd.Eval(ctx))
	paramType := cmd.Param.Type.Eval(ctx)
	arg := c.Arg.Eval(ctx)
	argType := arg.Parent(ctx)
	if !argType.LessEqual(ctx, paramType) {
		panic(TypeError)
	}
	return cmd.Body.Eval(ctx.Set(cmd.Param.Name, arg))
}
