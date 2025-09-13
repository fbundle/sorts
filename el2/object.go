package el2

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	TypeCheck(sa SortAttr, parent Sort) Sort
}

// Object - a sort
type Object struct {
	Sort
}

func (s Object) TypeCheck(sa SortAttr, parent Sort) Sort {
	must(sa).termOf(s, parent)
	return s.Sort
}

// Beta - beta reduction
type Beta struct {
	Cmd Form
	Arg Form
}

func (f Beta) TypeCheck(sa SortAttr, parent Sort) Sort {
	//TODO implement me
	panic("implement me")
}

// Lambda - lambda abstraction
type Lambda struct {
	Param Name
	Body  Form
}

func (l Lambda) TypeCheck(sa SortAttr, parent Sort) Sort {
	//TODO implement me
	panic("implement me")
}
