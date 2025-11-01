-- extended from (Thierry Coquand - An algorithm for type-checking dependent types)
-- the algorithm is able to type check dependently-typed λ-calculus
-- with type universe (type_0, type_1, ...)
-- added annotated term

-- TODO only Exp, typeCheck? are public

namespace EL2.Core

structure Map α where
  list: List (String × α)

def Map.toString (m: Map α) (toString: α → String): String :=
  "map(" ++ (String.join $ List.intersperse "," $ m.list.map (λ (key, val) =>
    s!"{key} {toString val}"
  )) ++ ")"

partial def Map.lookup? (map: Map α) (name: String): Option α :=
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
  -- typ 0 is at level 2
  -- typ N is at level N + 2
  -- small types are at level 1
  -- terms are at level 0
  | typ : (level: Nat) → Exp
  -- variable
  | var: (name: String) → Exp
  -- application
  | app: (cmd: Exp) → (arg: Exp) → Exp
  -- let binding: let name: type := value
  | bnd: (name: String) → (value: Exp) → (type: Exp) → (body: Exp) → Exp
  -- annotated term
  | ann: (term: Exp) → (type: Exp) → Exp
  -- λ abstraction
  | lam: (name: String) → (body: Exp) → Exp
  -- Π type: Π (name: type) body - type of abs
  | pi:  (name: String) → (type: Exp) → (body: Exp) → Exp
  -- inh
  | inh: (name: String) → (type: Exp) → (body: Exp) → Exp
  deriving Nonempty

def Exp.toString (e: Exp): String :=
  match e with
    | Exp.typ level => s!"type_{level}"
    | Exp.var name => name
    | Exp.app cmd arg => s!"({cmd.toString} {arg.toString})"
    | Exp.bnd name value type body => s!"let {name}: {type.toString} := {value.toString}\n{body.toString}"
    | Exp.ann term type => s!"({term.toString}: {type.toString})"
    | Exp.lam name body => s!"(λ{name} {body.toString})"
    | Exp.pi name type body => s!"(Π{name}: {type.toString}.{body.toString})"
    | Exp.inh name type body => s!"(*{name}: {type.toString}.{body.toString})"

instance: ToString Exp where
  toString := Exp.toString

inductive Val where
  -- typ_n
  | typ : (n: Nat) → Val
  -- generic value
  | gen: (i: Nat) → Val
  -- application
  | app: (cmd: Val) → (arg: Val) → Val
  -- with closure - a future value - evaluated by eval?
  | clos: (map: Map Val) → (exp: Exp) → Val

partial def Val.toString (v: Val): String :=
  match v with
    | Val.typ level => s!"type_{level}"
    | Val.gen i => s!"gen_{i}"
    | Val.app cmd arg => s!"({cmd.toString} {arg.toString})"
    | Val.clos map exp => s!"closure({map.toString Val.toString} {exp.toString})"

instance: ToString Val where
  toString := Val.toString

instance : ToString (Map Val) where
  toString (m: Map Val): String := m.toString Val.toString



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

    | Exp.ann term _ =>
      eval? env term

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
  (λ (o?: Option Bool) => do
    if o? = true then o? else
    dbg_trace s!"[DBG_TRACE] eqVal? {k} {u1} {u2}"
    o?
  ) $ do
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

def Ctx.toString (ctx: Ctx): String :=
  s!"ctx(k={ctx.k}, ρ={ctx.ρ}, Γ={ctx.Γ})"

instance: ToString Ctx where
  toString := Ctx.toString

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

partial def checkTypLevel? (checkExp?: Ctx → Exp → Val → Option Bool) (ctx: Ctx) (exp: Exp) (maxN: Nat): Option Nat :=
  -- if exp is of type TypeN for 0 ≤ N ≤ maxN
  -- return N
  (λ (o?: Option Nat) =>
    match o? with
      | some v =>
        some v
      | none =>
        dbg_trace s!"[DBG_TRACE] checkTypLevel? {ctx}\n\texp = {exp}"
        none
  ) $ do
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

mutual

