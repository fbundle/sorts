package sorts5

import "github.com/fbundle/sorts/form"

type Form = form.Form
type Name = form.Name

type Sort struct {
	Form      Form
	Level     func() int
	Parent    func() Form
	LessEqual func(dst Sort) bool
}
