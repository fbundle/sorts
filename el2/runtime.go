package el2

func DefaultRuntime() Runtime {
	return Runtime{
		runtimeSortAttr: runtimeSortAttr{
			initialHeader:  "Unit",
			terminalHeader: "Any",
		}.mustSortAttr(),
	}.
		NewListParser("->", toListParser(ListParseArrow("->"))).
		NewListParser("⊕", toListParser(ListParseSum("⊕"))).
		NewListParser("⊗", toListParser(ListParseProd("⊗"))).
		NewListParser("=>", ListParseLambda)
}

type Runtime struct {
	runtimeSortAttr
	runtimeFrame
	runtimeParser
}

func (u Runtime) NewTerm(name Name, parent Sort) (Runtime, Sort) {
	panic("unimplemented")
}
