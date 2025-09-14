package el2

import (
	"github.com/fbundle/sorts/el2/almost_sort"
	"github.com/fbundle/sorts/el2/almost_sort_extra"
	"github.com/fbundle/sorts/form"
)

func (r Context) Get(name form.Name) almost_sort.ActualSort {
	s, ok := r.frame.Get(name)
	if !ok {
		panic(TypeErr)
	}
	return s
}

func (r Context) Set(name form.Name, sort almost_sort.ActualSort) almost_sort_extra.Context {
	r.frame = r.frame.Set(name, sort)
	return r
}
