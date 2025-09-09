package main

import (
	"fmt"
	"log"

	"github.com/fbundle/sorts/el"
	"github.com/fbundle/sorts/el/ast"
	"github.com/fbundle/sorts/expr"
	sorts "github.com/fbundle/sorts/sorts/sorts_v3"
)

func main() {
	// Create a new interpreter
	interpreter := el.NewInterpreter()

	// Example 1: Simple arithmetic
	fmt.Println("=== Example 1: Simple Arithmetic ===")
	evalExample(interpreter, "(+ 2 3)")
	evalExample(interpreter, "(× 4 5)")
	evalExample(interpreter, "(+ (+ 1 2) (× 3 4))")

	// Example 2: Let bindings
	fmt.Println("\n=== Example 2: Let Bindings ===")
	evalExample(interpreter, "(let x 5 (+ x 2))")
	evalExample(interpreter, "(let x 3 (let y 4 (+ x y)))")

	// Example 3: Type constructors
	fmt.Println("\n=== Example 3: Type Constructors ===")
	evalExample(interpreter, "(Sum Nat Nat)")
	evalExample(interpreter, "(Prod Nat Nat)")

	// Example 4: Lambda expressions
	fmt.Println("\n=== Example 4: Lambda Expressions ===")
	evalExample(interpreter, "(x => (+ x 1))")

	// Example 5: Match expressions
	fmt.Println("\n=== Example 5: Match Expressions ===")
	evalExample(interpreter, "(match 1 1 2 0)") // if 1 matches 1, return 2, else 0

	// Example 6: Complex nested expression
	fmt.Println("\n=== Example 6: Complex Nested Expression ===")
	evalExample(interpreter, "(let f (x => (+ x 1)) (f 5))")
}

func evalExample(interpreter *el.Interpreter, input string) {
	fmt.Printf("Input: %s\n", input)

	// Parse the input
	tokens := expr.Tokenize(input)
	expr, _, err := expr.Parse(tokens)
	if err != nil {
		log.Printf("Failed to parse input: %v", err)
		return
	}

	// Convert to AST
	astExpr, err := ast.Parse(expr)
	if err != nil {
		log.Printf("Failed to convert to AST: %v", err)
		return
	}

	// Evaluate
	result, err := interpreter.Eval(astExpr)
	if err != nil {
		log.Printf("Failed to evaluate: %v", err)
		return
	}

	// Display result
	fmt.Printf("Result: %s\n", sorts.Name(result))
	fmt.Printf("Type: %s\n", sorts.Name(sorts.Parent(result)))
	fmt.Println()
}
