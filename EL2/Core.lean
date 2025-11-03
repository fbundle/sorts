-- extended from (Thierry Coquand - An algorithm for type-checking dependent types)
-- the algorithm is able to type check dependently-typed λ-calculus
-- with type universe (type_0, type_1, ...) and inhabit

namespace EL2.Core

structure Map α where
  list: List (String × α)

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
    | Exp.lam name body => s!"(λ {name} {body.toString})"
    | Exp.pi name type body => s!"(Π ({name}: {type.toString}) {body.toString})"
    | Exp.inh name type body => s!"* {name}: {type.toString}\n{body.toString})"

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

def Map.toString (m: Map α) (toString: α → String): String :=
  "map(" ++ (String.join $ List.intersperse " | " $ m.list.map (λ (key, val) =>
    s!"{key} → {toString val}"
  )) ++ ")"

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

-- WHNF
mutual
partial def whnf? (cmd: Val) (arg: Val): Option Val := do
  -- reduce (Val.app cmd arg)
  match cmd with
    | Val.clos env (Exp.lam name body) =>
      eval? (env.update name arg) body

    | _ => -- cannot reduce
      pure (Val.app cmd arg)

partial def eval? (env: Map Val) (exp: Exp): Option Val := do
  -- reduce (Val.clos env exp)
  match exp with
    | Exp.typ n => pure (Val.typ n)

    | Exp.var name =>
      env.lookup? name

    | Exp.app cmd arg =>
      whnf? (← eval? env cmd) (← eval? env arg)

    | Exp.bnd name val _ body =>
      eval? (env.update name (← eval? env val)) body

    | _ => pure (Val.clos env exp) -- skip lam pi inh
end

-- DEFINITIONAL EQUALITY
partial def eqVal? (k: Nat) (u1: Val) (u2: Val): Option Bool := do
    match (u1, u2) with
      | (Val.typ n1, Val.typ n2) => pure (n1 = n2)

      | (Val.gen k1, Val.gen k2) =>
        pure (k1 = k2)

      | (Val.app cmd1 arg1, Val.app cmd2 arg2) =>
        pure ((← eqVal? k cmd1 cmd2) ∧ (← eqVal? k arg1 arg2))

      | (Val.clos env1 (Exp.lam name1 body1), Val.clos env2 (Exp.lam name2 body2)) =>
        let v := Val.gen k
        eqVal? (k + 1)
          (← eval? (env1.update name1 v) body1)
          (← eval? (env2.update name2 v) body2)

      | (Val.clos env1 (Exp.pi name1 type1 body1), Val.clos env2 (Exp.pi name2 type2 body2)) =>
        let v := Val.gen k
        pure (
          (← eqVal? k (← eval? env1 type1) (← eval? env2 type2))
            ∧
          (
            ← eqVal? (k + 1)
            (← eval? (env1.update name1 v) body1)
            (← eval? (env2.update name2 v) body2)
          )
        )

      | _ => pure false

-- TYPE CHECKING
structure Ctx where
  maxN: Nat     -- max universe level
  debug: Bool   -- whether to dbg_trace

  k: Nat
  ρ : Map Val   -- name -> value
  Γ: Map Val    -- name -> type

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

def Ctx.nodebug (ctx: Ctx) : Ctx :=
  {ctx with debug := false}

def Ctx.printIfNone (ctx: Ctx) (msg: String) (o?: Option α): Option α :=
  match (o?, ctx.debug) with
    | (none, true) => dbg_trace msg; none
    | _ => o?

def Ctx.printIfFalse (ctx: Ctx) (msg: String) (o?: Option Bool): Option Bool :=
  if (o? ≠ true) ∧ (ctx.debug) then
    dbg_trace msg; none
  else
    o?


def emptyCtx: Ctx := {
  maxN := 5,
  debug := true,

  k := 0,
  ρ := emptyMap,
  Γ := emptyMap
}

partial def checkTypLevel? (checkExp?: Ctx → Exp → Val → Option Bool) (ctx: Ctx) (exp: Exp) (maxN: Nat): Option Nat :=
  -- if exp is of type TypeN for 0 ≤ N ≤ maxN
  -- return N
  ctx.printIfNone s!"[DBG_TRACE] checkTypLevel? {ctx}\n\texp = {exp}\n\tmaxLevel = {maxN}" do
  let rec loop (n: Nat): Option Nat := do
    if n > maxN then
      none
    else
      let b ← checkExp? ctx.nodebug exp (Val.typ n)
      if b then
        pure n
      else
        loop (n + 1)
  loop 0

mutual

