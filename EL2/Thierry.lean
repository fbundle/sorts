-- adapted from (Thierry Coquand - An algorithm for type-checking dependent types)

namespace EL2.Thierry

structure Env α where
  list: List (String × α)
  deriving Repr

partial def Env.lookup? (env: Env α) (name: String): Option α :=
  match env.list with
    | [] => none
    | (key, val) :: list =>
      if name = key then
        some val
      else
        {list := list: Env α}.lookup? name

partial def Env.update (env: Env α) (name: String) (val: α): Env α :=
  {list := (name, val) :: env.list}

def emptyEnv: Env α := {list := []}

inductive Exp where
  | type: Exp
  | var: (name: String) → Exp
  | app: (cmd: Exp) → (arg: Exp) → Exp
  | abs: (name: String) → Exp
  | bnd: (name: String) → (value: Exp) → (type: Exp) → (body: Exp) → Exp
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp

inductive Val where
  | type: Val
  | gen: (level: Int) → Val
  | app: (cmd: Val) → (arg: Val) → Val
  | clos: (env: Env Val) → (term: Exp) → Val

-- a short way of writing the whnf algorithm




end EL2.Thierry
