package el2

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	Sort() Sort // nullable
	TypeCheck(sa SortAttr, parent Sort) Sort
}

// Object - a sort
type Object struct {
	sort Sort
}

func (s Object) Sort() Sort {
	return s.sort
}

func (s Object) TypeCheck(sa SortAttr, parent Sort) Sort {
	must(sa).termOf(s.sort, parent)
	return s.sort
}

// Beta - beta reduction
type Beta struct {
	Cmd Form
	Arg Form
}

func (f Beta) Sort() Sort {
	return nil
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

func (l Lambda) Sort() Sort {
	return nil
}

func (l Lambda) TypeCheck(sa SortAttr, parent Sort) Sort {
	//TODO implement me
	panic("implement me")
}
