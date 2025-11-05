import EL2.Typer

namespace EL2.Reducer.Internal
open EL2.Typer

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

partial def reduce? (env: List (String × Exp)) (exp: Exp): Option Exp := do
  match exp with
    | Exp.typ level => some (Exp.typ level)
    | Exp.var name =>
      reduce? env (← lookup? env name)
    | Exp.app cmd arg =>
      let cmd ← reduce? env cmd
      let arg ← reduce? env arg
      match cmd with
        | Exp.lam name body =>
          reduce? (update env name arg) body
        | _ => none
    | Exp.lam _ _ => exp
    | Exp.pi _ _ _ => none
    | Exp.bnd name value _ body =>
      let value ← reduce? env value
      reduce? (update env name value) body

    | Exp.inh name _ body =>
      reduce? (update env name exp) body







end EL2.Reducer.Internal
