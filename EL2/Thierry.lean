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
  -- type_0 type_1 ...
  | type: Exp
  -- variable
  | var: (name: String) → Exp
  -- application
  | app: (cmd: Exp) → (arg: Exp) → Exp
  -- λ abstraction
  | abs: (name: String) → Exp
  -- let binding: let name: type := value
  | bnd: (name: String) → (value: Exp) → (type: Exp) → (body: Exp) → Exp
  -- Π type
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp

inductive Val where
  -- type_0 type_1 ...
  | type: Val
  -- generic value
  | gen: (level: Int) → Val
  -- application
  | app: (cmd: Val) → (arg: Val) → Val
  -- with closure
  | clos: (env: Env Val) → (term: Exp) → Val

-- a short way of writing the whnf algorithm




end EL2.Thierry
