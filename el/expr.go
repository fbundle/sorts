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

// Lambda - (=> param body)
type Lambda struct {
	Param Term
	Body  Expr
}

func (l Lambda) mustExpr() {}

func (l Lambda) Marshal() form.Form {
	return form.List{
		form.Term("=>"),
		l.Param.Marshal(),
		l.Body.Marshal(),
	}
}

// Define - (: name type)
type Define struct {
	Name Term
	Type Expr
}

func (d Define) mustExpr() {}

func (d Define) Marshal() form.Form {
	return form.List{
		form.Term(":"),
		d.Name.Marshal(),
		d.Type.Marshal(),
	}
}

// Assign - (:= name value)
type Assign struct {
	Name  Term
	Value Expr
}

func (a Assign) mustExpr() {}

func (a Assign) Marshal() form.Form {
	return form.List{
		form.Term(":="),
		a.Name.Marshal(),
		a.Value.Marshal(),
	}
}

// Chain - (chain expr1 expr2 ... exprn tail)
type Chain struct {
	Init []Expr
	Tail Expr
}

func (c Chain) mustExpr() {}

func (c Chain) Marshal() form.Form {
	forms := make([]form.Form, 0, 2+len(c.Init))
	forms = append(forms, form.Term("chain"))
	for _, expr := range c.Init {
		forms = append(forms, expr.Marshal())
	}
	forms = append(forms, c.Tail.Marshal())
	return form.List(forms)
}

// Match - (match cond comp1 value1 comp2 value2 ... compn valuen final)
type Match struct {
	Cond  Expr
	Cases []Case
	Final Expr
}
type Case struct {
	Comp  Expr
	Value Expr
}

func (m Match) mustExpr() {}

func (m Match) Marshal() form.Form {
	forms := make([]form.Form, 0, 2+2*len(m.Cases)+1)
	forms = append(forms, form.Term("match"))
	forms = append(forms, m.Cond.Marshal())
	for _, c := range m.Cases {
		forms = append(forms, c.Comp.Marshal())
		forms = append(forms, c.Value.Marshal())
	}
	forms = append(forms, m.Final.Marshal())
	return form.List(forms)
}
