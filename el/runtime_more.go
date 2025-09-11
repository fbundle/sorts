package el

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/sorts"
)

var defaultFrame Frame

func init() {
	defaultFrame = defaultFrame.SetExec("let", func(frame Frame, argList form.List) (Frame, sorts.Sort, form.Form, error) {
		for i:=2; i< len(list)
	})
}
