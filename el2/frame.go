package el2

import "github.com/fbundle/sorts/persistent/ordered_map"

type frame struct {
	frame ordered_map.OrderedMap[Name, Sort]
}

func (f frame) Get(name Name) (Sort, bool) {
	return f.frame.Get(name)
}

func (f frame) Set(name Name, sort Sort) frame {
	return frame{
		frame: f.frame.Set(name, sort),
	}
}
