# EL2 Language & Project Documentation

## Table of Contents
1. [Introduction](#introduction)
2. [Language Overview](#language-overview)
    - [Philosophy & Goals](#philosophy--goals)
    - [Main Concepts](#main-concepts)
    - [Syntax & Constructs](#syntax--constructs)
3. [Complete Example](#complete-example)
4. [Type System & Semantics](#type-system--semantics)
5. [Implementation Internals](#implementation-internals)
    - [Parser Architecture](#parser-architecture)
    - [Type Checker Kernel](#type-checker-kernel)
    - [File Structure](#file-structure)
6. [VS Code Extension](#vs-code-extension)
7. [Advanced Notes & Future Work](#advanced-notes--future-work)
8. [References & Resources](#references--resources)

---

## Introduction

EL2 is a minimal dependently-typed functional language, implemented in Lean 4, that aims to:
- Explore core mechanisms behind Calculus of Constructions (CoC) and move towards Calculus of Inductive Constructions (CIC).
- Serve as a learning, teaching, and experimentation platform for type theorists and language designers.
- Provide a simple, readable kernel focusing on type universes, lambda abstraction, and basic inductive types (simulated).

## Language Overview

### Philosophy & Goals
- Minimalism: The kernel focuses on the smallest adequate set of features supporting dependent types and universes.
- Clarity: The implementation is pedagogically clear, with direct mapping from type theory texts.
- Extensibility: The next step is to move beyond basic CoC into richer CIC territory.

### Main Concepts
- **Sorts**: `Type0`, `Type1`, ...
- **Type Universes**: Hierarchical structure analogous to Lean/Coq universes.
- **Inhabitation (`inh`)**: Introduce constant types, constructors, and simulated inductives.
- **Functions & Pi Types**: "hom" for dependent and non-dependent function types.
- **Lambda**: "lam" for abstraction.
- **Let Bindings**: Typed and untyped.
- **Pattern Matching**: Not directly implemented; to be desugared to recursors.

### Syntax & Constructs
- **Type Universes**: `Type0`, `Type1`, ...
- **Inhabit**: `inh name : Type ...`, `inh name : hom ...`
- **Functions**:\
  Dependent Pi: `hom (x : T) -> U`,\
  Non-dependent: `hom T U`
- **Lambda Abstractions**: `lam x y => body`
- **Let Bindings**:
  - Typed: `let x: T := y`
  - Untyped (syntactic sugar for application): `let x := y`
- **Comments**: `--` for single line

#### Example Snippets
```el2
inh Nat : Type0
inh zero : Nat
inh succ : hom Nat -> Nat

inh Vec : hom Nat Type0 -> Type0
inh nil : hom (T: Type0) -> (Vec zero T)
let n := (succ zero)
let singleton: Vec one Nat := (push zero Nat (nil Nat) zero)
```

See [example.el2](./example.el2) for a complete annotated sample.

## Complete Example

Excerpt with explanations (see `example.el2` for full listing):
```el2
inh Bool : Type0
inh true : Bool
inh false : Bool

-- Assume inductive type Nat & Vec exists for tutorial
inh Nat : Type0
inh zero : Nat
inh succ : hom Nat -> Nat
inh Nat_rec : hom
  (P : hom Nat -> Type0)
  (P zero)
  (hom (n : Nat) (P n) -> (P (succ n)))
  (n : Nat) -> (P n)

let one := (succ zero)
let two := (succ one)
let pure: hom Nat -> (Vec one Nat) := lam x => (push zero Nat (nil Nat) x)
let pure_two: (Vec one Nat) := (pure two)

let is_zero : hom Nat -> Bool := lam n =>
  (Nat_rec
    (lam _ => Bool)
    true
    (lam n _ => false)
    n)
```

## Type System & Semantics

- **Universes**:
  - `Type0` is the base universe (contains e.g. `Nat`, Pi types at level 0, etc.)
  - `TypeN` is universe level N.

- **Terms/Expressions**:
  - Variables, applications, lambda abstractions, dependent function types, let bindings, and inhabits.

- **Type Checking**:
  - Based on [Coquand's algorithm](obsolete/coquand/1-s2.0-0167642395000216-main.pdf), specialized for this kernel in `EL2/Core.lean`.
  - Inference, checking, and definitional equality operate mutually.
  - Simulated inductives: no positivity checking, relies on trusted `inh`.

- **Desugaring**:
  - Pattern matching is desugared to recursor use (see README.md for detailed equivalence).
  - Untyped `let` is sugar for immediate lambda/application.

## Implementation Internals

### Parser Architecture
- See `Parser/Combinator.lean` for a monadic combinator library.
- Grammar is implemented in `EL2/Parser.lean`:
  - Parsers for names, applications, universes, abstraction, Pi types, binding, and inhabitation.
  - Handles comments (removes all after `--`).
  - "hom", "lam", "inh", "let" recognized as keywords.

### Type Checker Kernel
- Central logic is in `EL2/Core.lean`.
- Core datatypes:
  - `Exp`: expressions
  - `Val`: evaluated values
  - `Map`: environment for values/types
  - `Ctx`: typing context (universes, envs)
- Implements:
  - Normalization (WHNF), evaluation, weak inference, universe level checks.
  - Definitional equality by recursive, mutual calls
  - Binding and environment manipulation

### File Structure
- `EL2/`: Core logic for expressions, evaluation, type-checking
- `Parser/`: Combinators and parser logic
- `example.el2`: Language and type system demonstration
- `Main.lean`: Main CLI for parsing, typechecking, and printing results
- `tools/vscode-extension/`: Syntax highlighting VSIX for EL2
- `obsolete/`: Previous/prototype implementations (including Haskell and alt. Lean code)

## VS Code Extension
- Syntax highlighting and configuration under `tools/vscode-extension/`:
  - Recognizes `.el2` files
  - Highlights keywords (`lam`, `let`, `inh`, `hom`), numbers, brackets, and comments
- Installation:
  1. Run `npm install -g vsce && vsce package` in the extension directory
  2. Open VS Code, run "Install from VSIX..."

## Advanced Notes & Future Work

- **Inductive Types**: No positivity check; relies on trusted environment. Next goal: full CIC-style inductives (automatic positivity check & recursors).
- **Pattern Matching**: Planned as high-level sugar over inductive recursors (see README for desugaring example).
- **Type Theoretic Foundations**:
  - Currently closely matches CoC. Future plans involve inductive families, pattern matching, etc.
- **Algorithmic Foundations**:
  - See [Coquand's original paper][1] and code in `obsolete/`.
- Further improvement: better error reporting, more advanced parsing, extend language user-friendliness.

## References & Resources

- [example.el2](./example.el2): Language feature demonstration
- [Coquand's original paper][1]
- [Lean 4 docs](https://leanprover.github.io/)
- [Minimalist Type Theory explainer](https://www.andres-erbsen.de/posts/2022-05-17-minimalist-type-theory.html)

[1]: obsolete/coquand/1-s2.0-0167642395000216-main.pdf
