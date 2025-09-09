package ast

type Expr interface{}

type Name string // can be a name, an integer or a string (string starts with double quote)

type Binding struct {
	Name Name
	Expr Expr
}
type Let struct {
	Bindings []Binding
	Final    Expr
}

type Case struct {
	Comp  Expr
	Value Expr
}
type Match struct {
	Cond    Expr
	Cases   []Case
	Default Expr
}

type Lambda struct {
	Param Name
	Body  Expr
}

type Call struct {
	Cmd  Name
	Args []Expr
}
