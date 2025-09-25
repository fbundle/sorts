package sorts

import "strconv"

func init() {
	ListParseFuncMap[TypeCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(TypeCmd), "value"})
		if len(list) != 1 {
			panic(err)
		}
		return Type{
			Body: ctx.Parse(list[0]),
		}
	}
}

const (
	TypeCmd Name = "&"
)

type Type struct {
	Body Sort
}

func (s Type) Form() Form {
	return List{TypeCmd, s.Body.Form()}
}

func (s Type) Parent(ctx Context) Sort {
	return s.Body.Parent(ctx).Parent(ctx)
}

func (s Type) Level(ctx Context) int {
	return s.Body.Level(ctx) + 1
}
func (s Type) LessEqual(ctx Context, d Sort) bool {
	return s.Body.Parent(ctx).LessEqual(ctx, d)
}

func (s Type) Eval(ctx Context) Sort {
	return s.Body.Parent(ctx)
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
	uuid uint64
	Type Sort
}

func (s Inhabited) Form() Form {
	return List{InhabitCmd, s.Type.Form(), Name(strconv.Itoa(int(s.uuid)))}
}

func (s Inhabited) Parent(ctx Context) Sort {
	return s.Type
}

func (s Inhabited) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Inhabited) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
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
