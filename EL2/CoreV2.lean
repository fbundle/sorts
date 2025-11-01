-- extended from (Thierry Coquand - An algorithm for type-checking dependent types)
-- the algorithm is able to type check dependently-typed λ-calculus
-- with type universe (type_0, type_1, ...)

-- TODO only Exp and typeCheck are public

namespace EL2.CoreV2

def traceOpt (err: String) (o: Option α): Option α :=
  match o with
    | some v => some v
    | none =>
      dbg_trace err
      none

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
  -- typ 0 is type of small types: Nat, Pi, etc.
  -- typ N is at level N + 2
  -- typ 0 is at level 2, small types are at level 1, terms are at level 0
  | typ : (n: Nat) → Exp
  -- variable
  | var: (name: String) → Exp
  -- application
  | app: (cmd: Exp) → (arg: Exp) → Exp
  -- let binding: let name: type := value
  -- TODO remove - just use λ
  | bnd: (name: String) → (value: Exp) → (type: Exp) → (body: Exp) → Exp
  -- λ abstraction
  | lam: (name: String) → (body: Exp) → Exp
  -- Π type: Π (name: type) body - type of abs
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp
  -- pair
  | pair: (fst: Exp) → (snd: Exp) → Exp
  -- Σ type: Σ (fstName: fstType) sndType
  | sigma: (fstName: String) → (fstType: Exp) → (sndType: Exp) → Exp
  -- fst snd
  | fst: (pair: Exp) → Exp
  | snd: (pair: Exp) → Exp
  -- Eq(a, b) equality type, refl - proof for equality - typecheck by definitional equality
  | eq: (a: Exp) →  (b: Exp) → Exp
  | refl: Exp
  deriving Repr

inductive Val where
  -- typ_n
  | typ : (n: Nat) → Val
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
    | Val.clos env (Exp.lam name body) =>
      eval? (env.update name arg) body

    | _ =>
      pure (Val.app cmd arg)

partial def eval? (env: Map Val) (exp: Exp): Option Val := do
  match exp with
    | Exp.typ n => pure (Val.typ n)

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

-- the conversion algorithm; the integer is used to represent the introduction of a fresh variable
partial def eqVal? (k: Nat) (u1: Val) (u2: Val): Option Bool := do
  -- definitional equality
  traceOpt s!"[DBG_TRACE] eqVal? {k} {repr u1} {repr u2}" do
    match (← whnf? u1, ← whnf? u2) with
      | (Val.typ n1, Val.typ n2) => pure (n1 = n2)

      | (Val.app cmd1 arg1, Val.app cmd2 arg2) =>
        pure ((← eqVal? k cmd1 cmd2) ∧ (← eqVal? k arg1 arg2))

      | (Val.gen k1, Val.gen k2) =>
        pure (k1 = k2)

      | (Val.clos env1 (Exp.lam name1 body1), Val.clos env2 (Exp.lam name2 body2)) =>
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

-- type checking and type inference

structure Ctx where
  maxN: Nat
  k: Nat
  ρ : Map Val -- name -> value
  Γ: Map Val -- name -> type
  deriving Repr

def Ctx.bind (ctx: Ctx) (name: String) (val: Val) (type: Val) : Ctx :=
  {
    ctx with
    ρ := ctx.ρ.update name val,
    Γ := ctx.Γ.update name type,
  }

def Ctx.intro (ctx: Ctx) (name: String) (type: Val) : Ctx × Val :=
  let val := Val.gen ctx.k
  ({
    ctx with
    k := ctx.k + 1,
    ρ := ctx.ρ.update name val,
    Γ := ctx.Γ.update name type,
  }, val)

def emptyCtx: Ctx := {
  maxN := 5,
  k := 0,
  ρ := emptyMap,
  Γ := emptyMap,
}


mutual

partial def inferTypLevel? (ctx: Ctx) (exp: Exp) (maxN: Nat): Option Nat :=
  -- if exp is of type TypeN for 0 ≤ N ≤ maxN
  -- return N
  -- helper for checkExp?
  let rec loop (n: Nat): Option Nat := do
    if n > maxN then
      none
    else
      let b ← checkExp? ctx exp (Val.typ n)
      if b then
        pure n
      else
        loop (n + 1)
  loop 0

partial def inferExp? (ctx: Ctx) (exp: Exp): Option Val := do
  -- infer the type of exp - helper for checkExp?
  traceOpt s!"[DBG_TRACE] inferExp? {repr ctx}\n\texp = {repr exp}" do
    match exp with
      | Exp.typ n => pure (Val.typ (n + 1))
      | Exp.var name => ctx.Γ.lookup? name

      | Exp.app cmd arg =>
        -- for Exp.app, cmd should be typ, var, or app
        -- TODO possibly annotated term ann (x: T)
        -- so that we can do (λx.x : A → A)y instead of let z: A → A := λx.x in y
        match (← whnf? (← inferExp? ctx cmd)) with
          | Val.clos env (Exp.pi name type body) =>
            if ← checkExp? ctx arg (Val.clos env type) then
              pure (Val.clos (env.update name (Val.clos ctx.ρ arg)) body)
            else
              none

          | _ => none

      | _ => none -- ignore these


