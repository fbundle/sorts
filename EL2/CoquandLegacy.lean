-- adapted from (Thierry Coquand - An algorithm for type-checking dependent types)
-- the algorithm is able to type check dependently-typed λ-calculus

namespace EL2.CoquandLegacy

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
  | typ: Exp
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
  | typ: Val
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
    | Exp.typ => pure Val.typ

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
      | (Val.typ, Val.typ) =>
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

def Ctx.bind (ctx: Ctx) (name: String) (val: Val) (type: Val) : Ctx := {
  k := ctx.k,
  ρ := ctx.ρ.update name val,
  Γ := ctx.Γ.update name type,
}

def Ctx.intro (ctx: Ctx) (name: String) (type: Val) : Ctx × Val :=
  let val := Val.gen ctx.k
  (ctx.bind name val type, val)

def emptyCtx: Ctx := {k := 0, ρ := emptyMap, Γ := emptyMap}


mutual

partial def checkType? (ctx: Ctx) (e: Exp): Option Bool :=
  checkExp? ctx e Val.typ

partial def checkExp? (ctx: Ctx) (exp: Exp) (val: Val): Option Bool := do
  -- check if expr e is of type v
  match exp with
    | Exp.abs name1 body1 =>
      match ← whnf? val with
        | Val.clos env2 (Exp.pi name2 type2 body2) =>
          let (subCtx, v) := ctx.intro name1 (Val.clos env2 type2)
          checkExp? subCtx body1 (Val.clos (env2.update name2 v) body2)
        | _ => none

    | Exp.pi name type body =>
      match ← whnf? val with
        | Val.typ =>
          let (subCtx, _) := ctx.intro name (Val.clos ctx.ρ type)
          pure (
            (← checkType? ctx type)
              ∧
            (← checkType? subCtx body)
          )
        | _ => none

    | Exp.bnd name value type body =>
      pure (
        (← checkType? ctx type)
          ∧
        (← checkExp? ctx value (Val.clos ctx.ρ type))
          ∧
        (
          ← checkExp? (ctx.bind name
            (← whnf? (Val.clos ctx.ρ value))
            (← whnf? (Val.clos ctx.ρ type))
          ) body val
        )
      )

    | _ => eqVal? ctx.k (← inferExp? ctx exp) val

partial def inferExp? (ctx: Ctx) (exp: Exp): Option Val := do
  -- infer type of expr e
  match exp with
    | Exp.var name => ctx.Γ.lookup? name
    | Exp.app cmd arg =>
      match (← whnf? (← inferExp? ctx cmd)) with
        | Val.clos env (Exp.pi name type body) =>
          if ← checkExp? ctx arg (Val.clos env type) then
            pure (Val.clos (env.update name (Val.clos ctx.ρ arg)) body)
          else
            none

        | _ => none
    | Exp.typ => Val.typ
    | _ => none
end

def typeCheck (m: Exp) (a: Exp): Option Bool := do
  pure (
    (← checkType? emptyCtx a)
      ∧
    (← checkExp? emptyCtx m (Val.clos emptyMap a))
  )

private def test :=
  typeCheck
    (Exp.abs "A" (Exp.abs "x" (Exp.var "x")))
    (Exp.pi "B" Exp.typ (Exp.pi "y" (Exp.var "B") (Exp.var "B")))

#eval test

namespace EL2.CoquandLegacy
