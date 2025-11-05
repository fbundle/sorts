import EL2.Typer

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

inductive ReExp where
  | const: (name: String) → ReExp
  | exp: (exp: Exp) → ReExp

instance: Coe Exp ReExp where
  coe (exp: Exp) := ReExp.exp exp

def ReExp.toString (re: ReExp): String :=
  match re

instance: ToString ReExp where
  toString := ReExp.toString

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
