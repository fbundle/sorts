package sorts

func NewChain(name Name, level int) Atom {
	return Atom{
		form: name,
		level: func(ctx Context) int {
			return level
		},
		parent: func(ctx Context) Sort {
			return NewChain(name, level+1)
		},
	}
}

func NewTerm(form Form, parent func(ctx Context) Sort) Atom {
	return Atom{
		form: form,
		level: func(ctx Context) int {
			return parent(ctx).Level(ctx) - 1
		},
		parent: parent,
	}
}

type Atom struct {
	form   Form
	level  func(ctx Context) int
	parent func(ctx Context) Sort
}

func (s Atom) Form() Form {
	return s.form
}

func (s Atom) Compile(ctx Context) Sort {
	// atom is created not compiled from
	return s
}

func (s Atom) Level(ctx Context) int {
	return s.level(ctx)
}

func (s Atom) Parent(ctx Context) Sort {
	return s.parent(ctx)
}

func (s Atom) LessEqual(ctx Context, d Sort) bool {
	return ctx.LessEqual(s.Form(), d.Form())
}

var _ Sort = Atom{}
