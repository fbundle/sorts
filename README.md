# SORTS

started from some basic dependent type in [lean4](https://github.com/leanprover/lean4), I stated this project with the goal to implement the full dependent type system so that it is capable for mathemtical proof. This is probably a decade-long project, hope it would last.  


# EXAMPLES

## EL — a tiny typed, expression language

EL is a small, homoiconic, s‑expression language with a minimal core and a typed evaluation model built on the `sorts` library. Code and data share the same shape: lists and terms. The evaluator performs stepwise resolution under a persistent frame, tracking both inferred sorts (types) and values.

### Features
- **Homoiconic forms**: programs are made of terms and lists.
- **Function application**: simple two‑element list `(f x)`.
- **Lambdas**: `=>` special form.
- **Let‑bindings**: `let` special form for introducing names with declared sorts.
- **Pattern matching**: `match` with exact and structural patterns.
- **Type constructors**: function arrow `->`, sum `⊕`, product `⊗`.
- **Universes and built‑ins**: `U_n`, `Any_n`, `Unit_n` provide universe levels and terminal/initial sorts.
- **Static type checking**: validates applications against arrow sorts, enforces special‑form arities, and tracks sorts across pattern matches and bindings.

### Structure

This project is organized into several packages:

- **`sorts`**: The core of the project, this package implements a dependent type system. It includes features like Pi (`Π`), Sigma (`Σ`), Arrow (`->`), Sum (`+`), and Product (`×`) types. It also defines the fundamental building blocks such as `Sort`, `Atom`, and `Dependent`.

- **`form`**: This package provides a parser and tokenizer for a simple S-expression-based language. It supports both prefix and infix notations, enabling a flexible and readable syntax.

- **`el`**: An expression language built on top of the `sorts` and `form` packages. It implements features like `let` bindings, `lambda` functions, and `match` expressions. A simple type checker is also included to ensure type safety.

- **`persistent`**: This package provides a set of persistent data structures, including `ordered_map`, `seq`, and `stack`. These data structures are used in the `el` package to implement the runtime environment and ensure immutability.

- **`cmd`**: This directory contains the main applications for the project.
  - **`cmd/el`**: A command-line interface (CLI) for the `el` language, allowing users to execute `el` code and see the results.

## Install / Build

### Requirements
- Go 1.25+

### Build the CLI
```bash
go build -o ./bin/el ./cmd/el
```

Or install into your `GOBIN`:
```bash
go install ./cmd/el
```

Verify:
```bash
echo "(let x U_1 undef x)" | ./bin/el
```

## Language overview

### Syntax
- **Terms**: bare identifiers like `Nat`, `succ`, `x`.
- **Lists**: `(head arg)`; the empty list is not allowed as a term.
- **Application**: a regular two‑element list. Multi‑arg functions are curried: `((add 2) 3)`.

### Special forms

- **Lambda**: `(=> param body)`
  - Produces a function value without immediately evaluating the body.

- **Let**: `(let name1 type1 value1 ... nameN typeN valueN final)`
  - Binds a sequence of names, each with a declared sort and a value; evaluates `final` in the extended frame.
  - Use the literal `undef` to declare a name without assigning a value yet; the name becomes a self‑reference (useful for recursive references and data declarations).

- **Match**: `(match cond comp1 value1 ... compN valueN final)`
  - Tries each pattern `compK` against `cond`. On first match, evaluates `valueK` in the possibly extended frame (pattern variables are bound). If none match, evaluates `final`.

### Patterns
- **Exact**: `(exact expr)` matches only when the resolved value is syntactically equal to `expr`.
- **Variable**: a bare term in pattern position binds to the matched value (and its sort).
- **Structural**: `(f x)` patterns match function calls and recurse over command and argument.

### Sorts (types)
- **Arrow**: `(-> A B)` is the sort of functions from `A` to `B`.
- **Sum**: `(⊕ A B)` disjoint union of sorts.
- **Product**: `(⊗ A B)` pair/product of sorts.
- **Universes**: `U_n` denotes universe level `n+1` of sorts; `Any_n` is the terminal sort at level `n+1`; `Unit_n` is the initial sort at level `n+1`.

### Static typing and guarantees
- **Application checking**: `(f x)` type checks iff `f : (-> A B)` and `x : A`, yielding result sort `B`.
- **Universe inference**: bare `U_n`/`Any_n`/`Unit_n` resolve to atoms at level `n+1`, ensuring level‑correctness for declarations.
- **Pattern binding**: variables bound in `match` patterns carry the matched value’s sort into the branch scope.
- **Special‑form arity**: `=>`, `let`, `match`, `->`, `⊕`, `⊗`, `exact` are arity‑checked at parse time.
- **Error reporting**: ill‑typed applications and malformed forms surface precise errors during resolution.

For the precise algorithms behind pattern matching, reverse pattern matching, and type checking (variable binding and function calling), see `el/ALGORITHM.md`.

Note: binding and match exhaustiveness checks are intentionally conservative and will be tightened further as the language evolves.

## Example

`el/example.el` demonstrates natural numbers, addition, and simple matching:

```el
(let
    Bool U_1                undef
    True Bool               undef
    False Bool              undef

    Nat U_1                 undef
    n0 Nat                  undef
    succ {Nat -> Nat}       undef

    n1 Nat (succ n0)
    n2 Nat (succ n1)
    n3 Nat (succ n2)
    n4 Nat (succ n3)

    x Any_0 {n1 ⊕ n2 ⊕ n3}
    x Any_0 {n1 ⊗ n2 ⊗ n3 ⊗ n4}

    is_two {Nat -> Bool} {x => (match x
        (exact (succ n1))   True
                            False
    )}

    add {Nat -> Nat -> Nat} {x => y => (match y
        (succ z)    (succ ((add x) z))
                    x
    )}

    # (is_two n2)
    (add n2 n3)             # output (succ (succ (succ (succ (succ n0)))))
)
```

Run it with the CLI:

```bash
./bin/el < el/example.el
```

You will see, for each top‑level form, the printed form, its inferred sort, and its resolved value, e.g.:

```text
expr	 (add n2 n3)
sort	 Nat
value	 (succ (succ (succ (succ (succ n0)))))
```

## CLI usage

- **Input**: reads UTF‑8 EL code from stdin; multiple top‑level forms are allowed one after another.
- **Output**: for each form, prints three lines: `expr`, `sort`, `value`.

Examples:
```bash
echo "(=> x x)" | ./bin/el
echo "(let id {Any_0 -> Any_0} {x => x} (id id))" | ./bin/el
```

## Dependent types and proof assistant roadmap

The `sorts` library already implements the following dependent types:

- **Dependent function (Π)**: Also known as a dependent product, this type represents a function where the type of the return value depends on the input value.
- **Dependent pair (Σ)**: Also known as a dependent sum, this type represents a pair where the type of the second element depends on the value of the first.

The following features are planned for the future:

- **Definitional equality and normalization**: Implementing a robust definitional equality check, likely using normalization-by-evaluation, is a crucial next step.
- **Universe polymorphism**: To avoid paradoxes and manage the hierarchy of types, the system will need to support universe polymorphism, including cumulativity and level inference for `Π` and `Σ` types.
- **Equality/type of paths**: The introduction of an identity type with `refl`, `cong`, and `transport` primitives will enable reasoning about equality within the type system.
- **Pattern coverage & termination**: To ensure the reliability of proofs, the system will need to perform totality checks for functions defined by `match`, verifying both pattern coverage and termination.
- **Holes and elaboration**: To improve the user experience and support interactive proof development, the system will feature typed holes, metavariables, and an elaboration mechanism to refine proofs.
- **Tactics/automation**: A small, trusted kernel will be extended with an optional tactic layer for proof search and automation, making it easier to construct complex proofs.

## Notes

- Unicode operators `⊕` and `⊗` are regular identifiers; ensure your editor saves files as UTF‑8.
- EL uses explicit braces `{ ... }` in examples purely for readability; they are not required by the reader—any list is written with parentheses.
- The implementation favors clarity; sorts are reported using names from the `sorts` package.


