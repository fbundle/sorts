package sorts5

type Atom struct {
	form   Form
	level  func() int
	parent func() Sort
}

func (a Atom) Form() Form {
	return a.form
}

func (a Atom) Level() int {
	return a.level()
}

func (a Atom) Parent() Sort {
	return a.parent()
}

func (a Atom) LessEqual(dst Sort) bool {

}
