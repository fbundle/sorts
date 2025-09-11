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

// Chain - chain expressions
type Chain struct {
	Init []Expr
	Tail Expr
}

func (c Chain) mustExpr() {}

// Match - match expression
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
