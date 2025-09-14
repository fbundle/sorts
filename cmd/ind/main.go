package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Constructor struct {
	Name     string
	TypeList []string
}

func (c Constructor) String() string {
	return fmt.Sprintf("\t| %s\t: %s", c.Name, strings.Join(c.TypeList, " -> "))
}

type Inductive struct {
	Name         string
	Constructors []Constructor
}

func (ind Inductive) String() string {
	lines := make([]string, 0, len(ind.Constructors)+1)
	lines = append(lines, fmt.Sprintf("inductive %s", ind.Name))
	for _, c := range ind.Constructors {
		lines = append(lines, c.String())
	}
	return strings.Join(lines, "\n")
}

func (ind Inductive) MustOk() Inductive {
	if len(ind.Name) == 0 {
		panic("inductive must have a name")
	}
	if !unicode.IsUpper([]rune(ind.Name)[0]) {
		panic("inductive name must be public")
	}
	for _, c := range ind.Constructors {
		if !unicode.IsUpper([]rune(c.Name)[0]) {
			panic("constructor name must be public")
		}
		if c.TypeList[len(c.TypeList)-1] != ind.Name {
			panic("constructor type must return inductive type")
		}
	}
	return ind
}

func (ind Inductive) Generate() string {
	ind = ind.MustOk()

	lines := make([]string, 0)
	push := func(format string, args ...interface{}) {
		lines = append(lines, fmt.Sprintf(format, args...))
	}

	push("package %s", strings.ToLower(ind.Name))

	push("/*")
	push("the code below was auto generated for inductive type")
	push(ind.String())
	push("*/")

	push("type %s interface {", ind.Name)
	push("\tattr%s()", ind.Name)
	push("}")

	for _, c := range ind.Constructors {
		argTypeList := c.TypeList[:len(c.TypeList)-1]

		push("type %s struct {", c.Name)
		for i, t := range argTypeList {
			push("\tField_%d %s", i, t)
		}
		push("}")

		push("func (o %s) attr%s() {}", c.Name, ind.Name)

		push("func (o %s) Unwrap() (%s) {", c.Name, strings.Join(argTypeList, ","))
		vals := make([]string, 0, len(argTypeList))
		for i, _ := range argTypeList {
			vals = append(vals, fmt.Sprintf("o.Field_%d", i))
		}
		push("\treturn %s", strings.Join(vals, " , "))
		push("}")

	}

	push("type Match[T any] struct {")
	for _, c := range ind.Constructors {
		push("\tMap%s func(%s) T", c.Name, c.Name)
	}
	push("}")

	push("func (m Match[T]) Apply(o %s) T {", ind.Name)

	push("\tswitch o := o.(type) {")
	for _, c := range ind.Constructors {
		push("\t\tcase %s:", c.Name)
		push("\t\t\treturn m.Map%s(o)", c.Name)
	}
	push("\t\tdefault:")
	push("\t\t\tpanic(\"unreachable\")")
	push("\t}")
	push("}")

	return strings.Join(lines, "\n")
}

func split(s string, sep string) []string {
	trimSpaces := func(parts []string) []string {
		newParts := make([]string, 0, len(parts))
		for _, part := range parts {
			newParts = append(newParts, strings.TrimSpace(part))
		}
		return newParts
	}
	return trimSpaces(strings.Split(s, sep))
}

func Parse(s string) Inductive {
	ind := Inductive{}

	parts := split(s, "|")

	part0 := parts[0]
	fields := strings.Fields(part0)
	if len(fields) < 2 || fields[0] != "inductive" {
		panic("inductive must start with \"inductive TypeName ...\"")
	}
	ind.Name = fields[1]

	for _, part := range parts[1:] {
		subParts := split(part, ":")
		if len(subParts) != 2 {
			panic("constructor must be in the form \"ConstructorName : ConstructorType\"")
		}
		c := Constructor{}
		c.Name = strings.TrimSpace(subParts[0])
		c.TypeList = split(subParts[1], "->")
		ind.Constructors = append(ind.Constructors, c)
	}
	return ind.MustOk()
}

func main() {
	filename := os.Args[1]

	if !strings.HasSuffix(filename, ".ind") {
		panic("induction filename must end with .ind")
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ind := Parse(string(b))

	outputFilename := strings.TrimSuffix(filename, ".ind") + ".go"
	err = os.WriteFile(outputFilename, []byte(ind.Generate()), 0644)
	if err != nil {
		panic(err)
	}
}
