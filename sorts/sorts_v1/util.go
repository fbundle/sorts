package sorts

const (
	// Unit - Unit can be cast into any type (initial object), it has a unique zero value
	// in some category, the initial object does not have any element - maybe I will resolve it later
	Unit = "unit"
	// Any - every type can be cast into Any (terminal object)
	Any = "any"
)

type rule struct {
	src string
	dst string
}

var leMap = make(map[rule]struct{})

func AddRule(src string, dst string) {
	leMap[rule{src, dst}] = struct{}{}
}

func le(src string, dst string) bool {
	if src == Unit || dst == Any {
		return true
	}
	if src == dst {
		return true
	}
	if _, ok := leMap[rule{src, dst}]; ok {
		return true
	}
	return false
}
