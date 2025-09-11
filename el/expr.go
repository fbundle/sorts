package el

import "github.com/fbundle/sorts/form"

type Expr interface {
	Marshal() form.Form
	mustExpr()
}

type Term string

func (t Term) mustExpr() {}

func (t Term) Marshal() form.Form {
	return form.Term(t)
}

// FunctionCall - (cmd arg1 arg2 ...)
type FunctionCall struct {
	Cmd  Expr
	Args []Expr
}

func (f FunctionCall) mustExpr() {}

func (f FunctionCall) Marshal() form.Form {
	forms := make([]form.Form, 0, 1+len(f.Args))
	forms = append(forms, f.Cmd.Marshal())
	for _, arg := range f.Args {
		forms = append(forms, arg.Marshal())
	}
	return form.List(forms)
}
