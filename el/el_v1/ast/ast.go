package ast

type Expr interface {
	AST()
}

type Name string // can be a name, an integer or a string (string starts with double quote)

func (n Name) AST() {}

type LetBinding struct {
	Name Name
	Expr Expr
}
type Let struct {
	Bindings []LetBinding
	Final    Expr
}

func (l Let) AST() {}

type MatchCase struct {
	Comp  Expr
	Value Expr
}

type Match struct {
	Cond    Expr
	Cases   []MatchCase
	Default Expr
}

func (m Match) AST() {}

type Lambda struct {
	Param Name
	Body  Expr
}

func (l Lambda) AST() {}

type FunctionCall struct {
	Cmd  Name
	Args []Expr
}

func (f FunctionCall) AST() {}
