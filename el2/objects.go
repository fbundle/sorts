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

func (f Beta) TypeCheck(sa SortAttr, parent Sort) Sort {
	//TODO implement me
	panic("implement me")
}

func ListParseLambda(parse ParseFunc, list List) AlmostSort {
	if len(list) != 2 {
		panic("lambda list must have two elements")
	}
	return Lambda{
		Param: list[0].(Name),
		Body:  parse(list[1]),
	}
}

// Lambda - lambda abstraction
type Lambda struct {
	Param Name
	Body  AlmostSort
}

func (l Lambda) TypeCheck(sa SortAttr, parent Sort) Sort {
	//TODO implement me
	panic("implement me")
}

type LetBinding struct {
	Name  Name
	Value AlmostSort
}
type Let struct {
	Bindings []LetBinding
	Final    AlmostSort
}

func (l Let) TypeCheck(sa SortAttr, parent Sort) Sort {
	panic("implement me")
}

type MatchCase struct {
	Pattern Form
	Value   AlmostSort
}
type Match struct {
	Cond  AlmostSort
	Cases []MatchCase
	Final AlmostSort
}

func (m Match) TypeCheck(sa SortAttr, parent Sort) Sort {
	panic("implement me")
}
