package sorts

import "fmt"

const (
	BetaCmd Name = "beta"
)

func init() {
	ListParseFuncMap[BetaCmd] = func(ctx Context, list List) (Context, Sort) {
		err := fmt.Errorf("beta must be (%s cmd arg1 ... argN)", BetaCmd)
	}
}

type Beta struct {
	Cmd  Sort
	Args []Sort
}

func (b Beta) Compile(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Form() Form {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Level(ctx Context) int {
	//TODO implement me
	panic("implement me")
}

func (b Beta) Parent(ctx Context) Sort {
	//TODO implement me
	panic("implement me")
}

func (b Beta) LessEqual(ctx Context, d Sort) bool {
	//TODO implement me
	panic("implement me")
}

var _ Sort = Beta{}
