import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2


def term : Term := bnd {
  init := [
    -- inductive Nat : U_1 where
    --  | zero: inh [] Nat 1
    --  | succ: (n: Nat) => inh [n] Nat 1
    {
      name := "Nat",
      value := inh {type := univ 1, cons := "Nat", args := []},
    },
    {
      name := "zero",
      value := inh {type := var "Nat", cons := "zero", args := []},
    },
    {
      name := "succ",
      value := lam {
        params := [{name := "n", type := var "Nat"}],
        type := var "Nat",
        body := inh {type := var "Nat", cons := "succ", args := [var "n"]},
      },
    },
    -- inductive Vec: U_1 where
    --  | nil: (T: U_1) => inh
    {
      name := "Vec",
      value := lam {
        params := [{name := "n", type := var "Nat"}, {name := "T", type := univ 1}],
        type := univ 1,
        body := inh {type := univ 1, cons := "Vec", args := [var "n", var "T"]},
      },
    },
    {
      name := "nil",
      value := lam {
        params := [{name := "T", type := univ 1}],
        type := app {cmd := var "Vec", args := [var "zero", var "T"]},
        body := inh {
          type := app {cmd := var "Vec", args := [var "zero", var "T"]},
          cons := "nil", args := [var "T"],
        },
      },
    },
    {
    name := "append",
      value := lam {
        params := [
          {name := "n", type := var "Nat"},
          {name := "T", type := univ 1},
          {name := "vec", type := app {cmd := var "Vec", args := [var "n", var "T"]}},
          {name := "last", type := var "T"},
        ],
        type := app {cmd := var "Vec", args := [app {cmd := var "succ", args:= [var "n"]}, var "T"]},
        body := inh {
          type := app {cmd := var "Vec", args := [app {cmd := var "succ", args:= [var "n"]}, var "T"]},
          cons := "append", args := [var "n", var "T", var "vec", var "last"],
        },
      },
    },
    -- code
    {
      name := "one",
      value := app {cmd := var "succ", args := [var "zero"]},
    },
    {
      name := "two",
      value := app {cmd := var "succ", args := [var "one"]},
    },
    {
      name := "three",
      value := app {cmd := var "succ", args := [var "two"]},
    },

    {
      name := "append_if_empty",
      value := lam {
        params := [
          {name := "n", type := var "Nat"},
          {name := "T", type := univ 1},
          {name := "vec", type := app {cmd := var "Vec", args := [var "n", var "T"]}},
          {name := "val", type := var "T"},
        ],
        type := mat {
          cond :=  var "n",
          cases := [
            {patCmd := "zero", patArgs := [], value := app {cmd := var "Vec", args := [var "one", var "T"]}},
            {patCmd := "succ", patArgs := ["_"], value := typ {value := app {cmd := var "Vec", args := [var "n", var "T"]}}},
          ]
        },
        body := mat {
          cond := var "n",
          cases := [
            {patCmd := "zero", patArgs := [], value := app {cmd := var "append", args := [var "n", var "T", var "vec", var "val"]}},
            {patCmd := "succ", patArgs := ["_"], value := var "vec"},
          ],
        },
      },
    },
  ],
  last := app {
    cmd := var "append_if_empty",
    args := [var "zero", var "Nat", var "nil", var "one"],
  },
}

end EL2_EXAMPLE


instance: EL2.Frame (Std.HashMap String EL2.InferedTerm) where
  set := Std.HashMap.insert
  get? := Std.HashMap.get?

def main  : IO Unit := do
  let term := EL2_EXAMPLE.term
  -- print program
  IO.println term
  -- reduce program
  let frame: Std.HashMap String EL2.InferedTerm := Std.HashMap.emptyWithCapacity

  match EL2.reduce? frame term with
    | some iterm => IO.println s!"term: {iterm.term}\ntype: {iterm.type}\nlevel: {iterm.level}"
    | none => IO.println "error"
