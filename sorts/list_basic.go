package sorts

import (
	"fmt"
)

const (
	ArrowCmd Name = "->"
)

func init() {
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) (Context, Sort) {
		err := fmt.Errorf("arrow must be (%s type1 type2)", ArrowCmd)
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
	A Sort
	B Sort
}

func (s Arrow) Compile(ctx Context) Sort {
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

var _ Sort = Arrow{}

const (
	ProdCmd Name = "⊗"
)

func init() {
	ListParseFuncMap[ProdCmd] = func(ctx Context, list List) (Context, Sort) {
		err := fmt.Errorf("prod must be (%s type1 type2)", ProdCmd)
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
	A Sort
	B Sort
}

func (s Prod) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Prod) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (s Prod) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Prod) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Prod) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Prod{}
