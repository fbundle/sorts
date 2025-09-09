package el

import (
	"testing"

	"github.com/fbundle/sorts/el/ast"
	"github.com/fbundle/sorts/expr"
	sorts "github.com/fbundle/sorts/sorts/sorts_v3"
)

func TestInterpreter(t *testing.T) {
	interpreter := NewInterpreter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple number",
			input:    "0",
			expected: "0",
		},
		{
			name:     "addition",
			input:    "(+ 1 2)",
			expected: "(1 + 2)",
		},
		{
			name:     "multiplication",
			input:    "(× 2 3)",
			expected: "(2 × 3)",
		},
		{
			name:     "nested arithmetic",
			input:    "(+ (+ 1 2) (× 2 3))",
			expected: "((1 + 2) + (2 × 3))",
		},
		{
			name:     "let binding",
			input:    "(let x 5 (+ x 3))",
			expected: "(5 + 3)",
		},
		{
			name:     "multiple let bindings",
			input:    "(let x 2 (let y 3 (+ x y)))",
			expected: "(2 + 3)",
		},
		{
			name:     "type constructor Sum",
			input:    "(Sum Nat Nat)",
			expected: "Nat + Nat",
		},
		{
			name:     "type constructor Prod",
			input:    "(Prod Nat Nat)",
			expected: "Nat × Nat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input
			tokens := expr.Tokenize(tt.input)
			expr, _, err := expr.Parse(tokens)
			if err != nil {
				t.Fatalf("Failed to parse input: %v", err)
			}

			// Convert to AST
			astExpr, err := ast.Parse(expr)
			if err != nil {
				t.Fatalf("Failed to convert to AST: %v", err)
			}

			// Evaluate
			result, err := interpreter.Eval(astExpr)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			// Check result
			resultName := sorts.Name(result)
			if resultName != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, resultName)
			}
		})
	}
}

func TestArithmeticOperations(t *testing.T) {
	interpreter := NewInterpreter()

	// Test addition
	result, err := interpreter.Eval(ast.Name("0"))
	if err != nil {
		t.Fatalf("Failed to evaluate 0: %v", err)
	}
	if sorts.Name(result) != "0" {
		t.Errorf("Expected 0, got %s", sorts.Name(result))
	}

	// Test addition function call
	addCall := ast.FunctionCall{
		Cmd:  "+",
		Args: []ast.Expr{ast.Name("1"), ast.Name("2")},
	}
	result, err = interpreter.Eval(addCall)
	if err != nil {
		t.Fatalf("Failed to evaluate addition: %v", err)
	}
	if sorts.Name(result) != "(1 + 2)" {
		t.Errorf("Expected (1 + 2), got %s", sorts.Name(result))
	}
}

func TestTypeConstructors(t *testing.T) {
	interpreter := NewInterpreter()

	// Test Sum type constructor
	sumCall := ast.FunctionCall{
		Cmd:  "Sum",
		Args: []ast.Expr{ast.Name("Nat"), ast.Name("Nat")},
	}
	result, err := interpreter.Eval(sumCall)
	if err != nil {
		t.Fatalf("Failed to evaluate Sum constructor: %v", err)
	}

	// Check that it's a Sum type
	if _, ok := result.(sorts.Sum); !ok {
		t.Errorf("Expected Sum type, got %T", result)
	}

	// Test Prod type constructor
	prodCall := ast.FunctionCall{
		Cmd:  "Prod",
		Args: []ast.Expr{ast.Name("Nat"), ast.Name("Nat")},
	}
	result, err = interpreter.Eval(prodCall)
	if err != nil {
		t.Fatalf("Failed to evaluate Prod constructor: %v", err)
	}

	// Check that it's a Prod type
	if _, ok := result.(sorts.Prod); !ok {
		t.Errorf("Expected Prod type, got %T", result)
	}
}

func TestLetBinding(t *testing.T) {
	interpreter := NewInterpreter()

	// Test simple let binding
	letExpr := ast.Let{
		Bindings: []ast.LetBinding{
			{Name: "x", Expr: ast.Name("5")},
		},
		Final: ast.FunctionCall{
			Cmd:  "+",
			Args: []ast.Expr{ast.Name("x"), ast.Name("3")},
		},
	}

	result, err := interpreter.Eval(letExpr)
	if err != nil {
		t.Fatalf("Failed to evaluate let expression: %v", err)
	}

	if sorts.Name(result) != "(5 + 3)" {
		t.Errorf("Expected (5 + 3), got %s", sorts.Name(result))
	}
}

func TestLambda(t *testing.T) {
	interpreter := NewInterpreter()

	// Test lambda expression
	lambdaExpr := ast.Lambda{
		Param: "x",
		Body: ast.FunctionCall{
			Cmd:  "+",
			Args: []ast.Expr{ast.Name("x"), ast.Name("1")},
		},
	}

	result, err := interpreter.Eval(lambdaExpr)
	if err != nil {
		t.Fatalf("Failed to evaluate lambda: %v", err)
	}

	// Check that it's an Arrow type
	if _, ok := result.(sorts.Arrow); !ok {
		t.Errorf("Expected Arrow type, got %T", result)
	}
}

func TestMatch(t *testing.T) {
	interpreter := NewInterpreter()

	// Test match expression
	matchExpr := ast.Match{
		Cond: ast.Name("1"),
		Cases: []ast.MatchCase{
			{
				Comp:  ast.Name("1"),
				Value: ast.Name("2"),
			},
		},
		Default: ast.Name("0"),
	}

	result, err := interpreter.Eval(matchExpr)
	if err != nil {
		t.Fatalf("Failed to evaluate match: %v", err)
	}

	if sorts.Name(result) != "2" {
		t.Errorf("Expected 2, got %s", sorts.Name(result))
	}
}

// Benchmark tests
func BenchmarkArithmetic(b *testing.B) {
	interpreter := NewInterpreter()

	addCall := ast.FunctionCall{
		Cmd:  "+",
		Args: []ast.Expr{ast.Name("1"), ast.Name("2")},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interpreter.Eval(addCall)
		if err != nil {
			b.Fatalf("Failed to evaluate: %v", err)
		}
	}
}

func BenchmarkTypeConstructor(b *testing.B) {
	interpreter := NewInterpreter()

	sumCall := ast.FunctionCall{
		Cmd:  "Sum",
		Args: []ast.Expr{ast.Name("Nat"), ast.Name("Nat")},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interpreter.Eval(sumCall)
		if err != nil {
			b.Fatalf("Failed to evaluate: %v", err)
		}
	}
}
