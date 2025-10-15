import EL.Code

namespace EL



private def parseListToParseForm (parseList: List Form → Option α) (head: String) (form: Form) : Option α := do
  match form with
    | .list (.name x :: xs) =>
      if head ≠ x then
        none
      else
        let a ← parseList xs
        pure a
    | _ => none

private def parseFormToParseForm (parse: Form → Option α) (convert: α → β) (form: Form): Option β := do
  let a ← parse form
  let b := convert a
  b

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

private def parseListAnnot [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Annot (Code β)) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure {name := name, type := type}

private def parseAnnot [Irreducible β] (parse: Form → Option (Code β)) (form: Form): Option (Annot (Code β)) :=
  parseListToParseForm (parseListAnnot parse) ":" form

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
    let params ← Util.optionMapAll paramForms (parseAnnot parse)

    let bodyForm ← list[list.length-1]?
    let body ← parse bodyForm

    pure (.pi {params := params, body := body})

private def parsePi  [Irreducible β] (parse: Form → Option (Code β)) (form: Form): Option (Pi (Code β)) := do
  sorry

private def parse1ListPi [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Pi (Code β)) := do
  let paramForms := list.extract 0 (list.length-1)
  let params ← Util.optionMapAll paramForms (parseAnnot parse)

  let bodyForm ← list[list.length-1]?
  let body ← parse bodyForm
  pure {params := params, body := body}

private def parseListInd [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Code β) := do
  let nameForm ← list[0]?
  let name ← parseAnnot parse nameForm

  let consForm ← list.extract 1 list.length
  let cons ← Util.optionMapAll consForm (parsePi parse)

  pure (.ind {name := name, cons := cons})

private def parseBeta [Irreducible β] (parse: Form → Option (Code β)) (form: Form): Option (Code β) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.optionMapAll xs parse
      pure (.beta {cmd := cmd, args := args})
    | _ => none

partial def parseCode [Irreducible β]
  (parseAtom: String → Option β)
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

  let parseList := parseCode parseAtom


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
  ]
  ++
  -- parse beta (default case)
  [parseBeta parseList]

  Util.applyOnce parseFuncList form

end EL
