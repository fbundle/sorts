# SORTS

The goal of this project is to implement minimal dependent type checker. Currently, it should be able to handle Calculus of Constructions (CoC) with Type universes. Next goal, fully Calculus of Inductive Constructions (CIC) just like lean4 or rocq

![example.el2](https://raw.githubusercontent.com/fbundle/el2/refs/heads/master/screenshots/screenshot2.png)

## TODO

pattern matching as syntactic sugar

```lean
inh Nat : Type0
inh zero : Nat
inh succ : Nat -> Nat

inh nat_rec :
  (P : Nat -> Type0) ->
  (P zero) ->
  ((n : Nat) -> (P n) -> (P (succ n))) ->
  (n : Nat) -> (P n)
```

then 
```lean
match n with
| zero => a
| succ m => b m
```

is equivalent to
```lean
nat_rec (λ _. T) a (λ m rec. b m) n
```

## TYPE CHECKING WITH COQUAND'S ALGORITHM

It is a difficult topic checking of an inductive type is well-defined or at least positive recurrent.

There are two main ways to typecheck inductive types: (1) using fixpoint combinator which is not decidable and (2) using the initial object $\mathbb{N}$ of the category $F$-algebras over the category of sets where $F(X) = 1 + X$ where $1$ is the singleton set and $1 + X$ is the disjoint union.

Both of which are not a simple weekend project, hence I decided to stop here. Currently, we simulate inductive types using `*` operator (or `inhabit`) which basically assume some type inhabits. 

The example above assumed `Nat` is a constant of type `type_0`, `zero` is a constant of type `Nat`, `succ` is a function `Nat -> Nat`, etc.

The original Coquand's algorithm can be found in `obsolete/coquand/1-s2.0-0167642395000216-main.pdf`, the implementation in Haskell is at `obsolete/coquand/app/Main.hs`, the implementation in Lean4 is at `obsolete/Coquand.lean`

I respectively added type universes, inhabit, annotated type, and desugaring for application of untyped lambda.