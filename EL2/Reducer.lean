import EL2.Typer

namespace EL2
inductive ReExp where
  | const: (name: String) → ReExp
  | exp: (exp: Exp) → ReExp


instance: Coe Exp ReExp where
  coe (exp: Exp) := ReExp.exp exp

def ReExp.toString? (re: ReExp): Option String :=
  match re with
    | ReExp.const name => some name
    | ReExp.exp $ Exp.typ level => some s!"Type{level}"
    | ReExp.exp $ Exp.var _ => none
    | ReExp.exp $ Exp.app cmd arg =>
      some s!"({ReExp.toString? cmd} {ReExp.toString? arg})"
    | ReExp.exp $ Exp.pi name typeA typeB =>
      match name with
        | "_" => some s!"Π {ReExp.toString? typeA} → {ReExp.toString? typeB}"
        | _   => some s!"Π ({name}: {ReExp.toString? typeA}) → {ReExp.toString? typeB}"
    | ReExp.exp $ Exp.lam name body =>
      s!"λ {name} => {ReExp.toString? body}"
    | ReExp.exp $ Exp.bnd _ _ _ _ => none
    | ReExp.exp $ Exp.inh _ _ _ => none

instance: ToString ReExp where
  toString (re: ReExp): String :=
    match re.toString? with
      | none => "none"
      | some s => s

end EL2

namespace EL2.Reducer
open EL2

-- Util
partial def lookup? (env: List (String × α)) (query: String): Option α :=
  match env with
    | [] => none
    | (key, val) :: rest =>
      if query = key then
        some val
      else
        lookup? rest query

partial def update (env: List (String × α)) (name: String) (val: α): List (String × α) :=
  (name, val) :: env

def printOption (msg: α → String) (o?: Option α): Option α :=
  match o? with
    | none => none
    | some a =>
        dbg_trace msg a ; some a

partial def reduce? (env: List (String × ReExp)) (re: ReExp): Option ReExp := do
  match re with
    | ReExp.const name => ReExp.const name
    | ReExp.exp $ Exp.typ level => some (Exp.typ level)
    | ReExp.exp $ Exp.var name =>
      reduce? env (← lookup? env name)
    | ReExp.exp $ Exp.app cmd arg =>
      let cmd ← reduce? env cmd
      let arg ← reduce? env arg
      match cmd with
        | Exp.lam name body =>
          reduce? (update env name arg) body
        | _ => none
    | ReExp.exp $ Exp.lam _ _ => re
    | ReExp.exp $ Exp.pi _ _ _ => re
    | ReExp.exp $ Exp.bnd name value _ body =>
      let value ← reduce? env value

      printOption (λ v => s!"[REDUCE] {name} = {v}") $
      reduce? (update env name value) body
      -- TODO print here

    | ReExp.exp $ Exp.inh name _ body =>
      reduce? (update env name (ReExp.const name)) body

end EL2.Reducer

namespace EL2
def reduce? (e: Exp): Option ReExp :=
  Reducer.reduce? [] e

end EL2
