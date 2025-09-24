package sorts

func init() {
	ListParseFuncMap[PiCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(PiCmd), "param", "body"})
		if len(list) != 2 {
			panic(err)
		}
		return Lambda{
			Param: compileAnnot(ctx, mustType[List](err, list[0])),
			Body:  ctx.Parse(list[1]),
		}
	}
}

const (
	PiCmd Name = "=>"
)

// Lambda - lambda abstraction (or Pi-type)
type Lambda struct {
	Param Annot
	Body  Sort
}

func (s Lambda) Form() Form {
	return List{PiCmd, s.Param.Form(), s.Body.Form()}
}

func (s Lambda) Level(ctx Context) int {
	panic("not_implemented")
}

func (s Lambda) Parent(ctx Context) Sort {
	return Lambda{
		Param: s.Param,
		Body: Type{
			Body: s.Body,
		},
	}
}
func (s Lambda) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}
func (s Lambda) Reduce(ctx Context) Sort {
	return s
}
