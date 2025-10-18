import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2

def el2 (_: Unit): List (Code Atom) :=
  [
    .bind_typ {name := "mk_Nat", params := []},
    .bind_mk {
      name := "mk_zero", params := [], type := {cmd := "mk_Nat", args := []},
    },
    .bind_mk {
      name := "mk_succ", params := [{name := "n", type := (.app {cmd := (.var "mk_Nat"), args := []})}], type := {cmd := "mk_Nat", args := []},
    },
    .bind_val {
      name := "Nat", value := (.app {cmd := (.var "mk_Nat"), args := []}),
    },
    .bind_val {
      name := "zero", value := (.app {cmd := (.var "mk_zero"), args := []}),
    },
    .bind_val {
      name := "succ", value := (.lam {
        params := [{name := "n", type := (.var "Nat")}],
        body := (.app {cmd := (.var "mk_succ"), args := [(.var "n")]}),
      }),
    },
    .bind_typ {name := "Vec", params := [{name := "T", type := (.var "U_2")}, {name := "n", type := (.var "Nat")}]},
    .bind_mk {
      name := "mk_nil", params := [
        {name := "T", type := (.var "U_2")},
      ], type := {cmd := "Vec", args := [(.var "T"), (.var "0")]},
    },
    .bind_mk {
      name := "mk_append", params := [
        {name := "T", type := (.var "U_2")}, {name := "n", type := (.var "Nat")},
        {name := "v", type := (.app {cmd := (.var "Vec"), args := [(.var "T"), (.var "n")]})}
      ], type := {cmd := "Vec", args := [(.var "T"), (.app {cmd := (.var "succ"), args := [(.var "n")]})]},
    },
  ]


end EL2_EXAMPLE





def main  : IO Unit := do
  IO.println s!"{repr (EL2_EXAMPLE.el2 ())}"
