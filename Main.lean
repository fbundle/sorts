
import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2

def el2 (_: Unit): List (Code Atom) :=
  [
    .bind_typ {name := "Nat", params := []},
    .bind_mk {
      name := "zero", params := [], type := {cmd := "Nat", args := []},
    },
    .bind_mk {
      name := "succ", params := [{name := "n", type := (.app {cmd := (.var "Nat"), args := []})}], type := {cmd := "Nat", args := []},
    },
    .bind_typ {name := "Vec", params := [{name := "T", type := (.var "U_2")}, {name := "n", type := (.var "Nat")}]},
    .bind_mk {
      name := "nil", params := [
        {name := "T", type := (.var "U_2")},
      ], type := {cmd := "Vec", args := [(.var "T"), (.var "0")]},
    },
    .bind_mk {
      name := "append", params := [
        {name := "T", type := (.var "U_2")}, {name := "n", type := (.var "Nat")},
        {name := "v", type := (.app {cmd := (.var "Vec"), args := [(.var "T"), (.var "n")]})}
      ], type := {cmd := "Vec", args := [(.var "T"), (.app {cmd := (.var "succ"), args := [(.var "n")]})]},
    },
  ]


end EL2_EXAMPLE





def main  : IO Unit := do
  IO.println s!"{repr (EL2_EXAMPLE.el2 ())}"
