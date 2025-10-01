package sorts

import "github.com/fbundle/sorts/slices_util"

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
	Name   Form // alternative Name
	Params []Annot
	Body   Code
}

func (s Pi) Form() Form {
	if s.Name != nil {
		return s.Name
	}
	var output List
	output = append(output, PiCmd)
	output = append(output, slices_util.Map(s.Params, func(param Annot) Form {
		return param.Form()
	})...)
	output = append(output, s.Body.Form())
	return output
}

func (s Pi) Level(ctx Context) int {
	// TODO - hard - probably need to include a "proof of level"
	panic("not_implemented")
}

func (s Pi) Parent(ctx Context) Sort {
	return Pi{
		Params: s.Params,
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
	Name   Form
	Params []Annot
	Body   Code
}

func (s Sigma) Form() Form {
	if s.Name != nil {
		return s.Name
	}
	var output List
	output = append(output, SigmaCmd)
	output = append(output, slices_util.Map(s.Params, func(param Annot) Form {
		return param.Form()
	})...)
	output = append(output, s.Body.Form())
	return output
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

// Inductive - TODO
type Inductive struct {
	Name Form
	Type Annot   // Nat: Any_2
	Cons []Annot // {nil, (x: Nat)} ->
}
