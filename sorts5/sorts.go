package sorts5

import (
	"github.com/fbundle/sorts/form2"
)

type Form = form2.Form

type Sort struct {
	Form      Form
	Level     func() int
	Parent    func() Sort
	LessEqual func(dst Sort) bool
}