partial def inferExp? (ctx: Ctx) (exp: Exp): Option Val :=
  -- infer the type of exp
  (λ (o?: Option Val) =>
    match o? with
      | some o => some o
      | none =>
        dbg_trace s!"[DBG_TRACE] inferExp? {ctx}\n\texp = {exp}"
        none
  ) $ do
    match exp with
      | Exp.typ n => pure (Val.typ (n + 1))
      | Exp.var name => ctx.Γ.lookup? name
      | Exp.ann term type =>
        let _ ← checkTypLevel? checkExp? ctx type ctx.maxN
        let b ← checkExp? ctx term (Val.clos ctx.ρ type)
        if b then
          pure (Val.clos ctx.ρ type)
        else
          none

      | Exp.app cmd arg =>
        -- for Exp.app, cmd should be typ, var, or app
        match (← whnf? (← inferExp? ctx cmd)) with
          | Val.clos env (Exp.pi name type body) =>
            if ← checkExp? ctx arg (Val.clos env type) then
              pure (Val.clos (env.update name (Val.clos ctx.ρ arg)) body)
            else
              none
          | _ => none

      | _ => none -- ignore these

partial def checkExp? (ctx: Ctx) (exp: Exp) (val: Val): Option Bool :=
  -- check if type of exp is val
  (λ (o? : Option Bool) =>
    if o? = true then
      o?
    else
      dbg_trace s!"[DBG_TRACE] checkExp? {ctx}\n\texp = {exp}\n\tval = {val}"
      o?
  ) $ do
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
            let i ← checkTypLevel? checkExp? ctx type n
            let (subCtx, _) := ctx.intro name (Val.clos ctx.ρ type)
            let j ← checkTypLevel? checkExp? subCtx body n
            pure (n = (max i j))
          | _ => none

      | Exp.bnd name value type body =>
        let _ ← checkTypLevel? checkExp? ctx type ctx.maxN
        if ¬ ((← checkExp? ctx value (Val.clos ctx.ρ type))) then
          none
        else
          checkExp? (ctx.bind name
            (← whnf? (Val.clos ctx.ρ value))
            (← whnf? (Val.clos ctx.ρ type))
          ) body val

      | Exp.app (Exp.lam name body) arg => -- process untyped lam (λx.y z)
        let argType ← whnf? (← inferExp? ctx arg)
        let argValue ← whnf? (Val.clos ctx.ρ arg)

        let subCtx := ctx.bind name argValue argType
        checkExp? subCtx body val

      | Exp.inh name type body =>
        let _ ← checkTypLevel? checkExp? ctx type ctx.maxN
        let (subCtx, _) := ctx.intro name (Val.clos ctx.ρ type)
        checkExp? subCtx body val


      | _ => eqVal? ctx.k (← inferExp? ctx exp) val

end

def typeCheck? (exp: Exp) (type: Exp): Option Bool :=
  -- typeCheck?
  -- some false - type check error
  -- none - parse error
  checkExp? emptyCtx exp (Val.clos emptyMap type)

def test1 :=
  typeCheck?
    (Exp.lam "B" (Exp.lam "y" (Exp.var "y")))
    (Exp.pi "A" (Exp.typ 0) (Exp.pi "x" (Exp.var "A") (Exp.var "A")))

def test2 :=
  typeCheck?
    (Exp.pi "A" (Exp.typ 0) (Exp.pi "x" (Exp.var "A") (Exp.var "A")))
    (Exp.typ 1)


def test3 :=
  typeCheck? (Exp.typ 0) (Exp.typ 1)


def test4 :=
  typeCheck?
    (Exp.pi "A" (Exp.typ 0) (Exp.pi "x" (Exp.var "A") (Exp.var "A")))
    (Exp.typ 1)


def test5 :=
  -- this is expected to fail
  typeCheck?
    (Exp.app (Exp.lam "x" (Exp.var "x")) (Exp.typ 0))
    (Exp.typ 0)

def test6 :=
  let e := ( id
    $ Exp.inh "Nat" (Exp.typ 0)
    --$ Exp.inh "zero" (Exp.var "Nat")
    --$ Exp.inh "succ" (Exp.pi "n" (Exp.var "Nat") (Exp.var "Nat"))
    $ Exp.typ 0
  )
  let t := Exp.typ 1
  typeCheck? e t



#eval test1
#eval test2
#eval test3
#eval test4
#eval test6

end EL2.Core
