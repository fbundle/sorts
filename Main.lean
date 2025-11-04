import EL2.Parser
import EL2.Core

open EL2.Parser
open EL2.Core

def s := "
  inh Nat : Type0
  inh zero : Nat
  inh succ: Nat -> Nat

  



  inh Vec : Nat -> Type0 -> Type0
  inh nil : (T: Type0) -> (Vec zero T)
  inh push : (n: Nat) -> (T: Type0) -> (v: (Vec n T)) -> (x: T) -> (Vec (succ n) T)
  let one := (succ zero)
  let two := (succ one)
  let pure: Nat -> (Vec one Nat) := lam x =>
    (push zero Nat (nil Nat) x)
  let pure_two: (Vec one Nat) := (pure two)
  Type0
"

def t := Exp.typ 1

-- private def tokens := tokenize s

def main  : IO Unit := do
  IO.println "--------------------------------------"
  match parse s.toList with
    | none => IO.println "parse_error"
    | some (e, rest) =>
      IO.println s!"{e}"
      match rest with
        | [] =>
          if true = typeCheck? e t then
            IO.println "passed"
          else
            IO.println "type_error"
        | _ => IO.println "parse_error"
