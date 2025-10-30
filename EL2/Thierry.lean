-- adapted from (Thierry Coquand - An algorithm for type-checking dependent types)

namespace EL2.Thierry

def lift (e: α) (o: Option β): Except α β :=
  match o with
    | none => Except.error e
    | some v => Except.ok v

structure Ctx α where
  list: List (String × α)
  deriving Repr

partial def Ctx.lookup? [Repr α] (ctx: Ctx α) (name: String): Except String α :=
  match ctx.list with
    | [] => Except.error s!"name {name} not found in {repr ctx}"
    | (key, val) :: list =>
      if name = key then
        Except.ok val
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
  | abs: (name: String) → (body: Exp) → Exp
  -- let binding: let name: type := value
  | bnd: (name: String) → (value: Exp) → (type: Exp) → (body: Exp) → Exp
  -- Π type
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp
  deriving Repr

inductive Val where
  -- type_0 type_1 ...
  | type: Val
  -- generic value at level
  | gen: (i: Nat) → Val
  -- application
  | app: (cmd: Val) → (arg: Val) → Val
  -- with closure
  | clos: (ctx: Ctx Val) → (term: Exp) → Val
  deriving Repr

abbrev Env := Ctx Val

-- a short way of writing the whnf algorithm
mutual
partial def app? (cmd: Val) (arg: Val): Except String Val := do
  match cmd with
    | Val.clos env (Exp.abs name body) =>
      eval? (env.update name arg) body

    | _ =>
      pure (Val.app cmd arg)

partial def eval? (env: Env) (exp: Exp): Except String Val := do
  match exp with
    | Exp.var name =>
      env.lookup? name

    | Exp.app cmd arg =>
      app? (← eval? env cmd) (← eval? env arg)

    | Exp.bnd name val _ body =>
      eval? (env.update name (← eval? env val)) body

    | Exp.type => pure Val.type

    | _ => pure (Val.clos env exp)

end

partial def whnf? (val: Val): Except String Val := do
  match val with
    | Val.app u w =>
      app? (← whnf? u) (← whnf? w)

    | Val.clos env e =>
      eval? env e

    | _ => pure val

-- the conversion algorithm; the integer is
-- used to represent the introduction of a fresh variable

partial def eqVal? (k: Nat) (u1: Val) (u2: Val): Except String Bool := do
  let wU1 ← whnf? u1
  let wU2 ← whnf? u2
  match (wU1, wU2) with
    | (Val.type, Val.type) => pure true

    | (Val.app t1 w1, Val.app t2 w2) =>
      pure ((← eqVal? k t1 t2) ∧ (← eqVal? k w1 w2))

    | (Val.gen k1, Val.gen k2) =>
      pure (k1 == k2)

    | (Val.clos env1 (Exp.abs x1 e1), Val.clos env2 (Exp.abs x2 e2)) =>
      let v := Val.gen k
      eqVal? (k + 1)
        (Val.clos (env1.update x1 v) e1)
        (Val.clos (env2.update x2 v) e2)

    | (Val.clos env1 (Exp.pi x1 a1 b1), Val.clos env2 (Exp.pi x2 a2 b2)) =>
      let v := Val.gen k
      pure (
        (← eqVal? k (Val.clos env1 a1) (Val.clos env2 a2))
          ∧
        (← eqVal? (k + 1)
          (Val.clos (env1.update x1 v) b1)
          (Val.clos (env2.update x2 v) b2)
        )
      )
    | _ => pure false

-- type checking and type inference


end EL2.Thierry
