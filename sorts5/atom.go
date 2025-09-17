package sorts5

func NewChain(name Name, level int) Atom {
	return Atom{
		form: name,
		level: func(frame Frame) int {
			return level
		},
		parent: func(frame Frame) Sort {
			return NewChain(name, level+1)
		},
	}
}

func NewTerm(form Form, parent func(frame Frame) Sort) Atom {
	return Atom{
		form: form,
		level: func(frame Frame) int {
			return parent(frame).Level(frame) - 1
		},
		parent: parent,
	}
}

type Atom struct {
	form   Form
	level  func(frame Frame) int
	parent func(frame Frame) Sort
}

func (s Atom) Form() Form {
	return s.form
}

func (s Atom) Compile(frame Frame) Sort {
	// atom is created not compiled from
	return s
}

func (s Atom) Level(frame Frame) int {
	return s.level(frame)
}

func (s Atom) Parent(frame Frame) Sort {
	return s.parent(frame)
}

func (s Atom) LessEqual(frame Frame, d Sort) bool {
	return FallbackLessEqual(s.Form(), d.Form())
}

var _ Sort = Atom{}
