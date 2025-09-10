package ast

import "github.com/fbundle/sorts/expr"

type AST interface {
	Marshal() []expr.Token
	astAttr()
}

type Term string

func (t Term) Marshal() []expr.Token {
	panic("not_implemented")
}
func (t Term) astAttr() {}

type Abstraction struct {
	Param Term
	Body  AST
}

func (a Abstraction) Marshal() []expr.Token {
	panic("not_implemented")
}
func (a Abstraction) astAttr() {}
func (a Abstraction) AlphaConversion(param Term) Abstraction {
	panic("not_implemented")
}

type Application struct {
	Func AST
	Arg  AST
}

func (a Application) Marshal() []expr.Token {
	panic("not_implemented")
}
func (a Application) astAttr() {}
func (a Application) BetaReduction() AST {
	panic("not_implemented")
}
