# EL Pattern Matching and Type Checking Algorithms

This document provides a detailed explanation of the pattern matching and type checking algorithms in the `el` package, including their design, implementation, and practical usage. It is based on a careful reading of the code in `match.go`, `type.go`, `expr_more.go`, `expr.go`, and `runtime.go`.

---

## Overview

The `el` language supports expressive pattern matching and type checking, inspired by functional programming languages. Pattern matching allows destructuring and binding of data, while type checking ensures expressions conform to expected sorts (types). Both are implemented in a way that supports extensibility and compositionality.

---

## Pattern Matching Algorithm

Pattern matching in `el` is centered around the `Match` expression and the `matchPattern` function. Patterns can be exact values, variable bindings, or structured patterns (e.g., function calls).

### Key Types
- **Expr**: The core interface for all expressions.
- **Match**: Represents a pattern match expression: `(match cond comp1 value1 ... compN valueN final)`.
- **Case**: A pair of a pattern (`Comp`) and a value (`Value`).
- **Exact, Term, FunctionCall**: Types of patterns that can be matched against.
- **Frame**: The environment mapping terms to values and sorts, threaded through all computations.

### High-Level Flow
1. **Evaluate the condition** to obtain its value and sort.
2. **Iterate through cases**:
   - For each case, attempt to match the condition value against the pattern using `matchPattern`.
   - If a match succeeds, update the environment (`Frame`) with any new bindings and evaluate the corresponding value.
3. **If no cases match**, evaluate and return the final/default value.

### Example (from `example.el`)
```lisp
(is_two {Nat -> Bool} {x => (match x
    (exact (succ n1))   True
                        False
)})
```
- If `x` is exactly `succ n1`, returns `True`; otherwise, returns `False`.

### Algorithm: `matchPattern`
Defined in `match.go`:
```go
func matchPattern(frame Frame, condSort sorts.Sort, condValue Expr, pattern Expr) (Frame, bool, error)
```
#### Pseudocode
```
function matchPattern(frame, condSort, condValue, pattern):
    if pattern is Exact:
        resolvedPattern = resolve(pattern.Expr, frame)
        return (frame, resolvedPattern == condValue)
    else if pattern is Term:
        frame' = frame.set(pattern, condSort, condValue)
        return (frame', true)
    else if pattern is FunctionCall:
        if condValue is FunctionCall:
            (frame', matched, err) = matchPattern(frame, cmdSort, cmdValue, pattern.Cmd)
            if not matched or err:
                return (frame', false, err)
            return matchPattern(frame', argSort, argValue, pattern.Arg)
        else:
            return (frame, false)
    else:
        return (frame, false)
```

### Algorithm: `Match.Resolve`
Defined in `expr_more.go`:
```go
func (m Match) Resolve(frame Frame) (Frame, sorts.Sort, Expr, error)
```
#### Pseudocode
```
function resolveMatch(m, frame):
    (frame, condSort, condValue, err) = resolve(m.Cond, frame)
    for each case in m.Cases:
        (frame', matched, err) = matchPattern(frame, condSort, condValue, case.Comp)
        if matched:
            return resolve(case.Value, frame')
    return resolve(m.Final, frame)
```

---

## Type Checking Algorithm

Type checking ensures that expressions conform to expected sorts. It is compositional and supports lambdas, pattern matches, and more.

### Key Types
- **Frame**: The environment for type checking, mapping terms to sorts and values.
- **typeCheckBinding**: Checks if an expression matches a given sort (type).
- **reverseMatchPattern**: Used for type checking pattern matches, binding variables with sorts.

### High-Level Flow
- For each binding or expression, check if it matches the expected sort.
- For lambdas, check the parameter and body types.
- For matches, check each case and the final/default value.

### Example (from `example.el`)
```lisp
(add {Nat -> Nat -> Nat} {x => y => (match y
    (succ z)    (succ ((add x) z))
                x
)})
```
- The type checker ensures that `add` is a function from two `Nat`s to a `Nat`, and that all branches of the match return the correct type.

