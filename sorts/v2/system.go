package sorts

type SortSystemOption func(*sortSystem)

func WithInitialName(name string) SortSystemOption {
	return func(ss *sortSystem) {
		ss.initialName = name
	}
}

func WithTerminalName(name string) SortSystemOption {
	return func(ss *sortSystem) {
		ss.terminalName = name
	}
}

func NewSortSystem(defaultName string, opts ...SortSystemOption) SortSystem {
	ss := &sortSystem{
		initialName:  "",
		terminalName: "",
		defaultName:  defaultName,
		lessEqualMap: make(map[rule]struct{}),
	}
	for _, opt := range opts {
		opt(ss)
	}
	return ss.validate()
}

type rule struct {
	src string
	dst string
}
type sortSystem struct {
	initialName  string // empty = noset
	terminalName string // empty = noset
	defaultName  string // must be non-empty, used when no parent is set
	lessEqualMap map[rule]struct{}
}

func (ss *sortSystem) LessEqual(src string, dst string) bool {
	if ss.isInitial(src) {
		return true
	}
	if ss.isTerminal(dst) {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := ss.lessEqualMap[rule{src: src, dst: dst}]; ok {
		return true
	}
	return false
}

func (ss *sortSystem) Default(level int) Sort {
	return atom{
		level:  level,
		name:   ss.defaultName,
		parent: nil,
	}
}

func (ss *sortSystem) Atom(level int, name string, parents ...Sort) Sort {
	if len(parents) >= 2 {
		panic("type_error")
	}
	var parent Sort = nil
	if len(parents) >= 1 {
		parent = parents[0]
	}
	if parent != nil && parent.Level(ss) != level+1 {
		panic("type_error")
	}
	return atom{
		level:  level,
		name:   name,
		parent: parent,
	}
}

func (ss *sortSystem) Arrow(param Sort, body Sort) Sort {
	return arrow{
		param: param,
		body:  body,
	}
}

func (ss *sortSystem) AddRule(src string, dst string) SortSystem {
	ss.lessEqualMap[rule{src: src, dst: dst}] = struct{}{}
	return ss
}

func (ss *sortSystem) isInitial(name string) bool {
	return len(ss.initialName) > 0 && ss.initialName == name
}
func (ss *sortSystem) isTerminal(name string) bool {
	return len(ss.terminalName) > 0 && ss.terminalName == name
}

func (ss *sortSystem) validate() *sortSystem {
	if len(ss.defaultName) == 0 {
		panic("validate_error: default name must not be empty")
	}
	if len(ss.initialName) > 0 && len(ss.terminalName) > 0 && ss.initialName == ss.terminalName {
		panic("validate_error: initial and terminal must be distinct")
	}
	return ss
}
