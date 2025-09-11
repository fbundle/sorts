package el

type Expr interface {
	mustExpr()
}

type Term string

func (t Term) mustExpr() {}

// Lambda - (=> param body)
type Lambda struct {
	Param Term
	Body  Expr
}

func (l Lambda) mustExpr() {}

// FunctionCall - (cmd arg1 arg2 ...)
type FunctionCall struct {
	Cmd  Expr
	Args []Expr
}

func (f FunctionCall) mustExpr() {}

// Define - (: name type)
type Define struct {
	Name Term
	Type Expr
}

func (d Define) mustExpr() {}

// Assign - (:= name value)
type Assign struct {
	Name  Term
	Value Expr
}

func (a Assign) mustExpr() {}

// Chain - (chain expr1 expr2 ... exprn tail)
type Chain struct {
	Init []Expr
	Tail Expr
}

func (c Chain) mustExpr() {}

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
