package sorts

import (
	"strconv"

	"github.com/fbundle/sorts/slices_util"
)

func init() {
	setFinalNameParseFunc(func(name Name) Code {
		return Var{Name: name}
	})
}

type Var struct {
	Name Name
}

func (c Var) Form() Form {
	return c.Name
}

func (c Var) Eval(ctx Context) Sort {
	return ctx.Get(c.Name)
}

const (
	TypeCmd Name = "&"
)

func init() {
	addListParseFunc(TypeCmd, func(parse func(form Form) Code, list List) Code {
		err := compileErr(list, []string{string(TypeCmd), "value"})
		if len(list) != 1 {
			panic(err)
		}
		return Type{
			Value: parse(list[0]),
		}
	})
}

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

func init() {
	addListParseFunc(InhabitCmd, func(parse func(form Form) Code, list List) Code {
		err := compileErr(list, []string{string(InhabitCmd), "type"})
		if len(list) != 1 {
			panic(err)
		}
		return Inhabited{
			uuid: nextCount(),
			Type: parse(list[0]),
		}
	})
}

type Inhabited struct {
	uuid uint64
	Type Code
}

func (c Inhabited) Form() Form {
	return List{InhabitCmd, c.Type.Form(), Name(strconv.Itoa(int(c.uuid)))}
}

func (c Inhabited) Eval(ctx Context) Sort {
	t := c.Type.Eval(ctx)
	switch t := t.(type) {
	case Pi:
		return Pi{
			Param: t.Param,
			Body: Inhabited{
				uuid: nextCount(),
				Type: t.Body,
			},
		}
	default:
		return NewTerm(c.Form(), t)
	}
}

func init() {
	const (
		LambdaCmd Name = "=>"
	)
	addListParseFunc(LambdaCmd, func(parse func(form Form) Code, list List) Code {
		err := compileErr(list, []string{
			string(LambdaCmd),
			makeForm(AnnotCmd, "name1", "type1"),
			"...",
			makeForm(AnnotCmd, "nameN", "typeN"),
			"body",
		}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}
		params := slices_util.Map(list[:len(list)-1], func(form Form) Annot {
			return compileAnnot(parse, mustType[List](err, form)[1:])
		})
		output := parse(list[len(list)-1])
		slices_util.ForEach(slices_util.Reverse(params), func(param Annot) {
			output = Pi{
				Param: param,
				Body:  output,
			}
		})
		return output
	})
	const ArrowCmd Name = "->"
	addListParseFunc(ArrowCmd, func(parse func(form Form) Code, list List) Code {
		// make builtin like succ
		// e.g. if arrow is Nat -> Nat
		// then its lambda is
		// (x: Nat) => Nat
		// or some mechanism to introduce arrow type from pi type
		// TODO - probably we don't need this anymore
		panic("not implemented")
	})
}

func init() {
	setFinalListParseFunc(func(parse func(form Form) Code, list List) Code {
		err := compileErr(list, []string{
			"cmd",
			"arg1",
			"...",
			"argN",
		}, "where N >= 0")
		if len(list) < 1 {
			panic(err)
		}
		output := parse(list[0])
		args := slices_util.Map(list[1:], func(form Form) Code {
			return parse(form)
		})
		slices_util.ForEach(args, func(arg Code) {
			output = Beta{
				Cmd: output,
				Arg: arg,
			}
		})
		return output
	})
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
