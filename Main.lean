import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2


def termList : List Term := [
  -- inductive Nat : U_1 where
  --  | zero: inh [] Nat 1
  --  | succ: (n: Nat) => inh [n] Nat 1
  .bind "Nat" (.inh (.univ 1) "Nat" []),
  .bind "zero" (.inh (.var "Nat") "zero" []),
  .bind "succ" (.lam [("n", (.var "Nat"))] (.inh (.var "Nat") "succ" [(.var "n")])),

  -- inductive Vec: U_1 where
  --  | nil: (T: U_1) => inh
  .bind "Vec" (.lam [("n", (.var "Nat")), ("T", (.univ 1))] (.inh (.univ 1) "Vec" [(.var "n"), (.var "T")])),
  .bind "nil" (.lam [("T", (.univ 1))] (.inh (.app (.var "Vec") [(.var "zero"), (.var "T")]) "nil" [(.var "T")])),
  .bind "append" (.lam [
    ("n", (.var "Nat")),
    ("T", (.univ 1)),
    ("vec", (.app (.var "Vec") [(.var "n"), (.var "T")])),
    ("last", (.var "T")),
  ] (.inh (.app (.var "Vec") [(.app (.var "succ") [(.var "n")]), (.var "T")]) "append" [(.var "n"), (.var "T"), (.var "vec"), (.var "last")])),

  -- code
  .bind "one" (.app (.var "succ") [(.var "zero")]),
  .bind "two" (.app (.var "succ") [(.var "one")]),
  .bind "three" (.app (.var "succ") [(.var "two")]),

  .bind "append_if_empty" (.lam [
    ("n", (.var "Nat")),
    ("T", (.univ 1)),
    ("vec", (.app (.var "Vec") [(.var "n"), (.var "T")])),
    ("val", (.var "T")),
  ] (.mat (.var "n") [
    ("zero", [], .app (.var "append") [(.var "n"), (.var "T"), (.var "vec"), (.var "val")]),
    ("succ", ["_"], (.var "vec")),
  ])),
]



end EL2_EXAMPLE


def main  : IO Unit := do
  let termList := EL2_EXAMPLE.termList
  -- print program
  IO.println (EL2.Term.list termList (EL2.Term.univ 0))
