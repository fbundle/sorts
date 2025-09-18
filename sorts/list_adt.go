package sorts

const (
	ArrowCmd Name = "->"
)

func init() {
	ListCompileFuncMap[ArrowCmd] = func(ctx Context, list List) Sort {
		err := compileErr(ArrowCmd, []string{"type1", "type2"})

		if len(list) != 2 {
			panic(err)
		}
		return Arrow{
			A: ctx.Compile(list[0]),
			B: ctx.Compile(list[1]),
		}
	}
}

type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) Form() Form {
	return List{ArrowCmd, s.A.Form(), s.B.Form()}
}

func (s Arrow) Level(ctx Context) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Arrow) Parent(ctx Context) Sort {
	return Arrow{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Arrow) LessEqual(ctx Context, d Sort) bool {
	if d, ok := d.(Arrow); ok {
		return d.A.LessEqual(ctx, s.A) && s.B.LessEqual(ctx, d.B)
	}
	return ctx.LessEqual(s.Form(), d.Form())
}

func (s Arrow) Reduce(ctx Context) Sort {
	panic("implement me")
}

var _ Sort = Arrow{}

const (
	ProdCmd Name = "⊗"
)

func init() {
	ListCompileFuncMap[ProdCmd] = func(ctx Context, list List) Sort {
		err := compileErr(ProdCmd, []string{"type1", "type2"})
		if len(list) != 2 {
			panic(err)
		}
		return Prod{
			A: ctx.Compile(list[0]),
			B: ctx.Compile(list[1]),
		}
	}
}

type Prod struct {
	A Sort
	B Sort
}

func (s Prod) Form() Form {
	return List{ProdCmd, s.A.Form(), s.B.Form()}
}
func (s Prod) Level(ctx Context) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Prod) Parent(ctx Context) Sort {
	return Prod{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Prod) LessEqual(ctx Context, d Sort) bool {
	if d, ok := d.(Prod); ok {
		return s.A.LessEqual(ctx, d.A) && s.B.LessEqual(ctx, d.B)
	}
	return ctx.LessEqual(s.Form(), d.Form())
}

func (s Prod) Reduce(ctx Context) Sort {
	panic("implement me")
}

var _ Sort = Prod{}

const (
	SumCmd Name = "⊕"
)

func init() {
	ListCompileFuncMap[SumCmd] = func(ctx Context, list List) Sort {
		err := compileErr(SumCmd, []string{"type1", "type2"})
		if len(list) != 2 {
			panic(err)
		}
		return Sum{
			A: ctx.Compile(list[0]),
			B: ctx.Compile(list[1]),
		}
	}
}

type Sum struct {
	A Sort
	B Sort
}

func (s Sum) Form() Form {
	return List{SumCmd, s.A.Form(), s.B.Form()}
}

func (s Sum) Level(ctx Context) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Sum) Parent(ctx Context) Sort {
	return Sum{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Sum) LessEqual(ctx Context, d Sort) bool {
	// interesting - (A + B) is the least upper bound of A and B
	// hence (A + B) <= C iff A <= C and B <= C
	return s.A.LessEqual(ctx, d) && s.B.LessEqual(ctx, d)
}

func (s Sum) Reduce(ctx Context) Sort {
	panic("implement me")
}
