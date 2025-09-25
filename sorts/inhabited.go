package sorts

func init() {
	ListParseFuncMap[InhabitCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(InhabitCmd), "type"})
		if len(list) != 1 {
			panic(err)
		}
		return Inhabited{
			Type: ctx.Parse(list[0]),
		}
	}
}

const (
	InhabitCmd Name = "*"
)

type Inhabited struct {
	Type Sort
}

func (s Inhabited) Form() Form {
	return List{InhabitCmd, s.Type.Form()}
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
	return NewTerm(s.Form(), s.Type)
}
