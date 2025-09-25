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

func NewTerm(form Form, parent Sort) Atom {
	return Atom{
		form: form,
		level: func(ctx Context) int {
			return parent.Level(ctx) - 1
		},
		parent: func(ctx Context) Sort {
			return parent
		},
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

func (s Atom) Level(ctx Context) int {
	return s.level(ctx)
}

func (s Atom) Parent(ctx Context) Sort {
	return s.parent(ctx)
}

func (s Atom) LessEqual(ctx Context, d Sort) bool {
	return ctx.LessEqual(s.Form(), d.Form())
}

const (
	PiCmd Name = "=>" // "Π"
)

type Pi struct {
	Param Annot
	Body  Code
}

func (s Pi) Form() Form {
	return List{PiCmd, s.Param.Form(), s.Body.Form()}
}

func (s Pi) Level(ctx Context) int {
	// TODO - hard - probably need to include a "proof of level"
	panic("not_implemented")
}

func (s Pi) Parent(ctx Context) Sort {
	return Pi{
		Param: s.Param,
		Body: Type{
			Value: s.Body,
		},
	}
}
func (s Pi) LessEqual(ctx Context, d Sort) bool {
	// TODO - probably use contravariant-covariant rules
	panic("not_implemented")
}

// Eval - Pi is both a Sort (Pi-type) and a Code (lambda abstraction)
func (s Pi) Eval(ctx Context) Sort {
	return s
}

const (
	SigmaCmd Name = "×"
)

type Sigma struct {
	Param Annot
	Body  Code
}

func (s Sigma) Form() Form {
	return List{SigmaCmd, s.Param.Form(), s.Body.Form()}
}

func (s Sigma) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (s Sigma) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (s Sigma) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

func (s Sigma) Eval(ctx Context) Sort {
	return s
}