### Algorithm: `typeCheckBinding`
Defined in `type.go`:
```go
func (frame Frame) typeCheckBinding(parentSort sorts.Sort, name Term, expr Expr) bool
```
#### Pseudocode
```
function typeCheckBinding(frame, parentSort, name, expr):
    if expr == Undef:
        return true
    else if expr is Lambda:
        callFrame = frame.set(name, parentSort, name)
        (paramSort, bodySort) = decompose(parentSort)
        callFrame = callFrame.set(expr.Param, paramSort, expr.Param)
        return typeCheckBinding(callFrame, bodySort, "", expr.Body)
    else if expr is Match:
        (frame, condSort, _, err) = resolve(expr.Cond, frame)
        for each case in expr.Cases:
            (matchedFrame, err) = reverseMatchPattern(frame, condSort, case.Comp)
            if not typeCheckBinding(matchedFrame, parentSort, "", case.Value):
                return false
        return typeCheckBinding(frame, parentSort, "", expr.Final)
    else:
        (_, sort, _, err) = resolve(expr, frame)
        return sort == parentSort
```

### Algorithm: `reverseMatchPattern`
Defined in `match.go`:
```go
func reverseMatchPattern(frame Frame, condSort sorts.Sort, pattern Expr) (Frame, error)
```
#### Pseudocode
```
function reverseMatchPattern(frame, condSort, pattern):
    if pattern is Exact:
        resolve(pattern.Expr, frame)
        return frame
    else if pattern is Term:
        return frame.set(pattern, condSort, pattern)
    else if pattern is FunctionCall:
        (frame, cmdSort, _) = resolve(pattern.Cmd, frame)
        (A, B) = decompose(cmdSort)
        if not subtype(B, parent(condSort)):
            error
        return reverseMatchPattern(frame, A, pattern.Arg)
    else:
        return frame
```

---

## Proof of Correctness (Informal)

### Pattern Matching
- **Soundness**: The `matchPattern` function only returns `true` if the value structurally matches the pattern, and all variable bindings are consistent with the environment. Recursive structure ensures that only matching subcomponents succeed.
- **Completeness**: All possible pattern forms (Exact, Term, FunctionCall) are handled. If a value matches a pattern, the algorithm will find it.
- **Termination**: The recursion always proceeds on structurally smaller subcomponents (e.g., function arguments), so it terminates for finite expressions.

### Type Checking
- **Soundness**: The `typeCheckBinding` function only returns `true` if the expression conforms to the expected sort, recursively checking all subcomponents and using the environment to ensure variable bindings are well-typed.
- **Completeness**: All expression forms are handled. If an expression is well-typed, the algorithm will accept it.
- **Preservation**: If an expression type checks, then evaluating it (using the same environment) will not produce a type error, assuming the sorts system is sound.

---

## Related Facts from Type Theory

- **Pattern Matching**: Pattern matching is a core feature in many functional languages (e.g., ML, Haskell, Agda). It enables concise and expressive destructuring of data and is closely related to inductive types and algebraic data types.
- **Type Soundness**: The combination of pattern matching and type checking is designed to ensure type soundness: "well-typed programs do not go wrong." This is typically formalized via progress and preservation theorems.
- **Environments/Frames**: The use of an environment (Frame) to track variable bindings and sorts is standard in type theory, corresponding to the context Γ in typing judgments.
- **Structural Recursion**: Both pattern matching and type checking algorithms are structurally recursive, mirroring the inductive definitions of expressions and types in type theory.
- **Subtyping and Sorts**: The use of `sorts.SubTypeOf` and `sorts.TermOf` reflects subtyping and type equality, which are central to advanced type systems.
- **Lambda Calculus**: The treatment of lambdas and function application is rooted in the simply-typed lambda calculus, extended with pattern matching.
- **Decidability**: The algorithms are designed to be decidable for the fragment of the language implemented, as all recursions are on finite structures.

---

## Glossary
- **Frame**: The environment mapping terms to values and sorts.
- **Sort**: The type of an expression.
- **Pattern**: An expression used to match against a value.
- **Case**: A pair of a pattern and a value in a match expression.
- **Lambda**: An anonymous function expression.
- **Exact**: A pattern that matches only if the value is exactly equal.
- **Term**: A variable or constant.
- **FunctionCall**: An application of a function to an argument.

---

## References
- `match.go`: `matchPattern`, `reverseMatchPattern`
- `expr_more.go`: `Match`, `Case`, `Match.Resolve`
- `type.go`: `typeCheckBinding`
- `expr.go`: `Expr`, `Term`, `FunctionCall`
- `runtime.go`: `Frame`, environment management

---

This document is based on a careful reading of the code and practical examples. For further questions or clarifications, please ask!
