package sorts

type Beta struct {
	Cmd  Sort
	Args []Sort
}

func (b Beta) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Beta{}
