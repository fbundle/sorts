package sorts5

import (
	"fmt"
)

const (
	ArrowCmd Name = "->"
)

var ArrowListParseFunc ListParseFunc[Context] = func(ctx Context, list List) (Context, Sort) {
	err := fmt.Errorf("arrow must be (%s type1 type2)", ArrowCmd)
	mustMatchHead(err, ArrowCmd, list)

}

type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) Compile(ctx Frame) Sort {
	s.A = s.A.Compile(ctx)
	s.B = s.B.Compile(ctx)
	return s
}

func (s Arrow) Form() Form {
	return List{ArrowCmd, s.A.Form(), s.B.Form()}
}

func (s Arrow) Level(ctx Frame) int {
	return max(s.A.Level(ctx), s.B.Level(ctx))
}

func (s Arrow) Parent(ctx Frame) Sort {
	return Arrow{
		A: s.A.Parent(ctx),
		B: s.B.Parent(ctx),
	}
}

func (s Arrow) LessEqual(ctx Frame, d Sort) bool {
	if d, ok := d.(Arrow); ok {
		return d.A.LessEqual(ctx, s.A) && s.B.LessEqual(ctx, d.B)
	}
	return ctx.LessEqual(s.Form(), d.Form())
}

var _ Sort = Arrow{}

type Prod struct {
}
