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

func newEmptyRuntime(InitialHeader Name, TerminalHeader Name) Runtime {
	r := Runtime{
		sortUniverse: sortUniverse{
			initialHeader:  InitialHeader,
			terminalHeader: TerminalHeader,
		},
		frame:  frame{},
		parser: parser{},
	}
	// parser depends on sortUniverse and frame
	r.parser.parseName = func(name Name) Sort {
		if sort, ok := r.frame.Get(name); ok {
			return sort
		}
		if sort, ok := r.sortUniverse.parseConstant(name); ok {
			return sort
		}
		panic("name not found")
	}
	return r
}

type Runtime struct {
	sortUniverse
	frame
	parser
}

func (u Runtime) NewTerm(name Name, parent Sort) (Runtime, Sort) {
	panic("unimplemented")
}
