package el

type Runtime struct {
	ParseLiteral func(s string) (Data, error)
}
