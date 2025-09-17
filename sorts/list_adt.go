package sorts

const (
	ArrowCmd Name = "->"
)

func init() {
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) (Context, Sort1) {
		err := parseErr(ArrowCmd, []string{"type1", "type2"})

		if len(list) != 2 {
			panic(err)
		}
		ctx, a := ctx.Parse(list[0])
		ctx, b := ctx.Parse(list[1])
		return ctx, Arrow{
			A: a,
			B: b,
		}
	}
}

type Arrow struct {
	A Sort1
	B Sort1
}

func (s Arrow) Compile(ctx Context) Sort1 {
	s.A = s.A.Compile(ctx)
	s.B = s.B.Compile(ctx)
	return s
}

func (s Arrow) Form() Form {
	return List{ArrowCmd, s.A.Form(), s.B.Form()}
}

func (s Arrow) Level(ctx Context) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Arrow) Parent(ctx Context) Sort1 {
	return Arrow{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Arrow) LessEqual(ctx Context, d Sort1) bool {
	if d, ok := d.(Arrow); ok {
		return d.A.LessEqual(ctx, s.A) && s.B.LessEqual(ctx, d.B)
	}
	return ctx.LessEqual(s.Form(), d.Form())
}

var _ Sort1 = Arrow{}

const (
	ProdCmd Name = "⊗"
)

func init() {
	ListParseFuncMap[ProdCmd] = func(ctx Context, list List) (Context, Sort1) {
		err := parseErr(ProdCmd, []string{"type1", "type2"})
		if len(list) != 2 {
			panic(err)
		}
		ctx, a := ctx.Parse(list[0])
		ctx, b := ctx.Parse(list[1])
		return ctx, Prod{
			A: a,
			B: b,
		}
	}
}

type Prod struct {
	A Sort1
	B Sort1
}

func (s Prod) Compile(ctx Context) Sort1 {
	s.A = s.A.Compile(ctx)
	s.B = s.B.Compile(ctx)
	return s
}

func (s Prod) Form() Form {
	return List{ProdCmd, s.A.Form(), s.B.Form()}
}

func (s Prod) Level(ctx Context) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Prod) Parent(ctx Context) Sort1 {
	return Prod{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Prod) LessEqual(ctx Context, d Sort1) bool {
	if d, ok := d.(Prod); ok {
		return s.A.LessEqual(ctx, d.A) && s.B.LessEqual(ctx, d.B)
	}
	return ctx.LessEqual(s.Form(), d.Form())
}

var _ Sort1 = Prod{}

const (
	SumCmd Name = "⊕"
)

func init() {
	ListParseFuncMap[SumCmd] = func(ctx Context, list List) (Context, Sort1) {
		err := parseErr(SumCmd, []string{"type1", "type2"})
		if len(list) != 2 {
			panic(err)
		}
		ctx, a := ctx.Parse(list[0])
		ctx, b := ctx.Parse(list[1])
		return ctx, Sum{
			A: a,
			B: b,
		}
	}
}

type Sum struct {
	A Sort1
	B Sort1
}

func (s Sum) Compile(ctx Context) Sort1 {
	s.A = s.A.Compile(ctx)
	s.B = s.B.Compile(ctx)
	return s
}

func (s Sum) Form() Form {
	return List{SumCmd, s.A.Form(), s.B.Form()}
}

func (s Sum) Level(ctx Context) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Sum) Parent(ctx Context) Sort1 {
	return Sum{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Sum) LessEqual(ctx Context, d Sort1) bool {
	// interesting - (A + B) is the least upper bound of A and B
	// hence (A + B) <= C iff A <= C and B <= C
	return s.A.LessEqual(ctx, d) && s.B.LessEqual(ctx, d)
}
