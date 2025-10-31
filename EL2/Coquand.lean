-- adapted from (Thierry Coquand - An algorithm for type-checking dependent types)

namespace EL2.Coquand

structure Ctx α where
  list: List (String × α)

instance [Repr α]: Repr (Ctx α) where
  reprPrec (ctx: Ctx α) (_: Nat): Std.Format := repr ctx.list

partial def Ctx.lookup? [Repr α] (ctx: Ctx α) (name: String): Option α :=
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
  | abs: (name: String) → (body: Exp) → Exp
  -- let binding: let name: type := value
  | bnd: (name: String) → (value: Exp) → (type: Exp) → (body: Exp) → Exp
  -- Π type
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp
  deriving Repr

inductive Val where
  -- type_0 type_1 ...
  | type: Val
  -- generic value
  | gen: (i: Nat) → Val
  -- application
  | app: (cmd: Val) → (arg: Val) → Val
  -- with closure
  | clos: (ctx: Ctx Val) → (term: Exp) → Val
  deriving Repr

abbrev Env := Ctx Val

-- a short way of writing the whnf algorithm
mutual
partial def app? (cmd: Val) (arg: Val): Option Val := do
  match cmd with
    | Val.clos env (Exp.abs name body) =>
      eval? (env.update name arg) body

    | _ =>
      pure (Val.app cmd arg)

partial def eval? (env: Env) (exp: Exp): Option Val := do
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

partial def whnf? (val: Val): Option Val := do
  match val with
    | Val.app u w =>
      app? (← whnf? u) (← whnf? w)

    | Val.clos env e =>
      eval? env e

    | _ => pure val

-- the conversion algorithm; the integer is
-- used to represent the introduction of a fresh variable

partial def eqVal? (k: Nat) (u1: Val) (u2: Val): Option Bool := do
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
mutual

partial def checkType? (kργ: Nat × Env × Env) (e: Exp): Option Bool :=
  checkExp? kργ e Val.type

partial def checkExp? (kργ: Nat × Env × Env) (e: Exp) (v: Val): Option Bool := do
  let (k, ρ, γ) := kργ
  let v ← whnf? v
  match e with
    | Exp.abs x n =>
      match v with
        | Val.clos env (Exp.pi y a b) =>
          let v := Val.gen k
          checkExp? (
            k+1,
            ρ.update x v,
            γ.update x (Val.clos env a),
          ) n (Val.clos (env.update y v) b)
        | _ => none
    | Exp.pi x a b =>
      match v with
        | Val.type =>
          pure (
            (← checkType? (k, ρ, γ) a)
              ∧
            (← checkType? (
                k + 1,
                ρ.update x (Val.gen k),
                γ.update x (Val.clos ρ a),
              ) b
            )
          )
        | _ => none

    | Exp.bnd x e1 e2 e3 =>
      pure (
        (← checkType? (k, ρ, γ) e2)
          ∧
        (← checkExp? (
          k,
          ρ.update x (← eval? ρ e1),
          γ.update x (← eval? ρ e2),
        ) e3 v)
      )
    | _ => eqVal? k (← inferExp? (k, ρ, γ) e) v

partial def inferExp? (kργ: Nat × Env × Env) (e: Exp): Option Val := do
  let (k, ρ, γ) := kργ
  match e with
    | Exp.var name => γ.lookup? name
    | Exp.app e1 e2 =>
      match (← whnf? (← inferExp? (k, ρ, γ) e1)) with
        | Val.clos env (Exp.pi x a b) =>
          if ← checkExp? (k, ρ, γ) e2 (Val.clos env a) then
            pure (Val.clos (env.update x (Val.clos ρ e2)) b)
          else
            none
        | _ => none
    | Exp.type => Val.type
    | _ => none
end

def typeCheck (m: Exp) (a: Exp): Option Bool := do
  pure (
    (← checkType? (0, emptyCtx, emptyCtx) a)
    ∧
    (← checkExp? (0, emptyCtx, emptyCtx) m (Val.clos emptyCtx a))
  )

private def test :=
  typeCheck
    (Exp.abs "A" (Exp.abs "x" (Exp.var "x")))
    (Exp.pi "A" Exp.type (Exp.pi "x" (Exp.var "A") (Exp.var "A")))

#eval test

end EL2.Coquand
