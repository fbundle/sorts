package ast

// AST - typed lambda calculus
type AST interface {
	astAttr()
}

type Term struct {
	Name string
	Type AST
}

func (t Term) astAttr() {}

type Abstraction struct {
	Type  AST
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