partial def inferExp? (ctx: Ctx) (exp: Exp): Option Val :=
  -- infer the type of exp
  ctx.printIfNone s!"[DBG_TRACE] inferExp? {ctx}\n\texp = {exp}" do
    match exp with
      | Exp.typ n => pure (Val.typ (n + 1))
      | Exp.var name => ctx.Γ.lookup? name

      | Exp.app cmd arg =>
        -- for Exp.app, cmd should be typ, var, or app
        match ← inferExp? ctx cmd with
          | Val.clos env (Exp.pi name type body) =>
            if ← checkExp? ctx arg (← eval? env type) then
              let argValue ← eval? ctx.ρ arg
              let subEnv := env.update name argValue
              pure (← eval? subEnv body)
            else
              none
          | _ => none

      | _ => none -- ignore these

partial def checkExp? (ctx: Ctx) (exp: Exp) (val: Val): Option Bool :=
  -- check if type of exp is val
  ctx.printIfFalse s!"[DBG_TRACE] checkExp? {ctx}\n\texp = {exp}\n\tval = {val}" do
    match exp with
      | Exp.lam name1 body1 =>
        match val with
          | Val.clos env2 (Exp.pi name2 type2 body2) =>
            let (subCtx, v) := ctx.intro name1 (← eval? env2 type2)
            checkExp? subCtx body1 (← eval? (env2.update name2 v) body2)
          | _ => none

      | Exp.pi name type body =>
        match val with
          | Val.typ n =>
            let i ← checkTypLevel? checkExp? ctx type ctx.maxN
            let (subCtx, _) := ctx.intro name (← eval? ctx.ρ type)
            let j ← checkTypLevel? checkExp? subCtx body ctx.maxN
            pure ((max i j) ≤ n)
          | _ => none

      | Exp.bnd name value type body =>
        let _ ← checkTypLevel? checkExp? ctx type ctx.maxN
        if ¬ (← checkExp? ctx value (← eval? ctx.ρ type)) then
          none
        else
          let subCtx := ctx.bind name
            (← eval? ctx.ρ value)
            (← eval? ctx.ρ type)

          checkExp? subCtx body val

      | Exp.app (Exp.lam name body) arg => -- desugar untyped lam (λx.y z)
        let argType ← ← inferExp? ctx arg
        let argValue ← eval? ctx.ρ arg

        let subCtx := ctx.bind name argValue argType
        checkExp? subCtx body val

      | Exp.inh name type body =>
        let _ ← checkTypLevel? checkExp? ctx type ctx.maxN
        let (subCtx, _) := ctx.intro name (← eval? ctx.ρ type)
        checkExp? subCtx body val


      | _ => eqVal? ctx.k (← inferExp? ctx exp) val

end

def typeCheck? (exp: Exp) (type: Exp): Option Bool := do
  -- typeCheck?
  -- some false - type check error
  -- none - parse error
  checkExp? emptyCtx exp (← eval? emptyMap type)

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

def piMany (params: List (String × Exp)) (body: Exp): Exp :=
  match params with
    | [] => body
    | (name, type) :: rest =>
      Exp.pi name type (piMany rest body)

def appMany (cmd: Exp) (args: List Exp): Exp :=
  match args with
    | [] => cmd
    | arg :: rest =>
      appMany (Exp.app cmd arg) rest

def test6 :=
  let e: Exp := ( id
    $ .inh "Nat" (.typ 0)
    $ .inh "zero" (.var "Nat")
    $ .inh "succ" (.pi "n" (.var "Nat") (.var "Nat"))
    $ .inh "Vec" (piMany [("n", .var "Nat"), ("T", .typ 0)] (.typ 0))
    $ .inh "nil" (piMany [("T", .typ 0)] (appMany (.var "Vec") [.var "zero", .var "T"]))
    $ .inh "push" (
      piMany [
        ("n", .var "Nat"),
        ("T", .typ 0),
        ("v", (appMany (.var "Vec") [.var "n", .var "T"])),
        ("x", .var "T"),
      ]
      (appMany (.var "Vec") [.app (.var "succ") (.var "n"), .var "T"])
    )
    $ .bnd "one" (.app (.var "succ") (.var "zero")) (.var "Nat")
    $ .bnd "singleton" (appMany (.var "push") [
      .var "zero",
      .var "Nat",
      (.app (.var "nil") (.var "Nat")),
      .var "zero",
    ]) (appMany (.var "Vec") [.var "one", .var "Nat"])
    $ .typ 0
  )
  let t := Exp.typ 1
  typeCheck? e t

#eval test1
#eval test2
#eval test3
#eval test4
#eval test6

end EL2.Core




-- TODO only Exp, typeCheck? are public
export EL2.Core (Exp typeCheck?)
