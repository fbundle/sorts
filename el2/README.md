# EL2 - Enhanced Expression Language

EL2 is an enhanced version of the EL expression language, built on top of the `sorts` dependent type system. It provides a more robust and feature-rich environment for working with dependent types, pattern matching, and functional programming constructs.

## Overview

EL2 extends the original EL language with improved type checking, better error handling, and enhanced compilation capabilities. It serves as a more mature implementation of the expression language, designed to support complex mathematical reasoning and proof development.

## Architecture

### Core Components

- **`context.go`**: Main context implementation that orchestrates compilation, frame management, and universe operations
- **`context_compiler.go`**: Compilation logic for converting forms into sorts
- **`context_frame.go`**: Frame management for name binding and variable scoping
- **`context_universe.go`**: Universe operations and type hierarchy management
- **`el_sorts/`**: Extended sort definitions and compilation functions
- **`universe/`**: Universe implementation with built-in type support

### Key Features

#### 1. Enhanced Type System
- **Universe levels**: Support for `Unit_n` and `Any_n` at different levels
- **Type constructors**: Arrow (`->`), Sum (`⊕`), Product (`⊗`) types
- **Dependent types**: Full support for dependent function and pair types
- **Type checking**: Comprehensive static type checking with detailed error reporting

#### 2. Advanced Language Constructs

**Lambda Abstraction**:
```el
(lambda param1 type1 ... paramN typeN body)
```
- Multi-parameter lambda functions
- Type annotations for parameters
- Automatic arrow type construction

**Beta Reduction**:
```el
(cmd arg)
```
- Function application with type checking
- Validates argument types against function domain

**Let Bindings**:
```el
(let name1 value1 ... nameN valueN final)
```
- Sequential name binding
- Support for inhabitant renaming
- Lexical scoping

**Pattern Matching**:
```el
(match cond pattern1 value1 ... patternN valueN final)
```
- Exact pattern matching with `(exact expr)`
- Type-safe pattern binding
- Exhaustiveness checking

**Inhabitants**:
```el
(inh type)
```
- Generate undefined terms of a given type
- Automatic name generation
- Useful for proof construction

**Type Annotation**:
```el
(type value)
```
- Explicit type annotations
- Type validation

#### 3. Compilation System

The compilation system in EL2 is built around a context-based approach:

- **Context Interface**: Unified interface for compilation, frame management, and universe operations
- **List Compilers**: Extensible system for handling different language constructs
- **Type Inference**: Automatic type inference with manual override capabilities
- **Error Handling**: Comprehensive error reporting with context information

#### 4. Universe Management

EL2 includes a sophisticated universe system:

- **Level Management**: Support for universe levels with proper hierarchy
- **Built-in Types**: `Unit` and `Any` types at various levels
- **Type Ordering**: Subtyping relationships and type ordering
- **Rule System**: Custom type ordering rules

## Usage

### Command Line Interface

EL2 can be used through the command-line interface:

```bash
# Build the el2 command
go build -o ./bin/el2 ./cmd/el2

# Run el2 on a file
./bin/el2 example.el
```

### Example Programs

**Basic Lambda and Application**:
```el
(lambda x Any_0 (x x))
```

**Let Bindings with Types**:
```el
(let
    id {Any_0 -> Any_0} (lambda x Any_0 x)
    (id id)
)
```

**Pattern Matching**:
```el
(let
    Bool U_1 undef
    True Bool undef
    False Bool undef
    
    is_zero {Nat -> Bool} (lambda n Nat (match n
        (exact n0) True
        False
    ))
    
    (is_zero n0)
)
```

**Sum and Product Types**:
```el
(let
    x {Nat ⊕ Bool} (inh {Nat ⊕ Bool})
    y {Nat ⊗ Bool} (inh {Nat ⊗ Bool})
    x
)
```

## Implementation Details

### Context System

The `Context` struct implements the `el_sorts.Context` interface and provides:

- **Frame Management**: Persistent ordered map for variable bindings
- **Universe Integration**: Access to type universe operations
- **Compiler Registry**: Extensible list compiler system
- **Type Operations**: Parent type resolution, level checking, subtyping

### Compilation Process

1. **Tokenization**: Input is tokenized using the form processor
2. **Parsing**: Forms are parsed into abstract syntax trees
3. **Compilation**: Forms are compiled into sorts using the context
4. **Type Checking**: Types are validated and inferred
5. **Output**: Results are formatted and displayed

### Error Handling

EL2 provides detailed error reporting:

- **Type Errors**: Clear messages for type mismatches
- **Parse Errors**: Syntax error reporting
- **Name Resolution**: Undefined variable errors
- **Arity Errors**: Incorrect number of arguments

## Extensibility

EL2 is designed to be extensible:

- **Custom Compilers**: Add new language constructs by implementing `ListCompileFunc`
- **Type Rules**: Extend the universe with custom type ordering rules
- **Built-in Types**: Add new built-in types to the universe
- **Pattern Matching**: Extend pattern matching capabilities

## Dependencies

EL2 depends on several packages:

- **`sorts`**: Core dependent type system
- **`form`**: S-expression parsing and tokenization
- **`form_processor`**: Form processing utilities
- **`persistent`**: Persistent data structures for immutable state

## Future Development

Planned enhancements for EL2 include:

- **Pattern Matching Improvements**: Full structural pattern matching
- **Type Inference**: More sophisticated type inference algorithms
- **Proof Tactics**: Built-in proof tactics and automation
- **Module System**: Support for modular development
- **IDE Integration**: Language server protocol support

## Related Documentation

- See the main project README for overall architecture
- Check `el/DOC.md` for detailed algorithm descriptions
- Refer to `sorts/` package documentation for type system details
