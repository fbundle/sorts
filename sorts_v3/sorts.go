package sorts

const (
	defaultName  = "type"
	InitialName  = "unit" // initial object
	TerminalName = "any"  // terminal object
)

func Name(s Sort) string {
	return s.attr().name
}

func Level(s Sort) int {
	return s.attr().level
}

func Parent(s Sort) Sort {
	return s.attr().parent
}

func LessEqual(x Sort, y Sort) bool {
	return x.attr().lessEqual(y)
}

type Sort interface {
	attr() sortAttr
}

type sortAttr struct {
	level     int                 // universe Level
	name      string              // every Sort is identified with a Name (string)
	parent    Sort                // (or Type) every Sort must have a Parent
	lessEqual func(dst Sort) bool // a partial order on sorts (subtype)
}
