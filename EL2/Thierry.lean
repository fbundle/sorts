-- adapted from (Thierry Coquand - An algorithm for type-checking dependent types)

namespace EL2.Thierry

structure Ctx α where
  list: List (String × α)
  deriving Repr

partial def Ctx.lookup? (ctx: Ctx α) (name: String): Option α :=
  match ctx.list with
    | [] => none
    | (key, val) :: list =>
      if name = key then
        some val
      else
        {list := list: Ctx α}.lookup? name

partial def Ctx.update (ctx: Ctx α) (name: String) (val: α): Ctx α :=
  {list := (name, val) :: ctx.list}

def emptyCtx: Ctx α := {list := []}

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
  -- generic value at level
  | gen: (level: Int) → Val
  -- application
  | app: (cmd: Val) → (arg: Val) → Val
  -- with closure
  | clos: (ctx: Ctx Val) → (term: Exp) → Val

abbrev Env := Ctx Val

-- a short way of writing the whnf algorithm

mutual
def app (u: Val) (v: Val): Val :=
  sorry

def eval (env: Env) (e: Exp): Val :=
  sorry

end



end EL2.Thierry
