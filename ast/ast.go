package ast

type astAttr struct {
	parent AST // type
}

// AST - typed lambda calculus
type AST interface {
	astAttr() astAttr
}

func Parent(a AST) AST {
	return a.astAttr().parent
}

// Term - variable or constant
type Term string

func (t Term) astAttr() astAttr {
	panic("not implemented")
}

// LetBinding - let binding
type LetBinding struct {
	Name   string
	Parent AST
	Expr   AST
}
type Let struct {
	Parent   AST // eg. Nat
	Bindings []LetBinding
	Expr     AST
}

func (a Let) astAttr() astAttr {
	panic("not implemented")
}

type MatchCase struct {
	Comp AST
	Expr AST
}

// Match - match expression
type Match struct {
	Parent  AST
	Cond    AST
	Cases   []MatchCase
	Default AST
}

func (a Match) astAttr() astAttr {
	panic("not implemented")
}

// Lambda - lambda abstraction
type Lambda struct {
	Parent AST  // eg. (-> Nat Nat)
	Param  Term // eg. x
	Body   AST  // eg. (+ x 1)
}

func (a Lambda) astAttr() astAttr {
	panic("not implemented")
}

// Beta - beta reduction
type Beta struct {
	Func AST // must be Lambda or Term after evaluated
	Arg  AST
}

func (a Beta) astAttr() astAttr {
	panic("not implemented")
}
