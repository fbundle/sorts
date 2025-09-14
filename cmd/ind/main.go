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
	header := fmt.Sprintf("inductive %s", i.Name)
	lines := make([]string, len(i.Constructors))
	for _, constructor := range i.Constructors {
		lines = append(lines, constructor.String())
	}
	return header + strings.Join(lines, "\n")
}

func main() {
	i := Inductive{
		Name:         "Nat",
		Constructors: []Constructor{},
	}
	fmt.Println(i.String())
}
