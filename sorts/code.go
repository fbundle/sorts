package sorts

import (
	"strconv"

	"github.com/fbundle/sorts/slices_util"
)

const (
	TypeCmd Name = "&"
)

func init() {
	ListParseFuncMap[TypeCmd] = func(ctx Context, list List) Code {
		err := compileErr(list, []string{string(TypeCmd), "value"})
		if len(list) != 1 {
			panic(err)
		}
		return Type{
			Value: ctx.Parse(list[0]),
		}
	}
}

type Type struct {
	Value Code
}

func (s Type) Form() Form {
	return List{TypeCmd, s.Value.Form()}
}

func (s Type) Eval(ctx Context) Sort {
	return s.Value.Eval(ctx).Parent(ctx)
}

const (
	InhabitCmd Name = "*"
)

func init() {
	ListParseFuncMap[InhabitCmd] = func(ctx Context, list List) Code {
		err := compileErr(list, []string{string(InhabitCmd), "type"})
		if len(list) != 1 {
			panic(err)
		}
		return Inhabited{
			uuid: nextCount(),
			Type: ctx.Parse(list[0]),
		}
	}
}

type Inhabited struct {
	uuid uint64
	Type Code
}

func (s Inhabited) Form() Form {
	return List{InhabitCmd, s.Type.Form(), Name(strconv.Itoa(int(s.uuid)))}
}

func (s Inhabited) Eval(ctx Context) Sort {
	t := s.Type.Eval(ctx)
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
		return NewTerm(s.Form(), t)
	}
}

func init() {
	const (
		LambdaCmd Name = "=>"
	)
	ListParseFuncMap[LambdaCmd] = func(ctx Context, list List) Code {
		err := compileErr(list, []string{string(LambdaCmd), "param1", "...", "paramN", "body"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		params := slices_util.Map(list[:len(list)-1], func(form Form) Annot {
			return compileAnnot(ctx, mustType[List](err, form))
		})
		output := ctx.Parse(list[len(list)-1])
		slices_util.ForEach(slices_util.Reverse(params), func(param Annot) {
			output = Pi{
				Param: param,
				Body:  output,
			}
		})
		return output
	}
	const ArrowCmd Name = "->"
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) Code {
		// make builtin like succ
		// e.g. if arrow is Nat -> Nat
		// then its lambda is
		// (x: Nat) => Nat
		// or some mechanism to introduce arrow type from pi type
		// TODO - probably we don't need this anymore
		panic("not implemented")
	}
}

func init() {
	DefaultParseFunc = func(ctx Context, list List) Code {
		err := compileErr(list, []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		output := ctx.Parse(list[0])
		args := slices_util.Map(list[1:], func(form Form) Code {
			return ctx.Parse(form)
		})
		slices_util.ForEach(args, func(arg Code) {
			output = Beta{
				Cmd: output,
				Arg: arg,
			}
		})
		return output
	}
}

type Beta struct {
	Cmd Code
	Arg Code
}

func (s Beta) Form() Form {
	return List{s.Cmd.Form(), s.Arg.Form()}
}

func (s Beta) Eval(ctx Context) Sort {
	cmd := mustType[Pi](TypeError, s.Cmd.Eval(ctx))
	paramType := cmd.Param.Type.Eval(ctx)
	arg := s.Arg.Eval(ctx)
	argType := arg.Parent(ctx)
	if !argType.LessEqual(ctx, paramType) {
		panic(TypeError)
	}
	return cmd.Body.Eval(ctx.Set(cmd.Param.Name, arg))
}
