import EL2.Parser
import EL2.Core

open EL2.Parser
open EL2.Core

def s := "
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
"

def t := Exp.typ 1

private def tokens := tokenize s

def main  : IO Unit := do
  IO.println "--------------------------------------"
  match parse tokens with
    | none => IO.println "parse_error"
    | some (rest, e) =>
      IO.println s!"{e}"
      match rest with
        | [] =>
          if true = typeCheck? e t then
            IO.println "passed"
          else
            IO.println "type_error"
        | _ => IO.println "parse_error"
