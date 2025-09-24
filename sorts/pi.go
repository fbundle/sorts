package sorts

const (
	PiCmd Name = "Π"
)

// Pi - lambda abstraction
type Pi struct {
	Param Annot
	Body  Sort
}

func (s Pi) Form() Form {
	return List{PiCmd, s.Param.Form(), s.Body.Form()}
}

func (s Pi) Level(ctx Context) int {
	panic("not_implemented")
}

func (s Pi) Parent(ctx Context) Sort {
	return Pi{
		Param: s.Param,
		Body: Type{
			Body: s.Body,
		},
	}
}
func (s Pi) LessEqual(ctx Context, d Sort) bool {
	panic("not_implemented")
}
func (s Pi) Reduce(ctx Context) Sort {
	return s
}
