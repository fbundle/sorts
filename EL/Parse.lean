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

private def parseBeta [Irreducible β] (parse: Form → Option (Code β)) (form: Form): Option (Beta (Code β)) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.optionMapAll xs parse
      pure {cmd := cmd, args := args}
    | _ => none

private def parseListAnnot [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Annot (Code β)) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure {name := name, type := type}

private def parseListBinding [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Binding (Code β)) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let valueForm ← list[1]?
  let value ← parse valueForm
  pure {name := name, value := value}

private partial def parseListTypeof [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Typeof (Code β)) := do
  let valueForm ← list[0]?
  let value ← parse valueForm
  pure {value := value}

private def parseListInh [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Inh (Code β)) := do
  let typeForm ← list[0]?
  let type ← parse typeForm
  pure {type := type}

private def parseListPi [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Pi (Code β)) := do
  let parseAnnot := parseListToParseForm (parseListAnnot parse) ":"

  let paramForms := list.extract 0 (list.length-1)
  let params ← Util.optionMapAll paramForms parseAnnot

  let bodyForm ← list[list.length-1]?
  let body ← parse bodyForm

  pure {params := params, body := body}

private def parseListInd [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Ind (Code β)) := do
  let parseAnnot := parseListToParseForm (parseListAnnot parse) ":"
  let parsePi := parseListToParseForm (parseListPi parse) "=>"

  let nameForm ← list[0]?
  let name ← parseAnnot nameForm

  let consForm := list.extract 1 list.length
  let cons ← Util.optionMapAll consForm parsePi

  pure {name := name, cons := cons}

private def parseListCase [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Case (Code β)) := do
  let condForm ← list[0]?
  let cond ← parseBeta parse condForm

  let valueForm ← list[1]?
  let value ← parse valueForm
  pure {cond := cond, value := value}

private def parseListMat [Irreducible β] (parse: Form → Option (Code β)) (list: List Form): Option (Mat (Code β)) := do
  let parseCase := parseListToParseForm (parseListCase parse) "->"

  let compForm ← list[0]?
  let comp ← parse compForm

  let casesForm := list.extract 1 list.length
  let cases ← Util.optionMapAll casesForm parseCase

  pure {comp := comp, cases := cases}


partial def parseCode [Irreducible β]
  (parseAtom: String → Option β)
  (form: Form): Option (Code β) := do

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

  let parse := parseCode parseAtom
  let parseFuncList: List (Form → Option (Code β)) :=
  -- parse name
  [
    parseAtomFunc,
    parseNameFunc,
  ]
  ++
  -- parse basic
  [
    parseFormToParseForm (parseListToParseForm (parseListBinding parse) "let") (λ x => (.binding x)),
    parseFormToParseForm (parseListToParseForm (parseListTypeof parse) "&") (λ x => (.typeof x)),
    parseFormToParseForm (parseListToParseForm (parseListInh parse) "*") (λ x => (.inh x)),
    parseFormToParseForm (parseListToParseForm (parseListPi parse) "=>") (λ x => (.pi x)),
    parseFormToParseForm (parseListToParseForm (parseListInd parse) "ind") (λ x => (.ind x)),
    parseFormToParseForm (parseListToParseForm (parseListMat parse) "match") (λ x => (.mat x)),
  ]
  ++
  -- parse beta (default case)
  [parseFormToParseForm (parseBeta parse) (λ x => (.beta x))]

  Util.applyOnce parseFuncList form

end EL
