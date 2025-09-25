package sorts

import (
	"strconv"
	"sync/atomic"

	"github.com/fbundle/sorts/slices_util"
)

func Builtin(ctx Context, name Name, paramTypes []Sort, retType Sort, body func(args []Form) Form) Context {
	if len(paramTypes) == 0 {
		return ctx.Set(name, NewTerm(body(nil), retType))
	}
	var count uint64
	params := slices_util.Map(paramTypes, func(paramType Sort) Annot {
		i := atomic.AddUint64(&count, 1)
		paramName := Name("x_" + strconv.Itoa(int(i)))
		return Annot{
			Name: paramName,
			Type: sortCode{paramType},
		}
	})
	sort := piWithName{
		name: name,
		Pi: Pi{
			Params: params,
			Body: builtinPiBody{
				params:  params,
				retType: retType,
				body:    body,
			},
		},
	}
	return ctx.Set(name, sort)
}

type sortCode struct {
	Sort
}

func (s sortCode) Eval(ctx Context) Sort {
	return s.Sort
}

type builtinPiBody struct {
	params  []Annot
	body    func(arg []Form) Form
	retType Sort
}

func (c builtinPiBody) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (c builtinPiBody) Eval(ctx Context) Sort {
	args := slices_util.Map(c.params, func(param Annot) Form {
		value := ctx.Get(param.Name)
		if !value.Parent(ctx).LessEqual(ctx, param.Type.Eval(ctx)) {
			panic(TypeError)
		}
		return value.Form()
	})
	ret := c.body(args)
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
