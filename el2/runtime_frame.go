package el2

import "github.com/fbundle/sorts/persistent/ordered_map"

type runtimeFrame struct {
	frame ordered_map.OrderedMap[Name, Sort]
}

func (f runtimeFrame) Get(name Name) (Sort, bool) {
	return f.frame.Get(name)
}

func (f runtimeFrame) Set(name Name, sort Sort) runtimeFrame {
	return runtimeFrame{
		frame: f.frame.Set(name, sort),
	}
}
