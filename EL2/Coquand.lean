-- adapted from (Thierry Coquand - An algorithm for type-checking dependent types)
-- the algorithm is able to type check dependently-typed λ-calculus

namespace EL2.Coquand

structure Map α where
  list: List (String × α)

instance [Repr α]: Repr (Map α) where
  reprPrec (map: Map α) (_: Nat): Std.Format := repr map.list

partial def Map.lookup? [Repr α] (map: Map α) (name: String): Option α :=
  match map.list with
    | [] => none
    | (key, val) :: list =>
      if name = key then
        some val
      else
        {list := list: Map α}.lookup? name

partial def Map.update (map: Map α) (name: String) (val: α): Map α :=
  {list := (name, val) :: map.list}

def emptyMap: Map α := {list := []}


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
  -- Π type: Π (name: type) body
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp
  deriving Repr

inductive Val where
  -- type_0 type_1 ...
  | type: Val
  -- generic value
  | gen: (i: Nat) → Val
  -- application
  | app: (cmd: Val) → (arg: Val) → Val
  -- with closure - a future value - evaluated by eval?
  | clos: (map: Map Val) → (exp: Exp) → Val
  deriving Repr

-- a short way of writing the whnf algorithm
mutual
partial def app? (cmd: Val) (arg: Val): Option Val := do
  match cmd with
    | Val.clos env (Exp.abs name body) =>
      eval? (env.update name arg) body

    | _ =>
      pure (Val.app cmd arg)

partial def eval? (env: Map Val) (exp: Exp): Option Val := do
  match exp with
    | Exp.type => pure Val.type

    | Exp.var name =>
      env.lookup? name

    | Exp.app cmd arg =>
      app? (← eval? env cmd) (← eval? env arg)

    | Exp.bnd name val _ body =>
      eval? (env.update name (← eval? env val)) body

    | _ => pure (Val.clos env exp)
end

partial def whnf? (val: Val): Option Val := do
  match val with
    | Val.app cmd arg =>
      app? (← whnf? cmd) (← whnf? arg)

    | Val.clos env exp =>
      eval? env exp

    | _ => pure val

-- definitional equality
-- the conversion algorithm; the integer is used to represent the introduction of a fresh variable
partial def eqVal? (k: Nat) (u1: Val) (u2: Val): Option Bool := do
  let b: Option Bool := do
    match (← whnf? u1, ← whnf? u2) with
      | (Val.type, Val.type) =>
        pure true

      | (Val.app cmd1 arg1, Val.app cmd2 arg2) =>
        pure ((← eqVal? k cmd1 cmd2) ∧ (← eqVal? k arg1 arg2))

      | (Val.gen k1, Val.gen k2) =>
        pure (k1 == k2)

      | (Val.clos env1 (Exp.abs name1 body1), Val.clos env2 (Exp.abs name2 body2)) =>
        let v := Val.gen k
        eqVal? (k + 1)
          (Val.clos (env1.update name1 v) body1)
          (Val.clos (env2.update name2 v) body2)

      | (Val.clos env1 (Exp.pi name1 type1 body1), Val.clos env2 (Exp.pi name2 type2 body2)) =>
        let v := Val.gen k
        pure (
          (← eqVal? k (Val.clos env1 type1) (Val.clos env2 type2))
            ∧
          (
            ← eqVal? (k + 1)
            (Val.clos (env1.update name1 v) body1)
            (Val.clos (env2.update name2 v) body2)
          )
        )
      | _ => pure false

  if b = false then
    dbg_trace s!"eqVal? failed: {repr u1} {repr u2}"
    b
  else
    b

-- type checking and type inference

structure Ctx where
  k: Nat
  ρ : Map Val -- name -> value
  Γ: Map Val -- name -> type

def Ctx.bind (ctx: Ctx) (name: String) (val: Val) (typ: Val) : Ctx :=
  {
    k := ctx.k,
    ρ := ctx.ρ.update name val,
    Γ := ctx.Γ.update name typ,
  }

def Ctx.intro (ctx: Ctx) (name: String) (typ: Val) : Ctx × Val :=
  let val := Val.gen ctx.k
  (ctx.bind name val typ, val)

def emptyCtx: Ctx := {k := 0, ρ := emptyMap, Γ := emptyMap}


mutual

partial def checkType? (ctx: Ctx) (e: Exp): Option Bool :=
  checkExp? ctx e Val.type

partial def checkExp? (ctx: Ctx) (e: Exp) (v: Val): Option Bool := do
  -- check if expr e is of type v
  match (e, ← whnf? v) with
    | (Exp.abs x n, Val.clos env (Exp.pi y a b)) =>
      let (subCtx, v) := ctx.intro x (Val.clos env a)
      checkExp? subCtx n (Val.clos (env.update y v) b)

    | (Exp.pi x a b, Val.type) =>
      if ¬ (← checkType? ctx a) then
        pure false
      else
        let (subCtx, _) := ctx.intro x (Val.clos ctx.ρ a)
        checkType? subCtx b

    | (Exp.bnd x e1 e2 e3, _) =>
      if ¬ (← checkType? ctx e2) then
        pure false
      else if ¬ (← checkExp? ctx e1 (Val.clos ctx.ρ e2)) then
        pure false
      else
        checkExp? (ctx.bind x
          (← whnf? (Val.clos ctx.ρ e1))
          (← whnf? (Val.clos ctx.ρ e2))
        ) e3 v

    | _ => eqVal? ctx.k (← inferExp? ctx e) v

partial def inferExp? (ctx: Ctx) (e: Exp): Option Val := do
  -- infer type of expr e
  match e with
    | Exp.var name => ctx.Γ.lookup? name
    | Exp.app e1 e2 =>
      match (← whnf? (← inferExp? ctx e1)) with
        | Val.clos env (Exp.pi x a b) =>
          if ← checkExp? ctx e2 (Val.clos env a) then
            pure (Val.clos (env.update x (Val.clos ctx.ρ e2)) b)
          else
            none
        | _ => none
    | Exp.type => Val.type
    | _ => none
end

def typeCheck (m: Exp) (a: Exp): Option Bool := do
  if ¬ (← checkType? emptyCtx a) then
    pure false
  else if ¬ (← checkExp? emptyCtx m (Val.clos emptyMap a)) then
    pure false
  else
    pure true

private def test :=
  typeCheck
    (Exp.abs "A" (Exp.abs "x" (Exp.var "x")))
    (Exp.pi "B" Exp.type (Exp.pi "y" (Exp.var "B") (Exp.var "B")))

#eval test

end EL2.Coquand
