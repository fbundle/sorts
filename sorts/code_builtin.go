package sorts

/*
type Param struct {
	Name Name
	Type Sort
}

func Builtin(ctx Context, name Name, params []Param, retType Sort, body func(args []Form) Form) Context {
	if len(params) == 0 {
		return ctx.Set(name, NewTerm(body(nil), retType))
	}
	var output Sort = newBuiltinPi(
		"",
		params[len(params)-1],
		retType,
		func(arg Form) Form {
			args := slices_util.Map(params[:len(params)-1], func(param Param) Form {

			})
		},
	)
	for i := len(params) - 2; i > 0; i-- {
		param := params[i]

	}

}

type sortCode struct {
	Sort
}

func (s sortCode) Eval(ctx Context) Sort {
	return s.Sort
}

type builtinPiBody struct {
	param   Param
	body    func(arg Form) Form
	retType Sort
}

func (c builtinPiBody) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (c builtinPiBody) Eval(ctx Context) Sort {
	value := ctx.Get(c.param.Name)
	if !value.Parent(ctx).LessEqual(ctx, c.param.Type) {
		panic(TypeError)
	}
	arg := value.Form()
	ret := c.body(arg)
	return NewTerm(ret, c.retType)
}

type piWithName struct {
	name Name
	Pi
}

func (c piWithName) Form() Form {
	if len(c.name) > 0 {
		return c.name
	}
	return c.Pi.Form()
}

func newBuiltinPi(name Name, param Param, retType Sort, body func(arg Form) Form) piWithName {
	return piWithName{
		name: name,
		Pi: Pi{
			Param: Annot{
				Name: param.Name,
				Type: sortCode{param.Type},
			},
			Body: builtinPiBody{
				param:   param,
				body:    body,
				retType: retType,
			},
		},
	}
}
*/
