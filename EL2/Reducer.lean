import EL2.Typer

namespace EL2.Reducer
inductive Form where
  | name: String → Form
  | list: List Form → Form
end EL2.Reducer

namespace EL2.Reducer.Internal
open EL2.Typer

partial def reduce? (env: List (String × Exp)) (exp: Exp): Option Form :=
  match exp with
    | Exp.typ level => some (Form.name s!"Type{level}")
    | Exp.var name =>






end EL2.Reducer.Internal
