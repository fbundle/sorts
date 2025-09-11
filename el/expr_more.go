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

func init() {
	RegisterListParser(":", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) != 2 {
			return nil, errors.New("define must have exactly 2 arguments: name and type")
		}
		name, ok := list[0].(form.Term)
		if !ok {
			return nil, errors.New("define name must be a term")
		}
		typeExpr, err := parseFunc(list[1])
		if err != nil {
			return nil, err
		}
		return Define{
			Name: Term(name),
			Type: typeExpr,
		}, nil
	})
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

func init() {
	RegisterListParser(":=", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) != 2 {
			return nil, errors.New("assign must have exactly 2 arguments: name and value")
		}
		name, ok := list[0].(form.Term)
		if !ok {
			return nil, errors.New("assign name must be a term")
		}
		value, err := parseFunc(list[1])
		if err != nil {
			return nil, err
		}
		return Assign{
			Name:  Term(name),
			Value: value,
		}, nil
	})
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

func init() {
	RegisterListParser("chain", func(parseFunc ParseFunc, list form.List) (Expr, error) {
		if len(list) < 2 {
			return nil, errors.New("chain must have at least 2 arguments: init expressions and tail")
		}
		// All arguments except the last are init expressions
		init := make([]Expr, len(list)-1)
		for i, arg := range list[:len(list)-1] {
			initExpr, err := parseFunc(arg)
			if err != nil {
				return nil, err
			}
			init[i] = initExpr
		}
		// Last argument is the tail
		tail, err := parseFunc(list[len(list)-1])
		if err != nil {
			return nil, err
		}
		return Chain{
			Init: init,
			Tail: tail,
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
