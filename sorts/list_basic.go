package sorts

import (
	"fmt"
)

const (
	ArrowCommand Name = "->"
)

var ArrowListParser = ListParser{
	Command: ArrowCommand,
	ListParse: func(ctx Context, list List) (Context, Sort) {
		err := fmt.Errorf("arrow must be (%s type1 type2)", ArrowCommand)

	},
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
	return List{ArrowCommand, s.A.Form(), s.B.Form()}
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

type Prod struct {
}
