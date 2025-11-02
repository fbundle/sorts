# SORTS

The goal of this project is to implement minimal dependent type checker

## TYPE CHECKING WITH COQUAND'S ALGORITHM

```lean
* Nat: type_0
* zero: Nat
* succ: Π(n: Nat) Nat
* Vec: Π(n: Nat) Π(T: type_0) type_0
* nil: Π(T: type_0) (Vec zero T)
* push: Π(n: Nat) Π(T: type_0) Π(v: Vec n T) Π(x: T) (Vec (succ n) T)
let one: Nat := succ zero
let singleton: Vec one Nat := push zero Nat (nil Nat) zero
type_0
```

It is a difficult topic checking of an inductive type is well-defined or at least positive recurrent.

There are two main ways to typecheck inductive types: (1) using fixpoint combinator which is not decidable and (2) using the initial object $\mathbb{N}$ of the category $F$-algebras over the category of sets where $F(X) = 1 + X$ where $1$ is the singleton set and $1 + X$ is the disjoint union.

Both of which are not a simple weekend project, hence I decided to stop here. Currently, we simulate inductive types using `*` operator (or `inhabit`) which basically assume some type inhabits. 

The example above assumed `Nat` is a constant of type `type_0`, `zero` is a constant of type `Nat`, `succ` is a function `Nat -> Nat`, etc.

The original Coquand's algorithm can be found in `exp/coquand/1-s2.0-0167642395000216-main.pdf`, the implementation in Haskell is at `exp/coquand/app/Main.hs`, the implementation in Lean4 is at `EL2/Core/Coquand.lean`

I respectively added type universes, inhabit, annotated type, and desugaring for application of untyped lambda.