package el

type Expr interface {
	mustElExpr()
}

type Term string

func (t Term) mustElExpr() {}

// Define - (or :) define a new variable
type Define struct {
	Name Term
	Type Expr
}

func (d Define) mustElExpr() {}

// Assign - (or :=) assign a value to a variable
type Assign struct {
	Name  Term
	Value Expr
}

func (a Assign) mustElExpr() {}

// Lambda - (or =>) lambda abstraction
type Lambda struct {
	Param Term
	Body  Expr
}

func (l Lambda) mustElExpr() {}

// FunctionCall - (cmd arg1 arg2 ...) function call
type FunctionCall struct {
	Cmd  Expr
	Args []Expr
}

func (f FunctionCall) mustElExpr() {}

type Let struct {
	Defines []Define
	Assigns []Assign
	Body    Expr
}