partial def checkExp? (ctx: Ctx) (exp: Exp) (val: Val): Option Bool := do
  -- check if type of exp is val
  traceOpt s!"[DBG_TRACE] checkExp? {repr ctx}\n\texp = {repr exp}\n\tval = {repr val}" do
    match exp with
      | Exp.lam name1 body1 =>
        match ← whnf? val with
          | Val.clos env2 (Exp.pi name2 type2 body2) =>
            let (subCtx, v) := ctx.intro name1 (Val.clos env2 type2)
            checkExp? subCtx body1 (Val.clos (env2.update name2 v) body2)
          | _ => none

      | Exp.pi name type body =>
        match ← whnf? val with
          | Val.typ n =>
            let i ← inferTypLevel? ctx type n
            let (subCtx, _) := ctx.intro name (Val.clos ctx.ρ type)
            let j ← inferTypLevel? subCtx body n
            pure (n = (max i j))
          | _ => none

      | Exp.bnd name value type body =>
        let _ ← inferTypLevel? ctx type ctx.maxN
        if ¬ (← checkExp? ctx value (Val.clos ctx.ρ type)) then
          none
        else
          checkExp? (ctx.bind name
            (← whnf? (Val.clos ctx.ρ value))
            (← whnf? (Val.clos ctx.ρ type))
          ) body val

      | Exp.pair fst snd =>
        match ← whnf? val with
          | Val.clos env2 (Exp.sigma fstName fstType sndType) =>
            if ¬ (← checkExp? ctx fst (Val.clos env2 fstType)) then
              none
            else
              let subCtx := ctx.bind fstName (← whnf? (Val.clos ctx.ρ fst)) (← whnf? (Val.clos env2 fstType))
              checkExp? subCtx snd (Val.clos env2 sndType)
          | _ => none

      | Exp.sigma fstName fstType sndType =>
        match ← whnf? val with
          | Val.typ n =>
            let i ← inferTypLevel? ctx fstType n
            let (subCtx, _) := ctx.intro fstName (Val.clos ctx.ρ fstType)
            let j ← inferTypLevel? subCtx sndType n
            pure (n = (max i j))
          | _ => none

      | Exp.fst pair =>
        let pairType ← inferExp? ctx pair
        match ← whnf? pairType with
          | Val.clos env (Exp.sigma _ fstType _) =>
            eqVal? ctx.k (Val.clos env fstType) val
          | _ => none

      | Exp.snd pair =>
        let pairType ← inferExp? ctx pair
        match ← whnf? pairType with
          | Val.clos env (Exp.sigma fstName fstType sndType) =>
            let (subCtx, v) := ctx.intro fstName (Val.clos env fstType)
            eqVal? subCtx.k (Val.clos (env.update fstName v) sndType) val
          | _ => none

      | Exp.eq a b =>
        match ← whnf? val with
          | Val.typ n =>
            let sameType ← eqVal? ctx.k (← inferExp? ctx a) (← inferExp? ctx b)
            if ¬ sameType then pure false else

            let i ← inferTypLevel? ctx a n
            let j ← inferTypLevel? ctx b n
            pure (n = (max i j))
          | _ => none

      | Exp.refl =>
        match ← whnf? val with
          | Val.clos env2 (Exp.eq a b) =>
            eqVal? ctx.k (Val.clos env2 a) (Val.clos env2 b)
          | _ => none

      | Exp.typ _ => eqVal? ctx.k (← inferExp? ctx exp) val
      | Exp.var _ => eqVal? ctx.k (← inferExp? ctx exp) val
      | Exp.app _ _ => eqVal? ctx.k (← inferExp? ctx exp) val


end

def typeCheck (m: Exp) (a: Exp): Option Bool := do
  -- typeCheck
  -- some false - type check error
  -- none - parse error
  checkExp? emptyCtx m (Val.clos emptyMap a)

def test1 :=
  typeCheck
    (Exp.lam "B" (Exp.lam "y" (Exp.var "y")))
    (Exp.pi "A" (Exp.typ 0) (Exp.pi "x" (Exp.var "A") (Exp.var "A")))

def test2 :=
  typeCheck
    (Exp.pi "A" (Exp.typ 0) (Exp.pi "x" (Exp.var "A") (Exp.var "A")))
    (Exp.typ 1)


def test3 :=
  typeCheck (Exp.typ 0) (Exp.typ 1)


def test4 :=
  typeCheck
    (Exp.pi "A" (Exp.typ 0) (Exp.pi "x" (Exp.var "A") (Exp.var "A")))
    (Exp.typ 1)


def test5 :=
  typeCheck
    (Exp.app (Exp.lam "x" (Exp.var "x")) (Exp.typ 0))
    (Exp.typ 0)

def test6 :=
  typeCheck
    (Exp.eq (Exp.typ 1) (Exp.typ 1))
    (Exp.typ 2)

def test7 :=
  typeCheck
    (Exp.eq (Exp.typ 0) (Exp.typ 1))
    (Exp.typ 2)

#eval test1
#eval test2
#eval test3
#eval test4
#eval test5
#eval test6
#eval test7

end EL2.CoreV2
