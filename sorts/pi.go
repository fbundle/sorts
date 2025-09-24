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

func (s Pi) Compile(ctx Context) Sort           {}
func (s Pi) Level(ctx Context) int              {}
func (s Pi) Parent(ctx Context) Sort            {}
func (s Pi) LessEqual(ctx Context, d Sort) bool {}
