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

inductive Code (β: Type) where
  | atom: β → Code β
  | beta: Beta (Code β) → Code β
  | annot: Annot (Code β) → Code β
  | binding: Binding (Code β) → Code β
  | typeof: Typeof (Code β) → Code β
  | inh: Inh (Code β) → Code β
  | pi: Pi (Code β) → Code β
  deriving Repr


abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form

private partial def parseWithHead (parseList: List Form → Option (Code β)) (head: String) (form: Form): Option (Code β) :=
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

private partial def parseListAnnot (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure (Code.annot {name := name, type := type})

private partial def parseListBinding (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let valueForm ← list[1]?
  let value ← parse valueForm
  pure (Code.binding {name := name, value := value})

private partial def parseListTypeof (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let valueForm ← list[0]?
  let value ← parse valueForm
  pure (Code.typeof {value := value})

private partial def parseListInh (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let typeForm ← list[0]?
  let type ← parse typeForm
  pure (Code.inh {type := type})

private partial def parseListPi (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  if list.length = 0 then
    none
  else
    let paramForms := list.extract 0 (list.length-1)
    let params ← Util.applyAll paramForms ((λ form => do
      let annotCode ← parseWithHead (parseListAnnot parse) ":" form
      match annotCode with
        | Code.annot annot => some annot
        | _ => none
    ): Form → Option (Annot (Code β)))
    let bodyForm ← list[list.length-1]?
    let body ← parse bodyForm
    pure (Code.pi {params := params, body := body})

private partial def parseBeta (parse: Form → Option (Code β)) (form: Form): Option (Code β) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.applyAll xs parse
      pure (Code.beta {cmd := cmd, args := args})
    | _ => none

partial def parse (parseName: Form → Option (Code β)) (form: Form): Option (Code β) := do
  let rec loop (parseLists: List (Form → Option (Code β))) (form: Form): Option (Code β) :=
    match parseLists with
      | [] => none
      | parseList :: parseLists =>
        match parseList form with
          | some code => code
          | none => loop parseLists form

  loop [
    parseName,
    parseWithHead (parseListAnnot (parse parseName)) ":",
    parseWithHead (parseListBinding (parse parseName)) ":=",
    parseWithHead (parseListTypeof (parse parseName)) "&",
    parseWithHead (parseListInh (parse parseName)) "*",
    parseWithHead (parseListPi (parse parseName)) "=>",
    parseBeta (parse parseName),
  ] form


private partial def parseName (form: Form): Option (Code String) :=
  match form with
    | .name name => some (.atom name)
    | _ => none

def _example: List (Code String) :=
  let source := "
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
  match Form.defaultParseAll source with
    | none => []
    | some xs =>

    Util.applySome xs (parse parseName)




#eval _example





end Code
