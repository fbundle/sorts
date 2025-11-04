# EL2 Language User Guide & Tutorial

Welcome to EL2—a minimalist dependently-typed language, inspired by the Calculus of Constructions, designed for learning and experimentation. This guide walks you through EL2's syntax and features with step-by-step annotated examples.

## Table of Contents
1. [What is EL2?](#what-is-el2)
2. [Getting Started](#getting-started)
3. [Project Structure](#project-structure)
4. [Language Essentials](#language-essentials)
    - [Types and Universes](#types-and-universes)
    - [Declaring Types & Values](#declaring-types--values)
    - [Functions: Pi and Lambda](#functions-pi-and-lambda)
    - [Let Bindings](#let-bindings)
    - [Inductive Types (Simulated)](#inductive-types-simulated)
    - [Recursion & Pattern Matching](#recursion--pattern-matching)
    - [Comments](#comments)
5. [Complete Example Walkthrough](#complete-example-walkthrough)
6. [Tips & Best Practices](#tips--best-practices)
7. [Common Errors & Troubleshooting](#common-errors--troubleshooting)
8. [Further Resources](#further-resources)

---

## What is EL2?

EL2 is a tiny dependently-typed toy language where you can:
- Define custom types and functions (with dependent types).
- Experiment with minimal type theory, similar to Lean or Coq.
- Learn and teach the foundations of dependently-typed functional programming.

## Getting Started

### Running EL2
To typecheck an EL2 file (e.g., `example.el2`), run:

```bash
lake exe Main.lean example.el2
```
(Passed/failed output is shown in the terminal.)

### (Optional) Syntax Highlighting
Install the VS Code extension from `tools/vscode-extension` using VSIX for better editing experience.

## Project Structure
- `example.el2` — Example EL2 code (start here!)
- `Main.lean` — Entry point for CLI typechecker
- `EL2/` — Language internals (not needed for user/programmer)
- `tools/vscode-extension/` — VS Code syntax highlighter

## Language Essentials

### Types and Universes
EL2 uses universe levels:
- `Type0` — The smallest universe (contains types like `Nat`, etc)
- `Type1`, `Type2`, ... — Higher universes as needed

```el2
Type0    -- a sort (the type of basic types)
Type1    -- the next universe
```

### Declaring Types & Values
To state that a type or value exists, use `inh` (short for "inhabit").

```el2
inh Bool : Type0      -- Declare Bool as a type
inh true : Bool       -- Declare the constant true
inh false : Bool      -- Declare the constant false
```

### Functions: Pi and Lambda
#### Dependent Functions (Pi types)
Use `hom` to express function (possibly dependent) types:

```el2
inh succ : hom Nat -> Nat             -- succ: Nat -> Nat
inh Vec : hom Nat Type0 -> Type0      -- Vec: (n : Nat) → Type0 → Type0
```

- `hom (x : T) -> U` = dependent function (`Π (x : T), U` in Lean/Coq)
- `hom T -> U` = non-dependent (`T → U`)

#### Lambda Abstraction
Use `lam ... => ...` to define anonymous functions:

```el2
let id: hom Nat -> Nat := lam x => x

let plus_one: hom Nat -> Nat := lam n => (succ n)
```

multi-paramters `lam` and `hom` will be automatically converted into single-parameter version

### Let Bindings
- **Simple Let** (no type annotation):
  ```el2
  let one := (succ zero)            -- like let one = succ zero
  let two := (succ one)
  ```
  This is just syntactic sugar for immediately applying a lambda:
    - `let x := y` means `((lam x => body) y)`.

- **Typed Let**: type will be checked by the kernel
  ```el2
  let pure: hom Nat -> (Vec one Nat) := lam x => (push zero Nat (nil Nat) x)
  let pure_two: (Vec one Nat) := (pure two)
  ```



### Inductive Types (Simulated)
Inductive types (like `Nat`, `Vec`, etc.) are declared with multiple `inh` lines.

```el2
inh Nat : Type0
inh zero : Nat
inh succ : hom Nat -> Nat
```

No positivity/termination checks are enforced—it is up to you to write safe declarations. Think of `inh` as axioms, if you make the wrong assumption, you can derive/proof anything

### Recursion & Pattern Matching
EL2 does not have first-class pattern matching yet, but you can code recursion via induction principles (recursors) you define or trust as `inh`.

#### Simulated Nat Recursor
```el2
inh Nat_rec : hom
  (P : hom Nat -> Type0)  -- motive
  (P zero)                -- base case
  (hom (n : Nat) (P n) -> (P (succ n)))  -- step
  (n : Nat)               -- input
  -> (P n)                -- output
```
Then, use it as:
```el2
let is_zero : hom Nat -> Bool := lam n =>
  (Nat_rec
    (lam _ => Bool)
    true
    (lam n _ => false)
    n)
```
This is equivalent to Lean/Coq-style pattern matching on `Nat`.

### Comments
Single-line: `-- this is a comment`

---
## Complete Example Walkthrough

Let's look at the start of `example.el2`:
```el2
inh Bool : Type0
inh true : Bool
inh false : Bool

-- Inductive Nat & recursor
inh Nat : Type0
inh zero : Nat
inh succ : hom Nat -> Nat
inh Nat_rec : hom
  (P : hom Nat -> Type0)
  (P zero)
  (hom (n : Nat) (P n) -> (P (succ n)))
  (n : Nat) -> (P n)

-- Vectors
inh Vec : hom Nat Type0 -> Type0
inh nil : hom (T: Type0) -> (Vec zero T)
inh push : hom (n: Nat) (T: Type0) (v: (Vec n T)) (x: T) -> (Vec (succ n) T)

-- Working with let
let one := (succ zero)
let pure: hom Nat -> (Vec one Nat) := lam x => (push zero Nat (nil Nat) x)

-- Using recursion/induction
let is_zero : hom Nat -> Bool := lam n =>
  (Nat_rec
    (lam _ => Bool)
    true
    (lam n _ => false)
    n)

Type0  -- final expression (sort)
```
- Each `inh` declares a type or constant or function.
- `let` binds a value while optionally giving a type.
- The final term should be a `TypeN` (well-typed file/program).

---

## Tips & Best Practices
- When defining recursors, check you provide the correct arity and types.
- Add types to complex `let` bindings for better error messages.
- Use VS Code syntax highlighting for easier development.
- Keep inductives simple; EL2 does not enforce all type-theoretic safety checks.

---
## Common Errors & Troubleshooting

| Error           | Cause & Solution                                  |
|-----------------|--------------------------------------------------|
| `type_error`    | Your top-level term is not well-typed. Check all types, especially in lets & recursors. |
| `parse_error`   | Syntax issue—ensure all parentheses, `inh`, `hom`, and `lam` are spelled/cased right. |
| Misapplied recursor | Check argument count & types for your recursor/inference principle. |
| Unexpected variables | Make sure variables are properly introduced (in Pi/lambda etc). |

If in doubt, simplify: comment out parts of your file and rebuild incrementally.

## Further Resources
- See `example.el2` for full-featured, annotated examples
- Read about dependent type theory: "Calculus of Constructions", Lean/Coq docs
- [Minimalist Type Theory explainer](https://www.andres-erbsen.de/posts/2022-05-17-minimalist-type-theory.html)
- [Lean 4 docs](https://leanprover.github.io/)
- (EL2 source code, unless you want to contribute, can be ignored for end users)

---
Happy experimenting with EL2! For improvements, suggestions, or showcasing cool EL2 code, open an issue or PR.
