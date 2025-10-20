import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2
notation "atom" x => Term.atom x
notation "univ" x => Term.univ x
notation "var" x => Term.var x
notation "lst" x => Term.lst x
notation "bind_typ" x => Term.bind_typ x
notation "bind_val" x => Term.bind_val x
notation "bind_mk" x => Term.bind_mk x
notation "typ" x => Term.typ x
notation "lam" x => Term.lam x
notation "app" x => Term.app x
notation "mat" x => Term.mat x


def termList : List Term := [
  bind_typ {name := "Nat", params := [], level := 1},
  bind_mk {
    name := "zero",
    params := [],
    type := {cmd := "Nat", args := []},
  },
  bind_mk {
    name := "succ",
    params := [{name := "n", type := app {cmd := var "Nat", args := []}}],
    type := {cmd := "Nat", args := []},
  },
  bind_typ {
    name := "Vec",
    params := [{name := "T", type := univ 1}, {name := "n", type := var "Nat"}],
    level := 1,
  },
  bind_mk {
    name := "nil", params := [
      {name := "T", type := univ 1},
    ],
    type := {cmd := "Vec", args := [var "T", var "0"]},
  },
  bind_mk {
    name := "append", params := [
      {name := "T", type := univ 1},
      {name := "n", type := var "Nat"},
      {name := "v", type := app {cmd := var "Vec", args := [var "T", var "n"]}},
      {name := "x", type := var "T"},
    ],
    type := {cmd := "Vec", args := [var "T", app {cmd := var "succ", args := [var "n"]}]},
  },

  -- code
  bind_val {
    name := "one",
    value := app {cmd := var "succ", args := [var "zero"]},
  },
  bind_val {
    name := "two",
    value := app {cmd := var "succ", args := [var "one"]},
  },
  bind_val {
    name := "three",
    value := app {cmd := var "succ", args := [var "two"]},
  },

  bind_val {
    name := "f", value := lam {
      params := [{name := "_", type := var "Nat"}],
      body := lst {
        init := [
          bind_val {
            name := "l", value := app {cmd := var "nil", args := [var "Nat"]},
          },
          bind_val {
            name := "l", value := app {cmd := var "append", args := [var "Nat", var "zero", var "l", var "one"]},
          },
          bind_val {
            name := "l", value := app {cmd := var "append", args := [var "Nat", var "one", var "l", var "two"]},
          },
          bind_val {
            name := "l", value := app {cmd := var "append", args := [var "Nat", var "two", var "l", var "three"]},
          },
        ],
        last := var "l",
      },
    }
  },

  bind_val {
    name := "is_pos", value := lam {
      params := [{name := "n", type := var "Nat"}],
      body := mat {
        cond := var "n",
        cases := [
          {pattern := {cmd := "zero", args := []}, value := var "zero"},
          {pattern := {cmd := "succ", args := ["m"]}, value := var "one"},
        ],
      },
    }
  },

  app {cmd := var "f", args := [var "zero"]},
  app {cmd := var "is_pos", args := [var "one"]},
]



end EL2_EXAMPLE


def printLines [ToString α] (lines: List α) : IO Unit :=
  lines.forM IO.println



instance: EL2.Context (Std.HashMap String α) α where
  insert (m: Std.HashMap String α) (key: String) (val: α) := m.insert key val
  get? (m: Std.HashMap String α) (key: String) := m.get? key

def ctx : Std.HashMap String EL2.Term := Std.HashMap.emptyWithCapacity


def main  : IO Unit := do
  let termList := EL2_EXAMPLE.termList
  -- print program
  IO.println (lst {init := termList, last := (univ 0)})

  -- type check
  IO.println "type checking ..."
  IO.println ""

  let (ctx, typeList) := EL2.Util.optionCtxMap termList ctx EL2.inferType?
  for (term, type) in List.zip termList typeList do
    IO.println s!"term: {term}"
    IO.println s!"type: {type}"
    IO.println ""

  if h: typeList.length < termList.length then
    let nextTerm := termList[typeList.length]'h
    IO.println s!"type check error at: {nextTerm}"
  else
    pure ()
