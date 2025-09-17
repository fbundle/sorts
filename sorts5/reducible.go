package sorts5

type Beta struct {
	Cmd  Sort
	Args []Sort
}

func (b Beta) Compile(frame Frame) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Level(frame Frame) int {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Parent(frame Frame) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) LessEqual(frame Frame, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Beta{}
