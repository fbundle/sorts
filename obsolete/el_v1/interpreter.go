package el

import (
	"fmt"

	"github.com/fbundle/sorts/el_v2/ast"
	sorts2 "github.com/fbundle/sorts/obsolete/sorts_v3"
	"github.com/fbundle/sorts/persistent/ordered_map"
)

// Value represents a value in our interpreter, which is essentially a sort
type Value = sorts2.Sort

// FrameEntry represents an entry in the frame (environment)
type FrameEntry struct {
	Value Value
	Type  Value
}

// Frame represents the environment/context for evaluation
type Frame = ordered_map.OrderedMap[string, FrameEntry]

// Interpreter represents the main interpreter
type Interpreter struct {
	frame Frame
}

// NewInterpreter creates a new interpreter with built-in functions
func NewInterpreter() *Interpreter {
	frame := ordered_map.EmptyOrderedMap[string, FrameEntry]()

	// Add built-in natural numbers
	frame = frame.Set("0", FrameEntry{Value: sorts2.Zero, Type: sorts2.Nat})
	frame = frame.Set("1", FrameEntry{Value: sorts2.NewAtom(0, "1", sorts2.Nat), Type: sorts2.Nat})
	frame = frame.Set("2", FrameEntry{Value: sorts2.NewAtom(0, "2", sorts2.Nat), Type: sorts2.Nat})
	frame = frame.Set("3", FrameEntry{Value: sorts2.NewAtom(0, "3", sorts2.Nat), Type: sorts2.Nat})
	frame = frame.Set("4", FrameEntry{Value: sorts2.NewAtom(0, "4", sorts2.Nat), Type: sorts2.Nat})
	frame = frame.Set("5", FrameEntry{Value: sorts2.NewAtom(0, "5", sorts2.Nat), Type: sorts2.Nat})

	// Add built-in types
	frame = frame.Set("Nat", FrameEntry{Value: sorts2.Nat, Type: sorts2.NewAtom(1, "Type", nil)})
	frame = frame.Set("Type", FrameEntry{Value: sorts2.NewAtom(1, "Type", nil), Type: sorts2.NewAtom(2, "Type1", nil)})

	// Add built-in arithmetic operations
	frame = frame.Set("+", FrameEntry{
		Value: sorts2.NatToNatToNat.Intro("add", func(a sorts2.Sort) sorts2.Sort {
			sorts2.MustTermOf(a, sorts2.Nat)
			return sorts2.NatToNat.Intro("add_a", func(b sorts2.Sort) sorts2.Sort {
				sorts2.MustTermOf(b, sorts2.Nat)
				return addNats(a, b)
			})
		}),
		Type: sorts2.NatToNatToNat,
	})

	frame = frame.Set("×", FrameEntry{
		Value: sorts2.NatToNatToNat.Intro("mul", func(a sorts2.Sort) sorts2.Sort {
			sorts2.MustTermOf(a, sorts2.Nat)
			return sorts2.NatToNat.Intro("mul_a", func(b sorts2.Sort) sorts2.Sort {
				sorts2.MustTermOf(b, sorts2.Nat)
				return mulNats(a, b)
			})
		}),
		Type: sorts2.NatToNatToNat,
	})

	// Add Sum and Prod type constructors for any universe level
	frame = frame.Set("Sum", FrameEntry{
		Value: sorts2.NewAtom(0, "Sum", sorts2.NewAtom(1, "TypeConstructor", nil)),
		Type:  sorts2.NewAtom(1, "TypeConstructor", nil),
	})

	frame = frame.Set("Prod", FrameEntry{
		Value: sorts2.NewAtom(0, "Prod", sorts2.NewAtom(1, "TypeConstructor", nil)),
		Type:  sorts2.NewAtom(1, "TypeConstructor", nil),
	})

	// Add Pi and Sigma type constructors
	frame = frame.Set("Pi", FrameEntry{
		Value: sorts2.NewAtom(0, "Pi", sorts2.NewAtom(1, "TypeConstructor", nil)),
		Type:  sorts2.NewAtom(1, "TypeConstructor", nil),
	})

	frame = frame.Set("Sigma", FrameEntry{
		Value: sorts2.NewAtom(0, "Sigma", sorts2.NewAtom(1, "TypeConstructor", nil)),
		Type:  sorts2.NewAtom(1, "TypeConstructor", nil),
	})

	return &Interpreter{frame: frame}
}

