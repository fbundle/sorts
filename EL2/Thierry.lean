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
  | gen: (level: Int) → Val
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
    | Val.clos env (Exp.abs name body) => eval? (env.update name arg) body
    | _ => pure (Val.app cmd arg)

partial def eval? (env: Env) (exp: Exp): Except String Val := do
  match exp with
    | Exp.var name =>
      env.lookup? name
    | Exp.app cmd arg =>
      let cmdVal ← eval? env cmd
      let argVal ← eval? env arg
      app? cmdVal argVal
    | Exp.bnd name val _ body =>
      let valVal ← eval? env val
      eval? (env.update name valVal) body
    | Exp.type => pure Val.type
    | _ => pure (Val.clos env exp)

end

partial def whnf? (val: Val): Except String Val := do
  match val with
    | Val.app u w =>
      let rU ← whnf? u
      let rW ← whnf? w
      app? rU rW
    | Val.clos env e =>
      eval? env e
    | _ => pure val


end EL2.Thierry
