package el

import (
	"errors"

	"github.com/fbundle/sorts/form"
)

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

func init() {
	RegisterListParser("=>", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) != 2 {
			return nil, errors.New("lambda must have exactly 2 arguments: param and body")
		}
		param, ok := list[0].(form.Term)
		if !ok {
			return nil, errors.New("lambda parameter must be a term")
		}
		body, err := parseFunc(list[1])
		if err != nil {
			return nil, err
		}
		return Lambda{
			Param: Term(param),
			Body:  body,
		}, nil
	})
}

// Let - (chain name1 type1 value1 ... namen typen valuen tail)
type Let struct {
	Bindings []Binding
	Final    Expr
}

type Binding struct {
	Name  Term
	Type  Expr
	Value Expr
}

func (l Let) mustExpr() {}

func (l Let) Marshal() form.Form {
	forms := make([]form.Form, 0, 2+3*len(l.Bindings))
	forms = append(forms, form.Term("chain"))
	for _, binding := range l.Bindings {
		forms = append(forms, binding.Name.Marshal())
		forms = append(forms, binding.Type.Marshal())
		forms = append(forms, binding.Value.Marshal())
	}
	forms = append(forms, l.Final.Marshal())
	return form.List(forms)
}

func init() {
	RegisterListParser("let", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) < 1 {
			return nil, errors.New("let must have at least 1: final")
		}
		if (len(list)-1)%3 != 0 {
			return nil, errors.New("let bindings must be in triplets: name, type, value")
		}
		bindings := make([]Binding, 0, (len(list)-1)/3)
		for i := 0; i < len(list)-1; i += 3 {
			name, ok := list[i].(form.Term)
			if !ok {
				return nil, errors.New("let binding name must be a term")
			}
			typ, err := parseFunc(list[i+1])
			if err != nil {
				return nil, err
			}
			val, err := parseFunc(list[i+2])
			if err != nil {
				return nil, err
			}
			bindings = append(bindings, Binding{
				Name:  Term(name),
				Type:  typ,
				Value: val,
			})
		}
		final, err := parseFunc(list[len(list)-1])
		if err != nil {
			return nil, err
		}
		return Let{
			Bindings: bindings,
			Final:    final,
		}, nil
	})
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

func init() {
	RegisterListParser("match", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) < 3 {
			return nil, errors.New("match must have at least 3 arguments: condition, cases, and final")
		}
		// First argument is the condition
		cond, err := parseFunc(list[0])
		if err != nil {
			return nil, err
		}
		// Parse cases - they come in pairs (comp, value)
		cases := make([]Case, 0)
		remainingArgs := list[1:]
		// Check if we have an odd number of remaining arguments (final case)
		var final Expr
		if len(remainingArgs)%2 != 1 {
			return nil, errors.New("match must have a final case")
		}
		// Last argument is the final/default case
		final, err = parseFunc(remainingArgs[len(remainingArgs)-1])
		if err != nil {
			return nil, err
		}
		remainingArgs = remainingArgs[:len(remainingArgs)-1]
		// Parse case pairs
		for i := 0; i < len(remainingArgs); i += 2 {
			comp, err := parseFunc(remainingArgs[i])
			if err != nil {
				return nil, err
			}
			value, err := parseFunc(remainingArgs[i+1])
			if err != nil {
				return nil, err
			}
			cases = append(cases, Case{
				Comp:  comp,
				Value: value,
			})
		}
		return Match{
			Cond:  cond,
			Cases: cases,
			Final: final,
		}, nil
	})
}

type Arrow struct {
	A Expr
	B Expr
}

func (a Arrow) mustExpr() {}
func (a Arrow) Marshal() form.Form {
	return form.List{
		form.Term("->"),
		a.A.Marshal(),
		a.B.Marshal(),
	}
}

func init() {
	RegisterListParser("->", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) != 2 {
			return nil, errors.New("arrow must have exactly 2 arguments: a and b")
		}
		a, err := parseFunc(list[0])
		if err != nil {
			return nil, err
		}
		b, err := parseFunc(list[1])
		if err != nil {
			return nil, err
		}
		return Arrow{
			A: a,
			B: b,
		}, nil
	})
}

type Sum struct {
	A Expr
	B Expr
}

func (s Sum) mustExpr() {}
func (s Sum) Marshal() form.Form {
	return form.List{
		form.Term("⊕"),
		s.A.Marshal(),
		s.B.Marshal(),
	}
}
func init() {
	RegisterListParser("⊕", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) != 2 {
			return nil, errors.New("arrow must have exactly 2 arguments: a and b")
		}
		a, err := parseFunc(list[0])
		if err != nil {
			return nil, err
		}
		b, err := parseFunc(list[1])
		if err != nil {
			return nil, err
		}
		return Sum{
			A: a,
			B: b,
		}, nil
	})
}