// Eval evaluates an AST expression and returns a value
func (i *Interpreter) Eval(expr ast.Expr) (Value, error) {
	switch e := expr.(type) {
	case ast.Name:
		return i.evalName(e)
	case ast.Let:
		return i.evalLet(e)
	case ast.Match:
		return i.evalMatch(e)
	case ast.Lambda:
		return i.evalLambda(e)
	case ast.FunctionCall:
		return i.evalFunctionCall(e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

// evalName evaluates a name (variable lookup)
func (i *Interpreter) evalName(name ast.Name) (Value, error) {
	entry, ok := i.frame.Get(string(name))
	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", name)
	}
	return entry.Value, nil
}

// evalLet evaluates a let binding
func (i *Interpreter) evalLet(let ast.Let) (Value, error) {
	// Create a new frame with the bindings
	newFrame := i.frame

	for _, binding := range let.Bindings {
		value, err := i.Eval(binding.Expr)
		if err != nil {
			return nil, fmt.Errorf("error evaluating binding %s: %w", binding.Name, err)
		}

		// For now, we'll infer the type (in a full implementation, we'd have type annotations)
		entry := FrameEntry{
			Value: value,
			Type:  i.inferType(value),
		}
		newFrame = newFrame.Set(string(binding.Name), entry)
	}

	// Create a new interpreter with the extended frame
	newInterpreter := &Interpreter{frame: newFrame}
	return newInterpreter.Eval(let.Final)
}

// evalMatch evaluates a match expression
func (i *Interpreter) evalMatch(match ast.Match) (Value, error) {
	// Evaluate the condition
	condValue, err := i.Eval(match.Cond)
	if err != nil {
		return nil, fmt.Errorf("error evaluating match condition: %w", err)
	}

	// Try to match against each case
	for _, case_ := range match.Cases {
		compValue, err := i.Eval(case_.Comp)
		if err != nil {
			return nil, fmt.Errorf("error evaluating match case: %w", err)
		}

		// Check if the condition matches this case using sorts equality
		if i.valuesEqual(condValue, compValue) {
			return i.Eval(case_.Value)
		}
	}

	// If no case matches, evaluate the default
	if match.Default != nil {
		return i.Eval(match.Default)
	}

	return nil, fmt.Errorf("no matching case found in match expression")
}

// evalLambda evaluates a lambda expression
func (i *Interpreter) evalLambda(lambda ast.Lambda) (Value, error) {
	// For now, create a simple arrow type
	// In a full implementation, we'd need proper type inference
	return sorts2.Arrow{
		A: sorts2.NewAtom(0, "param", sorts2.Nat),  // Assume parameter is Nat for now
		B: sorts2.NewAtom(0, "result", sorts2.Nat), // Assume result is Nat for now
	}, nil
}

// evalFunctionCall evaluates a function call
func (i *Interpreter) evalFunctionCall(call ast.FunctionCall) (Value, error) {
	// Look up the function
	entry, ok := i.frame.Get(string(call.Cmd))
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", call.Cmd)
	}

	// Evaluate all arguments
	args := make([]Value, len(call.Args))
	for j, arg := range call.Args {
		value, err := i.Eval(arg)
		if err != nil {
			return nil, fmt.Errorf("error evaluating argument %d: %w", j, err)
		}
		args[j] = value
	}

	// Apply the function
	return i.applyFunction(entry.Value, args)
}

// applyFunction applies a function to arguments
func (i *Interpreter) applyFunction(fn Value, args []Value) (Value, error) {
	if len(args) == 0 {
		return fn, nil
	}

	// Handle arithmetic operations specially
	if len(args) == 2 {
		if i.isArithmeticOperation(fn) {
			return i.applyArithmeticOperation(fn, args[0], args[1])
		}
	}

	// Handle type constructors
	if i.isTypeConstructor(fn) {
		return i.applyTypeConstructor(fn, args)
	}

	// For now, handle simple cases
	switch f := fn.(type) {
	case sorts2.Arrow:
		if len(args) != 1 {
			return nil, fmt.Errorf("arrow function expects 1 argument, got %d", len(args))
		}
		return f.Elim(fn, args[0]), nil
	case sorts2.Pi:
		if len(args) != 1 {
			return nil, fmt.Errorf("pi function expects 1 argument, got %d", len(args))
		}
		return f.Elim(fn, args[0]), nil
	default:
		// Try to apply as a function created by Intro
		// This is a simplified approach - in practice, we'd need more sophisticated handling
		return i.applyIntroFunction(fn, args)
	}
}

// applyIntroFunction applies a function created by Intro
func (i *Interpreter) applyIntroFunction(fn Value, args []Value) (Value, error) {
	// This is a simplified implementation
	// In practice, we'd need to track the function definition and apply it properly
	return fn, nil
}

// isArithmeticOperation checks if a function is an arithmetic operation
func (i *Interpreter) isArithmeticOperation(fn Value) bool {
	name := sorts2.Name(fn)
	return name == "add" || name == "mul" || name == "+" || name == "×"
}

// applyArithmeticOperation applies an arithmetic operation
func (i *Interpreter) applyArithmeticOperation(fn Value, a, b Value) (Value, error) {
	name := sorts2.Name(fn)

	// Ensure both arguments are natural numbers
	if !i.isNaturalNumber(a) || !i.isNaturalNumber(b) {
		return nil, fmt.Errorf("arithmetic operations require natural numbers")
	}

	switch name {
	case "add", "+":
		return addNats(a, b), nil
	case "mul", "×":
		return mulNats(a, b), nil
	default:
		return nil, fmt.Errorf("unknown arithmetic operation: %s", name)
	}
}

// isTypeConstructor checks if a function is a type constructor
func (i *Interpreter) isTypeConstructor(fn Value) bool {
	name := sorts2.Name(fn)
	return name == "Sum" || name == "Prod" || name == "Pi" || name == "Sigma"
}

// applyTypeConstructor applies a type constructor
func (i *Interpreter) applyTypeConstructor(fn Value, args []Value) (Value, error) {
	name := sorts2.Name(fn)

	if len(args) != 2 {
		return nil, fmt.Errorf("type constructor %s expects 2 arguments, got %d", name, len(args))
	}

	a, b := args[0], args[1]

	switch name {
	case "Sum":
		return sorts2.Sum{A: a, B: b}, nil
	case "Prod":
		return sorts2.Prod{A: a, B: b}, nil
	case "Pi":
		// For Pi types, we need to create a dependent type
		return sorts2.Pi{
			A: a,
			B: sorts2.Dependent{
				Name: "B",
				Apply: func(x sorts2.Sort) sorts2.Sort {
					return b // In a full implementation, we'd substitute x in b
				},
			},
		}, nil
	case "Sigma":
		// For Sigma types, we need to create a dependent type
		return sorts2.Sigma{
			A: a,
			B: sorts2.Dependent{
				Name: "B",
				Apply: func(x sorts2.Sort) sorts2.Sort {
					return b // In a full implementation, we'd substitute x in b
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown type constructor: %s", name)
	}
}

// isNaturalNumber checks if a value is a natural number
func (i *Interpreter) isNaturalNumber(value Value) bool {
	// Check if the value is of type Nat
	parent := sorts2.Parent(value)
	return sorts2.Name(parent) == "Nat"
}

// valuesEqual checks if two values are equal using sorts equality
func (i *Interpreter) valuesEqual(v1, v2 Value) bool {
	// For now, use a simple name-based comparison
	// In a full implementation, we'd use the sorts equality system
	return sorts2.Name(v1) == sorts2.Name(v2)
}

// inferType infers the type of a value
func (i *Interpreter) inferType(value Value) Value {
	// This is a simplified type inference
	// In practice, we'd need more sophisticated type inference
	return sorts2.Parent(value)
}

// Helper functions for arithmetic operations

// addNats adds two natural numbers
func addNats(a, b sorts2.Sort) sorts2.Sort {
	// For now, create a symbolic addition
	// In a full implementation, we'd need proper natural number arithmetic
	return sorts2.NewAtom(0, fmt.Sprintf("(%s + %s)", sorts2.Name(a), sorts2.Name(b)), sorts2.Nat)
}

// mulNats multiplies two natural numbers
func mulNats(a, b sorts2.Sort) sorts2.Sort {
	// For now, create a symbolic multiplication
	// In a full implementation, we'd need proper natural number arithmetic
	return sorts2.NewAtom(0, fmt.Sprintf("(%s × %s)", sorts2.Name(a), sorts2.Name(b)), sorts2.Nat)
}
