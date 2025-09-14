package el2

import "github.com/fbundle/sorts/el2/el_almost_sort"

func DefaultRuntime() Parser {
	return Parser{
		sortUniverse: sortUniverse{
			initialHeader:  "Unit",
			terminalHeader: "Any",
		}.mustSortAttr(),
	}.
		NewListParser("->", toListParser(ListParseArrow("->"))).
		NewListParser("⊕", toListParser(ListParseSum("⊕"))).
		NewListParser("⊗", toListParser(ListParseProd("⊗"))).
		NewListParser("=>", el_almost_sort.ListParseLambda)
}
