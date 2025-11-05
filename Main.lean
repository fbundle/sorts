import EL2.Parser
import EL2.Core

open EL2.Parser
open EL2.Core

def s := "
  inh Bool : Type0
  inh true : Bool
  inh false : Bool

  -- assume the inductive type Nat exists with induction Nac_rec
  -- for CIC systems like lean
  -- inductive type will be checked for positive recurrent and
  -- induction will be automatically generated

  inh Nat : Type0
  inh zero : Nat
  inh succ : hom Nat -> Nat

  inh Nat_rec : hom
    (P : hom Nat -> Type0)
    (P zero)
    (hom (n : Nat) (P n) -> (P (succ n)))
    (n : Nat) -> (P n)

  -- assume (Vec n T) exists for all type (T: Type0)
  inh Vec : hom Nat Type0 -> Type0
  inh nil : hom (T: Type0) -> (Vec zero T)
  inh push : hom (n: Nat) (T: Type0) (v: (Vec n T)) (x: T) -> (Vec (succ n) T)

  -- some example code

  -- for let syntax, instead of `let x (: typeX) := y in z`
  -- we use new line or semicolon ; instead of in
  -- `let x = y in z`
  -- let syntax without type annotation is basically just syntactic sugar for
  -- ((λ x z) y)
  -- which is just name binding
  let one := (succ zero)
  let two := (succ one)

  -- `let x: typeX := y in z`
  -- let syntax with type annotation will be type-checked by the kernel
  let pure: hom Nat -> (Vec one Nat) := lam x =>
    (push zero Nat (nil Nat) x)
  let pure_two: (Vec one Nat) := (pure two)

  -- pattern matching is just induction
  let is_zero : hom Nat -> Bool := lam n =>
    (Nat_rec
      (lam _ => Bool)      --  for every n, the return type is Bool
      true                --  case zero
      (lam n _ => false)  --  case succ
      n                   --  apply inductive statement on n
    )

  Type0
"

def t := Exp.typ 1



def main (args : List String): IO Unit := do
  IO.println "-------------------------------------------------------------------"
  match args with
    | [] => IO.println "args_empty: use `el2 <filename>`"
    | filename :: _ =>
      let content ← IO.FS.readFile filename
      match parse content.toList with
        | none => IO.println s!"parse_error"
        | some (rest, e) =>
          match rest with
            | [] =>
              if true = typeCheck? e t then
                IO.println "passed"
              else
                IO.println "type_error"
            | _ => IO.println "parse_error"
