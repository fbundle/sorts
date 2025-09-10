package ast

type AST interface {
	astAttr()
}

type Term string

func (t Term) astAttr() {}

type Abstraction struct {
	Param Term
	Body  AST
}

func (a Abstraction) astAttr() {}

type Application struct {
	Func AST
	Arg  AST
}

func (a Application) astAttr() {}

func AlphaConversion(ast AST, oldName Term, newName Term) AST {
	panic("not implemented")
}
func BetaReduction(ast Application) AST {
	panic("not implemented")
}
