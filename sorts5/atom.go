package sorts5

type Atom struct {
	form   Name
	level  func(ctx Context) int
	parent func(ctx Context) Sort
}

func (s Atom) Form() Form {
	return s.form
}

func (s Atom) Compile(ctx Context) Sort {
	// atom is created not compiled from
	return s
}

func (s Atom) Level(ctx Context) int {
	return s.level(ctx)
}

func (s Atom) Parent(ctx Context) Sort {
	return s.parent(ctx)
}

func (s Atom) LessEqual(ctx Context, dst Sort) bool {

}

var _ Sort = Atom{}
