import EL2.EL2
import WBT.WBTArr
import Std

open EL2.Term
open WBT

def term : Term :=
  let x := bnd {
    init := [
      -- inductive Nat : U_1 where
      --  | zero: inh [] Nat 1
      --  | succ: (n: Nat) => inh Nat succ n
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
          body := inh {type := var "Nat", cons := "succ", args := [var "n"]},
        },
      },
      -- inductive Vec: U_1 where
      --  | nil: (T: U_1) => inh (Vec zero T) nil T
      --  | append: (n: Nat) (T: U_1) (vec: Vec n T) (last: T) => inh (Vec (succ n) T) append n T vec last
      {
        name := "Vec",
        value := lam {
          params := [{name := "n", type := var "Nat"}, {name := "T", type := univ 1}],
          body := inh {type := univ 1, cons := "Vec", args := [var "n", var "T"]},
        },
      },
      {
        name := "nil",
        value := lam {
          params := [{name := "T", type := univ 1}],
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
          body := mat {
            cond := var "n",
            cases := [
              {patCons := "zero", patArgs := [], value := app {cmd := var "append", args := [var "n", var "T", var "vec", var "val"]}},
              {patCons := "succ", patArgs := ["_"], value := var "vec"},
            ],
          },
        },
      },
      {
        name := "id_Succ",
        value := lam {
          params := [
            {
              name := "f",
              type := lam {
                params := [
                  {name := "m", type := var "Nat"},
                ],
                body := var "Nat",
              },
            },
          ],
          body := var "f",
        }
      },
      {
        name := "id_AppendIfEmpty",
        value := lam {
          params := [
            {
              name := "f",
              type := lam {
                params := [
                  {name := "m", type := var "Nat"},
                  {name := "T", type := univ 1},
                  {name := "vec", type := app {cmd := var "Vec", args := [var "m", var "T"]}},
                  {name := "val", type := var "T"},
                ],
                body := mat {
                  cond := var "m",
                  cases := [
                    {
                      patCons := "zero", patArgs := [],
                      value := app {cmd := var "Vec", args := [app {cmd := var "succ", args := [var "m"]}, var "T"]},
                    },
                    {
                      patCons := "succ", patArgs := ["_"],
                      value := app {cmd := var "Vec", args := [var "m", var "T"]},
                    },
                  ],
                },
              },
            }
          ],
          body := var "f",
        }
      },
    ],
    last := app {cmd := var "id_AppendIfEmpty", args := [var "append_if_empty"]},
    --last := var "append_if_empty",
    --last := app {
    --  cmd := var "append_if_empty",
    --  args := [var "zero", var "Nat", app {cmd := var "nil", args := [var "Nat"]}, var "one"],
    --},
  }
  x




instance : Map (Std.HashMap String α) α where
  size := Std.HashMap.size
  set := Std.HashMap.insert
  get? := Std.HashMap.get?


def ctxEmpty : Std.HashMap String Infer.InferedType := Std.HashMap.emptyWithCapacity

def main  : IO Unit := do
  -- print program
  -- IO.println s!"[PRINT]\n{term}"
  -- reduce program
  match Infer.inferType? ctxEmpty term with
    | some it => IO.println s!"[OK] type:\n{it.type}"
    | none => IO.println "[ERR]"
