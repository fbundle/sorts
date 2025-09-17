package sorts5

import "github.com/fbundle/sorts/form"

type ListParseFunc = func(form.List) Sort

var listParseFuncMap = map[Name]ListParseFunc{}

func AddListParseFunc(name Name, fn ListParseFunc) {
	listParseFuncMap[name] = fn
}
