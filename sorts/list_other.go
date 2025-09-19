package sorts

const (
	InspectCmd = "inspect"
)

func init() {
	ListCompileFuncMap[InspectCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{InspectCmd, "value"})
		if len(list) != 1 {
			panic(err)
		}
		return Inspect{
			Value: ctx.Compile(list[0]).TypeCheck(ctx),
		}
	}
}

type Inspect struct {
	Value Sort
}

func (s Inspect) Form() Form {
	return List{Name(InspectCmd), s.Value.Form()}
}

func (s Inspect) TypeCheck(ctx Context) Sort {
	return Inspect{
		Value: s.Value.TypeCheck(ctx),
	}
}

func (s Inspect) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Inspect) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Inspect) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Inspect) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

const (
	TypeCmd = "type"
)

func init() {
	ListCompileFuncMap[TypeCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{TypeCmd, "value"})
		if len(list) != 1 {
			panic(err)
		}
		return Type{
			Value: ctx.Compile(list[0]).TypeCheck(ctx),
		}
	}
}

type Type struct {
	Value Sort
}

func (s Type) Form() Form {
	return List{Name(TypeCmd), s.Value.Form()}
}

func (s Type) TypeCheck(ctx Context) Sort {
	return Inspect{
		Value: s.Value.TypeCheck(ctx),
	}
}

func (s Type) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Type) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Type) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Type) Reduce(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}
