import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2


def termList : List Term := [
  -- inductive Nat : U_1 where
  --  | zero: Nat
  --  | succ: (n: Nat) → Nat
  .bind "zero" (.inh "zero" [] (.typ "Nat" [] 1)),
  .bind "succ" (.lam [⟨"n", (.var "Nat") ⟩] (.inh "succ" [(.var "n")] (.typ "Nat" [] 1))),

  -- inductive Vec: U_1 where

  .bind "nil" (.lam [⟨"T", (.univ 1)⟩] (.inh "nil" [] (.typ "Vec" [(.var "zero"), (.var "T")] 1))),
  .bind "append" (.lam [
    ⟨"n", (.var "Nat")⟩,
    ⟨"T", (.univ 1)⟩,
    ⟨"vec", (.typ "Vec" [(.var "zero"), (.var "T")] 1)⟩,
    ⟨"last", (.var "T")⟩,
  ] (.inh "append" [(.var "n"), (.var "T"), (.var "vec"), (.var "last")] (.var "Vec"))),

  -- code
  .bind "one" (.app (.var "succ") [(.var "zero")]),
  .bind "two" (.app (.var "succ") [(.var "one")]),
  .bind "three" (.app (.var "succ") [(.var "two")]),

  .bind "append_if_empty" (.lam [
    ⟨"n", (.var "Nat")⟩,
    ⟨"T", (.univ 1)⟩,
    ⟨"vec", (.typ "Vec" [(.var "n"), (.var "T")] 1)⟩,
    ⟨"val", (.var "T")⟩,
  ] (.mat (.var "n") [
    ⟨"zero", [], .app (.var "append") [(.var "n"), (.var "T"), (.var "vec"), (.var "val")]⟩,
    ⟨"succ", ["_"], (.var "vec")⟩,
  ])),
]



end EL2_EXAMPLE


def main  : IO Unit := do
  let termList := EL2_EXAMPLE.termList
  -- print program
  IO.println (EL2.Term.list termList (EL2.Term.univ 0))
