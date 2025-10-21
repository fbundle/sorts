import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2


def termList : List Term := [
  -- inductive Nat : U_1 where
  --  | zero: inh [] Nat 1
  --  | succ: (n: Nat) => inh [n] Nat 1
  bind {
    name := "Nat",
    value := inh {type := univ 1, cons := "Nat", args := []},
  },
  bind {
    name := "zero",
    value := inh {type := var "Nat", cons := "zero", args := []},
  },
  bind {
    name := "succ",
    value := lam {
      params := [{name := "n", type := var "Nat"}],
      body := inh {type := var "Nat", cons := "succ", args := [var "n"]},
    },
  },

  -- inductive Vec: U_1 where
  --  | nil: (T: U_1) => inh
  bind {
    name := "Vec",
    value := lam {
      params := [{name := "n", type := var "Nat"}, {name := "T", type := univ 1}],
      body := inh {type := univ 1, cons := "Vec", args := [var "n", var "T"]},
    },
  },
  bind {
    name := "nil",
    value := lam {
      params := [{name := "T", type := univ 1}],
      body := inh {
        type := app {cmd := var "Vec", args := [var "zero", var "T"]},
        cons := "nil", args := [var "T"],
      },
    },
  },
  bind {
    name := "append",
    value := lam {
      params := [
        {name := "n", type := var "Nat"},
        {name := "T", type := univ 1},
        {name := "vec", type := app {cmd := var "Vec", args := [var "n", var "T"]}},
        {name := "last", type := var "T"},
      ],
      body := inh {
        type := app {cmd := var "Vec", args := [app {cmd := var "succ", args:= [var "n"]}, var "T"]},
        cons := "append", args := [var "n", var "T", var "vec", var "last"],
      },
    },
  },
  -- code
  bind {
    name := "one",
    value := app {cmd := var "succ", args := [var "zero"]},
  },
  bind {
    name := "two",
    value := app {cmd := var "succ", args := [var "one"]},
  },
  bind {
    name := "three",
    value := app {cmd := var "succ", args := [var "two"]},
  },

  bind {
    name := "append_if_empty",
    value := lam {
      params := [
        {name := "n", type := var "Nat"},
        {name := "T", type := univ 1},
        {name := "vec", type := app {cmd := var "Vec", args := [var "n", var "T"]}},
        {name := "val", type := var "T"},
      ],
      body := mat {
        cond := var "n",
        cases := [
          {patCmd := "zero", patArgs := [], value := app {cmd := var "append", args := [var "n", var "T", var "vec", var "val"]}},
          {patCmd := "succ", patArgs := ["_"], value := var "vec"},
        ],
      },
    },
  },
]



end EL2_EXAMPLE


def main  : IO Unit := do
  let termList := EL2_EXAMPLE.termList
  -- print program
  IO.println (lst {init := termList, last := univ 0})
