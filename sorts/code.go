package sorts

import (
	"strconv"

	"github.com/fbundle/sorts/slices_util"
)

func init() {
	ListParseFuncMap[TypeCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(TypeCmd), "value"})
		if len(list) != 1 {
			panic(err)
		}
		return Type{
			Value: ctx.Parse(list[0]),
		}
	}
}

const (
	TypeCmd Name = "&"
)

type Type struct {
	Value Code
}

func (s Type) Form() Form {
	return List{TypeCmd, s.Value.Form()}
}

func (s Type) Eval(ctx Context) Sort {
	return s.Value.Parent(ctx)
}

func init() {
	ListParseFuncMap[InhabitCmd] = func(ctx Context, list List) Sort {
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

const (
	InhabitCmd Name = "*"
)

type Inhabited struct {
	code
	uuid uint64
	Type Sort
}

func (s Inhabited) Form() Form {
	return List{InhabitCmd, s.Type.Form(), Name(strconv.Itoa(int(s.uuid)))}
}

func (s Inhabited) Eval(ctx Context) Sort {
	t := s.Type.Eval(ctx)
	switch t := t.(type) {
	case Lambda:
		return Lambda{
			Param: t.Param,
			Body: Inhabited{
				uuid: nextCount(),
				Type: t.Body,
			}.Eval(ctx),
		}
	default:
		return NewTerm(s.Form(), t)
	}
}

func init() {
	ListParseFuncMap[LambdaCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(LambdaCmd), "param1", "...", "paramN", "body"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		params := slices_util.Map(list[:len(list)-1], func(form Form) Annot {
			return compileAnnot(ctx, mustType[List](err, form))
		})
		output := ctx.Parse(list[len(list)-1])
		slices_util.ForEach(slices_util.Reverse(params), func(param Annot) {
			output = Lambda{
				Param: param,
				Body:  output,
			}
		})
		return output
	}
	const ArrowCmd Name = "->"
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) Sort {
		// make builtin like succ
		// e.g. if arrow is Nat -> Nat
		// then its lambda is
		// (x: Nat) => Nat
		// or some mechanism to introduce arrow type from pi type
		// TODO - probably we don't need this anymore
		panic("not implemented")
	}
}

const (
	LambdaCmd Name = "=>"
)

// Lambda - lambda abstraction (or Lambda-type)
type Lambda struct {
	Param Annot
	Body  Sort
}

func (s Lambda) Form() Form {
	return List{LambdaCmd, s.Param.Form(), s.Body.Form()}
}

func (s Lambda) Level(ctx Context) int {
	panic("not_implemented")
}

func (s Lambda) Parent(ctx Context) Sort {
	return Lambda{
		Param: s.Param,
		Body: Type{
			Value: s.Body,
		},
	}
}
func (s Lambda) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}

func init() {
	DefaultParseFunc = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{"cmd", "arg1", "...", "argN"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		output := ctx.Parse(list[0])
		args := slices_util.Map(list[1:], func(form Form) Sort {
			return ctx.Parse(form)
		})
		slices_util.ForEach(args, func(arg Sort) {
			output = Beta{
				Cmd: output,
				Arg: arg,
			}
		})
		return output
	}
}

type Beta struct {
	code
	Cmd Sort
	Arg Sort
}

func (s Beta) Form() Form {
	return List{s.Cmd.Form(), s.Arg.Form()}
}

func (s Beta) Eval(ctx Context) Sort {
	return s.Cmd.Body.Eval(ctx.Set(s.Cmd.Param.Name, s.Arg))
}
