package sorts

// Dept - represent a type B(x) depends on Sort x
type Dept struct {
	Form   Form
	Lambda Lambda
}

/*
// Inductive - inductive type
type Inductive interface {
	Sort
	Iter(yield func(name form.Name, constr func([]Sort) Inductive) bool)
}
*/

// Pi - dependent function type Π_{x: A} B(x)
type Pi struct {
	H Name
	A Sort
	B Dept
}

func (s Pi) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form: List{s.H, a.Form(s.A), a.Form(s.B)},
	}
}
