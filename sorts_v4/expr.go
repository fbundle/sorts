package sorts

type Expr interface{}

type Term string

type Node struct {
	Cmd  Expr
	Args Expr
}
