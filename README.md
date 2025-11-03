# SORTS

The goal of this project is to implement minimal dependent type checker

## TYPE CHECKING WITH COQUAND'S ALGORITHM

```lean
inh Nat : Type0 in
inh zero : Nat in
inh succ : Π Nat -> Nat in
inh Vec : Π Nat Type0 -> Type0 in
inh nil : Π (T: Type0) -> (Vec zero T) in
inh push : Π (n: Nat) (T: Type0) (v: (Vec n T)) (x: T) -> (Vec (succ n) T) in
let one: Nat := (succ zero) in
let two: Nat := (succ one) in
let pure: Π Nat -> (Vec one Nat) := λ x =>
  (push zero Nat (nil Nat) x)
in
let pure_two: (Vec one Nat) := (pure two) in
Type0
```

It is a difficult topic checking of an inductive type is well-defined or at least positive recurrent.

There are two main ways to typecheck inductive types: (1) using fixpoint combinator which is not decidable and (2) using the initial object $\mathbb{N}$ of the category $F$-algebras over the category of sets where $F(X) = 1 + X$ where $1$ is the singleton set and $1 + X$ is the disjoint union.

Both of which are not a simple weekend project, hence I decided to stop here. Currently, we simulate inductive types using `*` operator (or `inhabit`) which basically assume some type inhabits. 

The example above assumed `Nat` is a constant of type `type_0`, `zero` is a constant of type `Nat`, `succ` is a function `Nat -> Nat`, etc.

The original Coquand's algorithm can be found in `exp/coquand/1-s2.0-0167642395000216-main.pdf`, the implementation in Haskell is at `exp/coquand/app/Main.hs`, the implementation in Lean4 is at `EL2/Core/Coquand.lean`

I respectively added type universes, inhabit, annotated type, and desugaring for application of untyped lambda.