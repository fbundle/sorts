import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2

def el2 (_: Unit): Term Atom :=
  .list {
    init := [
      .bind_typ {name := "Nat", params := [], parent := (.var "U_2")},
      .bind_mk {
        name := "zero", params := [], type := {cmd := "Nat", args := []},
      },
      .bind_mk {
        name := "succ", params := [{name := "n", type := (.app {cmd := (.var "Nat"), args := []})}], type := {cmd := "Nat", args := []},
      },
      .bind_typ {name := "Vec", params := [{name := "T", type := (.var "U_2")}, {name := "n", type := (.var "Nat")}], parent := (.var "U_2")},
      .bind_mk {
        name := "nil", params := [
          {name := "T", type := (.var "U_2")},
        ], type := {cmd := "Vec", args := [(.var "T"), (.var "0")]},
      },
      .bind_mk {
        name := "append", params := [
          {name := "T", type := (.var "U_2")}, {name := "n", type := (.var "Nat")},
          {name := "v", type := (.app {cmd := (.var "Vec"), args := [(.var "T"), (.var "n")]})},
          {name := "x", type := (.var "T")},
        ], type := {cmd := "Vec", args := [(.var "T"), (.app {cmd := (.var "succ"), args := [(.var "n")]})]},
      },

      -- code
      .bind_val {
        name := "one", value := (.app {cmd := (.var "succ"), args := [(.var "zero")]}),
      },
      .bind_val {
        name := "two", value := (.app {cmd := (.var "succ"), args := [(.var "one")]}),
      },
      .bind_val {
        name := "three", value := (.app {cmd := (.var "succ"), args := [(.var "two")]}),
      },

      .bind_val {
        name := "f", value := (.lam {
          params := [{name := "_", type := (.var "Nat")}],
          body := .list {
            init := [
              .bind_val {
                name := "l", value := (.app {cmd := (.var "nil"), args := [(.var "Nat")]}),
              },
              .bind_val {
                name := "l", value := (.app {cmd := (.var "append"), args := [(.var "Nat"), (.var "zero"), (.var "l"), (.var "one")]}),
              },
              .bind_val {
                name := "l", value := (.app {cmd := (.var "append"), args := [(.var "Nat"), (.var "one"), (.var "l"), (.var "two")]}),
              },
              .bind_val {
                name := "l", value := (.app {cmd := (.var "append"), args := [(.var "Nat"), (.var "two"), (.var "l"), (.var "three")]}),
              },
            ],
            tail := (.var "l"),
          },
        })
      },

      .bind_val {
        name := "is_pos", value := (.lam {
          params := [{name := "n", type := (.var "Nat")}],
          body := (.mat {
            cond := (.var "n"),
            cases := [
              {pattern := {cmd := "zero", args := []}, value := (.var "zero")},
              {pattern := {cmd := "succ", args := ["m"]}, value := (.var "one")},
            ],
          }),
        })
      },

      (.app {cmd := (.var "f"), args := [(.var "zero")]}),
    ],
    tail :=  (.app {cmd := (.var "is_pos"), args := [(.var "one")]})
  }


end EL2_EXAMPLE


def printLines [ToString α] (lines: List α) : IO Unit :=
  lines.forM IO.println


def main  : IO Unit := do
  let code := EL2_EXAMPLE.el2 ()
  IO.println code
