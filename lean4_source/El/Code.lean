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
  | name: String → Code β
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

private def parseWithHead (parseList: List Form → Option (Code β)) (head: String) (form: Form): Option (Code β) :=
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

private def parseListAnnot (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure (Code.annot {name := name, type := type})

private def parseListBinding (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let valueForm ← list[1]?
  let value ← parse valueForm
  pure (Code.binding {name := name, value := value})

private partial def parseListTypeof (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let valueForm ← list[0]?
  let value ← parse valueForm
  pure (Code.typeof {value := value})

private def parseListInh (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let typeForm ← list[0]?
  let type ← parse typeForm
  pure (Code.inh {type := type})

private def parseListPi (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
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

private def parseBeta (parse: Form → Option (Code β)) (form: Form): Option (Code β) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.applyAll xs parse
      pure (Code.beta {cmd := cmd, args := args})
    | _ => none

partial def parse (parseAtom: String → Option β) (form: Form): Option (Code β) := do
  let rec loop (parseLists: List (Form → Option (Code β))) (form: Form): Option (Code β) :=
    match parseLists with
      | [] => none
      | parseList :: parseLists =>
        match parseList form with
          | some code => code
          | none => loop parseLists form

  let parseName (form: Form): Option (Code β) :=
    match form with
      | .name name =>
        match parseAtom name with
          | some atom => some (.atom atom)
          | none => some (.name name)
      | .list _ => none

  loop [
    parseName,
    parseWithHead (parseListAnnot (parse parseAtom)) ":",
    parseWithHead (parseListBinding (parse parseAtom)) ":=",
    parseWithHead (parseListTypeof (parse parseAtom)) "&",
    parseWithHead (parseListInh (parse parseAtom)) "*",
    parseWithHead (parseListPi (parse parseAtom)) "=>",
    parseBeta (parse parseAtom),
  ] form



inductive Builtin where
  | univ: Int → Builtin
  deriving Repr

private def parseAtom (s: String): Option Builtin := do
    let s ← s.dropPrefix? "U_"
    let s := s.toString
    let i ← s.toInt?

    pure (.univ i) -- universe level i


def _example: List (Code Builtin) :=
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

    Util.applySome xs (parse parseAtom)




#eval _example





end Code
