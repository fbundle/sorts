# EL Interpreter

A Lean4-inspired interpreter for the EL language that supports dependent types using the sorts system.

## Features

- **Natural Numbers**: Built-in support for natural numbers at universe level 0
- **Arithmetic Operations**: Addition (`+`) and multiplication (`×`) operations
- **Dependent Types**: Full support for Sum, Prod, Pi, and Sigma types
- **Let Bindings**: Local variable bindings with type inference
- **Lambda Expressions**: Function definitions (simplified implementation)
- **Match Expressions**: Pattern matching on values
- **Type System**: Universe levels with types at level 1 and above

## Architecture

The interpreter is built on top of the `sorts_v3` system and uses:

- **Frame Management**: `ordered_map.OrderedMap` for environment management
- **Value Representation**: Values are represented as sorts from the sorts system
- **Type Checking**: Integrated with the dependent type system

## Usage

```go
// Create a new interpreter
interpreter := el.NewInterpreter()

// Parse and evaluate an expression
tokens := expr.Tokenize("(+ 2 3)")
expr, _, _ := expr.Parse(tokens)
astExpr, _ := ast.Parse(expr)
result, _ := interpreter.Eval(astExpr)
```

## Supported Expressions

### Arithmetic
- `(+ 1 2)` - Addition
- `(× 3 4)` - Multiplication
- `(+ (+ 1 2) (× 3 4))` - Nested operations

### Let Bindings
- `(let x 5 (+ x 2))` - Simple binding
- `(let x 3 (let y 4 (+ x y)))` - Nested bindings

### Type Constructors
- `(Sum Nat Nat)` - Sum type
- `(Prod Nat Nat)` - Product type
- `(Pi A B)` - Pi type (dependent function)
- `(Sigma A B)` - Sigma type (dependent pair)

### Match Expressions
- `(match 1 1 2 0)` - Pattern matching

### Lambda Expressions
- `(x => (+ x 1))` - Function definition (simplified)

## Type System

The interpreter supports a universe hierarchy:

- **Level 0**: Natural numbers (`Nat`)
- **Level 1**: Types (`Type`)
- **Level 2+**: Higher-order types

### Built-in Types

- `Nat`: Natural numbers
- `Type`: Type universe
- `Sum A B`: Sum type (A + B)
- `Prod A B`: Product type (A × B)
- `Pi A B`: Dependent function type (Π(x:A)B(x))
- `Sigma A B`: Dependent pair type (Σ(x:A)B(x))

## Implementation Details

### Frame Management

The interpreter uses a persistent ordered map to manage the environment:

```go
type Frame = ordered_map.OrderedMap[string, FrameEntry]

type FrameEntry struct {
    Value Value
    Type  Value
}
```

### Value Representation

Values are represented as sorts from the sorts system:

```go
type Value = sorts.Sort
```

### Arithmetic Operations

Arithmetic operations are implemented symbolically:

- `(+ a b)` creates `(a + b)`
- `(× a b)` creates `(a × b)`

## Testing

Run the tests with:

```bash
go test ./el_v2/... -v
```

## Examples

See `examples/interpreter_demo.go` for a complete demonstration of the interpreter's capabilities.

## Limitations

- Lambda expressions have simplified implementation
- No proper type inference for lambda parameters
- Arithmetic operations are symbolic (no evaluation to concrete numbers)
- Limited pattern matching capabilities

## Future Improvements

- Full lambda evaluation with proper variable binding
- Concrete arithmetic evaluation
- Enhanced type inference
- More sophisticated pattern matching
- Error handling improvements
