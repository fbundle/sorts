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
    let params ← Util.optionMapAll paramForms ((λ form => do
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
      let args ← Util.optionMapAll xs parse
      pure (Code.beta {cmd := cmd, args := args})
    | _ => none

partial def parse (parseAtom: String → Option β) (form: Form): Option (Code β) := do
  let parseAtomFunc (form: Form): Option (Code β) :=
    match form with
      | .name name =>
        match parseAtom name with
          | some atom => some (.atom atom)
          | none => none
      | _ => none

  let parseNameFunc (form: Form): Option (Code β) :=
    match form with
      | .name name => some (.name name)
      | _ => none

  let parseList := parse parseAtom

  Util.applyOnce [
    parseAtomFunc,
    parseNameFunc,
    parseWithHead (parseListAnnot parseList) ":",
    parseWithHead (parseListBinding parseList) ":=",
    parseWithHead (parseListTypeof parseList) "&",
    parseWithHead (parseListInh parseList) "*",
    parseWithHead (parseListPi parseList) "=>",
    parseBeta parseList,
  ] form

inductive Atom where
  | univ: Int → Atom
  | integer: Int → Atom
  deriving Repr

private def parseInteger (s: String): Option Atom := do
  let i ← s.toInt?
  pure (.integer i) -- integer i

private def parseUniverse (s: String): Option Atom := do
  let s ← s.dropPrefix? "U_"
  let s := s.toString
  let i ← s.toInt?
  pure (.univ i) -- universe level i

private def parseAtom := Util.applyOnce [
  parseInteger,
  parseUniverse,
  λ _ => none,
]




def _example: List (Code Atom) :=
  let source := "
    (:= Nat (*U_2))
    (:= n0 (*Nat))
    (:= succ (*(-> Nat)))

    (:= n1 (succ n0))
    (:= n2 (succ n0))
    (:= x 3)
    (:= y 4)

    (+ x y)
  "
  match Form.defaultParseAll source with
    | none => []
    | some xs =>

    Util.optionMap xs (parse parseAtom)




#eval _example





end Code
