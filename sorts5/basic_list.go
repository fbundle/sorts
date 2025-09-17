package sorts5

const (
	ArrowCmd Name = "->"
)

type Arrow struct {
	A Sort
	B Sort
}

func (s Arrow) Form() Form {
	return List{ArrowCmd, s.A.Form(), s.B.Form()}
}

func (s Arrow) Level() int {
	//TODO implement me
	panic("implement me")
}

func (s Arrow) Parent() Sort {
	//TODO implement me
	panic("implement me")
}

func (s Arrow) LessEqual(dst Sort) bool {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Arrow{}
