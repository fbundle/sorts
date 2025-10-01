package sorts_parser

import "github.com/fbundle/sorts/sorts"

func compileBinding(parse func(form Form) Code, list List) sorts.Binding {
	err := compileErr(list, []string{string(sorts.BindingCmd), "name", "value"})
	if len(list) != 2 {
		panic(err)
	}
	return sorts.Binding{
		Name:  mustType[Name](err, list[0]),
		Value: parse(list[1]),
	}
}
func compileAnnot(parse func(form Form) Code, list List) sorts.Annot {
	err := compileErr(list, []string{string(sorts.AnnotCmd), "name", "type"})
	if len(list) != 2 {
		panic(err)
	}
	return sorts.Annot{
		Name: mustType[Name](err, list[0]),
		Type: parse(list[1]),
	}
}
