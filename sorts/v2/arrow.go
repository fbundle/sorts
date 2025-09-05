package sorts

type arrow struct {
	param Sort
	pody  Sort
}

func (s arrow) Level(ss SortSystem) int {
	//TODO implement me
	panic("implement me")
}

func (s arrow) Name(ss SortSystem) string {
	//TODO implement me
	panic("implement me")
}

func (s arrow) Parent(ss SortSystem) Sort {
	//TODO implement me
	panic("implement me")
}

func (s arrow) LessEqual(ss SortSystem, dst Sort) bool {
	//TODO implement me
	panic("implement me")
}
