package el2

func mustSort(as AlmostSort) Sort {
	return as.(ActualSort).sort
}

// AlmostSort - almost a sort - for example, a lambda
type AlmostSort interface {
	TypeCheck(sa SortAttr, parent Sort) Sort // not nullable
}

// ActualSort - a sort
type ActualSort struct {
	sort Sort
}

func (s ActualSort) MustSort() Sort {
	return s.sort
}

func (s ActualSort) TypeCheck(sa SortAttr, parent Sort) Sort {
	must(sa).termOf(s.sort, parent)
	return s.sort
}

func ListParseBeta(parse ParseFunc, list List) AlmostSort {
	if len(list) != 2 {
		panic("beta list must have two elements")
	}
	return Beta{
		Cmd: parse(list[0]),
		Arg: parse(list[1]),
	}
}

// Beta - beta reduction
type Beta struct {
	Cmd AlmostSort
	Arg AlmostSort
}

func (f Beta) MustSort() Sort {
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

func (l Lambda) MustSort() Sort {
	return nil
}

func (l Lambda) TypeCheck(sa SortAttr, parent Sort) Sort {
	//TODO implement me
	panic("implement me")
}

type Let struct {
}
type Match struct {
}
