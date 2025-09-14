package frame

import (
	"github.com/fbundle/sorts/form"
	"github.com/fbundle/sorts/persistent/ordered_map"
	"github.com/fbundle/sorts/sorts"
)

type Frame struct {
	frame ordered_map.OrderedMap[form.Name, sorts.Sort]
}

func (f Frame) Get(name form.Name) (sorts.Sort, bool) {
	return f.frame.Get(name)
}

func (f Frame) Set(name form.Name, sort sorts.Sort) Frame {
	return Frame{
		frame: f.frame.Set(name, sort),
	}
}
