package sorts5

func NewChain(name Name, level int) Atom {
	return Atom{
		form: name,
		level: func(ctx Frame) int {
			return level
		},
		parent: func(ctx Frame) Sort {
			return NewChain(name, level+1)
		},
	}
}

func NewAtom(form Form, level func(ctx Frame) int, parent func(ctx Frame) Sort) Atom {
	return Atom{
		form:   form,
		level:  level,
		parent: parent,
	}
}

type Atom struct {
	form   Form
	level  func(ctx Frame) int
	parent func(ctx Frame) Sort
}

func (s Atom) Form() Form {
	return s.form
}

func (s Atom) Compile(ctx Frame) Sort {
	// atom is created not compiled from
	return s
}

func (s Atom) Level(ctx Frame) int {
	return s.level(ctx)
}

func (s Atom) Parent(ctx Frame) Sort {
	return s.parent(ctx)
}

func (s Atom) LessEqual(ctx Frame, d Sort) bool {
	return FallbackLessEqual(s.Form(), d.Form())
}

var _ Sort = Atom{}
