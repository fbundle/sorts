import El.Form
import El.Util

namespace Code

structure Beta (α: Type) where
  cmd: α
  args: List α
  deriving Repr

structure Annot (α: Type) where
  name: String
  type: α
  deriving Repr

structure Binding (α: Type) where
  name: String
  value: α
  deriving Repr

structure Typeof (α: Type) where
  value: α
  deriving Repr

structure Inh (α: Type) where -- Inhabited
  type: α
  deriving Repr

structure Pi (α: Type) where -- Pi or Lambda
  params: List (Annot α)
  body: α
  deriving Repr

inductive Code where
  | atom: String → Code
  | beta: Beta Code → Code
  | annot: Annot Code → Code
  | binding: Binding Code → Code
  | typeof: Typeof Code → Code
  | inh: Inh Code → Code
  | pi: Pi Code → Code
  deriving Repr


abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form


mutual

private partial def parseListAnnot (list: List Form): Option Code := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure (Code.annot {name := name, type := type})

private partial def parseListBinding (list: List Form): Option Code := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let valueForm ← list[1]?
  let value ← parse valueForm
  pure (Code.binding {name := name, value := value})

private partial def parseListTypeof (list: List Form): Option Code := do
  let valueForm ← list[0]?
  let value ← parse valueForm
  pure (Code.typeof {value := value})

private partial def parseListInh (list: List Form): Option Code := do
  let typeForm ← list[0]?
  let type ← parse typeForm
  pure (Code.inh {type := type})

private partial def parseListPi (list: List Form): Option Code := do
  if list.length = 0 then
    none
  else
    let paramForms := list.extract 0 (list.length-1)
    let params ← Util.applyMany paramForms ((λ form =>
      match parseWithHead parseListAnnot ":" form with
        | Code.annot annot => some annot
        | _ => none
    ): Form → Option (Annot Code))
    let bodyForm ← list[list.length-1]?
    let body ← parse bodyForm
    pure (Code.pi {params := params, body := body})

private partial def parseWithHead (parseList: List Form → Option Code) (head: String) (form: Form): Option Code :=
  match form with
    | .name _ => none
    | .list list =>
      match list with
        | (.name name) :: xs =>
          if name ≠ head then
            none
          else
            parseList xs
        | _ => none


private partial def parseAtom (form: Form): Option Code :=
  match form with
    | .name name => some (.atom name)
    | _ => none

private partial def parseBeta (form: Form): Option Code := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.applyMany xs parse
      pure (Code.beta {cmd := cmd, args := args})
    | _ => none

partial def parse (form: Form): Option Code := do
  let rec loop (parseLists: List (Form → Option Code)) (form: Form): Option Code :=
    match parseLists with
      | [] => none
      | parseList :: parseLists =>
        match parseList form with
          | some code => code
          | none => loop parseLists form

  loop [
    parseAtom,
    parseWithHead parseListAnnot ":",
    parseWithHead parseListBinding ":=",
    parseWithHead parseListTypeof "&",
    parseWithHead parseListInh "*",
    parseWithHead parseListPi "=>",
    parseBeta,
  ] form

end



def _example := "
  (:= Nat (*U_2))
  (:= 0 (*Nat))
  (:= succ (*(-> Nat)))

  (:= 1 (succ 0))
  (:= 2 (succ 0))
  (:= 3 (succ 0))
  (:= 4 (succ 0))
  (:= x 3)
  (:= y 4)

  (+ x y)
"

#eval _example

#eval Form.defaultParser.tokenize _example

#eval (Form.defaultParseAll _example).get!




end Code
