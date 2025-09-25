package sorts

import (
	"github.com/fbundle/sorts/slices_util"
)

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

func (s Atom) Eval(ctx Context) Sort {
	return s
}

func init() {
	ListParseFuncMap[PiCmd] = func(ctx Context, list List) Sort {
		err := compileErr(list, []string{string(PiCmd), "param1", "...", "paramN", "body"}, "where N >= 0")
		if len(list) != 2 {
			panic(err)
		}
		params := slices_util.Map(list[:len(list)-1], func(form Form) Annot {
			return compileAnnot(ctx, mustType[List](err, form))
		})
		output := ctx.Parse(list[len(list)-1])
		slices_util.ForEach(slices_util.Reverse(params), func(param Annot) {
			output = Pi{
				Param: param,
				Body:  output,
			}
		})
		return output
	}
	const ArrowCmd Name = "->"
	ListParseFuncMap[ArrowCmd] = func(ctx Context, list List) Sort {
		// make builtin like succ
		// e.g. if arrow is Nat -> Nat
		// then its lambda is
		// (x: Nat) => Nat
		// or some mechanism to introduce arrow type from pi type
		// TODO - probably we don't need this anymore
		panic("not implemented")
	}
}

const (
	PiCmd Name = "=>"
)

// Pi - lambda abstraction (or Pi-type)
type Pi struct {
	Param Annot
	Body  Sort
}

func (s Pi) Form() Form {
	return List{PiCmd, s.Param.Form(), s.Body.Form()}
}

func (s Pi) Level(ctx Context) int {
	panic("not_implemented")
}

func (s Pi) Parent(ctx Context) Sort {
	return Pi{
		Param: s.Param,
		Body: Type{
			Body: s.Body,
		},
	}
}
func (s Pi) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}

func (s Pi) Eval(ctx Context) Sort {
	return s
}
