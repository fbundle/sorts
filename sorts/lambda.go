package sorts

func init() {
	ListParseFuncMap[LambdaCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(LambdaCmd), "param", "body"})
		if len(list) != 2 {
			panic(err)
		}
		return Lambda{
			Param: compileAnnot(ctx, mustType[List](err, list[0])),
			Body:  ctx.Parse(list[1]),
		}
	}
	const ArrowCmd Name = "->"
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) Sort {
		// make builtin like succ
		// e.g. if arrow is Nat -> Nat
		// then its lambda is
		// (x: Nat) => Nat
		// or some mechanism to introduce arrow type from pi type
		panic("not implemented")
	}
}

const (
	LambdaCmd Name = "=>"
)

// Lambda - lambda abstraction (or Pi-type)
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
			Body: s.Body,
		},
	}
}
func (s Lambda) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}
func (s Lambda) Eval(ctx Context) Sort {
	return s
}
