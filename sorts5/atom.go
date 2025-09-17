package sorts5

const (
	initialName  = "Unit"
	terminalName = "Any"
)

type Atom struct {
	form   Form
	level  func() int
	parent func() Sort
}

func (s Atom) Form() Form {
	return s.form
}

func (s Atom) Level() int {
	return s.level()
}

func (s Atom) Parent() Sort {
	return s.parent()
}

func (s Atom) LessEqual(dst Sort) bool {
	// compare form
	sName, ok1 := isName(s.Form())
	dName, ok2 := isName(dst.Form())
	if ok1 && sName == initialName {
		return true
	}
	if ok2 && dName == terminalName {
		return true
	}
	if ok1 && ok2 {
		if _, ok3 := ruleMap[[2]Name{sName, dName}]; ok3 {
			return true
		}
	}
	return false
}

var ruleMap = map[[2]Name]struct{}{}

func AddRule(src Name, dst Name) {
	ruleMap[[2]Name{src, dst}] = struct{}{}
}
