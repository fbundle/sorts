package sorts

const (
	AnnotCmd Name = ":"
)

func init() {
	ListParseFuncMap[AnnotCmd] = func(ctx Context, list List) Sort {
		err := parseErr(AnnotCmd, []string{"name", "type"})
		if len(list) != 2 {
			panic(err)
		}

		return Annot{
			Name: mustName(err, list[0]),
			Type: ctx.Parse(list[1]),
		}
	}
}

type Annot struct {
	Name Name
	Type Sort
}

func (s Annot) Form() Form {
	return List{AnnotCmd, s.Name, s.Type.Form()}
}

func (s Annot) Compile(ctx Context) Sort {
	return Annot{
		Name: s.Name,
		Type: s.Type.Compile(ctx),
	}
}

func (s Annot) Level(ctx Context) int {
	return s.Type.Level(ctx)
}

func (s Annot) Parent(ctx Context) Sort {
	return s.Type.Parent(ctx)
}

func (s Annot) LessEqual(ctx Context, d Sort) bool {
	return s.Type.LessEqual(ctx, d)
}

func (s Annot) Reduce(ctx Context) Sort {
	return Annot{
		Name: s.Name,
		Type: s.Type.Reduce(ctx),
	}
}

var _ Sort = Annot{}
