import El.Code

namespace Code

private def parseWithHead [Irreducible β] (parseList: List Form → Option (Code β)) (head: String) (form: Form): Option (Code β) :=
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

private def parseListAnnot [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure (.annot {name := name, type := type})

private def parseListBinding [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let valueForm ← list[1]?
  let value ← parse valueForm
  pure (.binding {name := name, value := value})

private partial def parseListTypeof [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let valueForm ← list[0]?
  let value ← parse valueForm
  pure (.typeof {value := value})

private def parseListInh [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let typeForm ← list[0]?
  let type ← parse typeForm
  pure (.inh {type := type})

private def parseListPi [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  if list.length = 0 then
    none
  else
    let paramForms := list.extract 0 (list.length-1)
    let params ← Util.optionMapAll paramForms ((λ form => do
      let annotCode ← parseWithHead (parseListAnnot parse) ":" form
      match annotCode with
        | .annot annot => some annot
        | _ => none
    ): Form → Option (Annot (Code β)))
    let bodyForm ← list[list.length-1]?
    let body ← parse bodyForm
    pure (.pi {params := params, body := body})

private def parseListArrow [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let aForm ← list[0]?
  let a ← parse aForm
  let bForm ← list[1]?
  let b ← parse bForm
  pure (.arrow {a := a, b := b})

private def parseListOther [Irreducible β] (head: String) (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let args ← Util.optionMapAll list parse
  pure (.other {head := head, args := args})

private def parseBeta [Irreducible β] (parse: Form → Option (Code β)) (form: Form): Option (Code β) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.optionMapAll xs parse
      pure (.beta {cmd := cmd, args := args})
    | _ => none

partial def parse [Irreducible β]
  (parseAtom: String → Option β)
  (otherHeadList: List String)
  (form: Form): Option (Code β) := do
  let makeParseAtomFunc (parseAtom: String → Option β) (form: Form): Option (Code β) :=
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

  let parseList := parse parseAtom otherHeadList


  let parseFuncList :=
  -- parse name
  [
    makeParseAtomFunc parseAtom,
    parseNameFunc,
  ]
  ++
  -- parse basic
  [
    parseWithHead (parseListAnnot parseList) ":",
    parseWithHead (parseListBinding parseList) ":=",
    parseWithHead (parseListTypeof parseList) "&",
    parseWithHead (parseListInh parseList) "*",
    parseWithHead (parseListPi parseList) "=>",
    parseWithHead (parseListArrow parseList) "->",
  ]
  ++
  -- parse builtin
  otherHeadList.map (λ head =>
    parseWithHead (parseListOther head parseList) head
  )
  ++
  -- parse beta (default case)
  [parseBeta parseList]

  Util.applyOnce parseFuncList form

end Code
