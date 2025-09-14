package el2

func DefaultRuntime() Runtime {
	return Runtime{
		sortUniverse: sortUniverse{
			initialHeader:  "Unit",
			terminalHeader: "Any",
		}.mustSortAttr(),
	}.
		NewListParser("->", toListParser(ListParseArrow("->"))).
		NewListParser("⊕", toListParser(ListParseSum("⊕"))).
		NewListParser("⊗", toListParser(ListParseProd("⊗"))).
		NewListParser("=>", ListParseLambda)
}
