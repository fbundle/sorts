package el

type Expr interface {
	mustExpr()
}

type Term string

func (t Term) mustExpr() {}

// Lambda - (or =>) lambda abstraction
type Lambda struct {
	Param Term
	Body  Expr
}

func (l Lambda) mustExpr() {}

// FunctionCall - (cmd arg1 arg2 ...) function call
type FunctionCall struct {
	Cmd  Expr
	Args []Expr
}

func (f FunctionCall) mustExpr() {}

// Define - (or :) define a new variable
type Define struct {
	Name Term
	Type Expr
}

func (d Define) mustExpr() {}

// Assign - (or :=) assign a value to a variable
type Assign struct {
	Name  Term
	Value Expr
}

func (a Assign) mustExpr() {}

type Let struct {
	Defines []Define
	Assigns []Assign
	Body    Expr
}

func (l Let) mustExpr() {}

type Case struct {
	Comp  Expr
	Value Expr
}
type Match struct {
	Cases []Case
	Else  Expr
}

func (m Match) mustExpr() {}
