package main

import (
	"fmt"
	"strings"
)

type Constructor struct {
	Name     string
	TypeList []string
}

func (c Constructor) String() string {
	return fmt.Sprintf("\t| %s : %s", c.Name, strings.Join(c.TypeList, " -> "))
}

type Inductive struct {
	Name         string
	Constructors []Constructor
}

func (i Inductive) String() string {
	lines := make([]string, 0, len(i.Constructors)+1)
	lines = append(lines, fmt.Sprintf("inductive %s", i.Name))
	for _, constructor := range i.Constructors {
		lines = append(lines, constructor.String())
	}
	return strings.Join(lines, "\n")
}

func main() {
	i := Inductive{
		Name: "Nat",
		Constructors: []Constructor{
			{Name: "zero", TypeList: []string{"Nat"}},
			{Name: "succ", TypeList: []string{"Nat", "Nat"}},
		},
	}
	fmt.Println(i.String())
}
