import EL.Code

namespace EL


private def parseBetaFunc [Irreducible β] (parse: Form → Option (Code β)) (form: Form): Option (Beta (Code β)) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.optionMapAll xs parse
      pure {cmd := cmd, args := args}
    | _ => none

private structure ParseList α β where
  parseHead: String
  parseList (parse: Form → Option β) (list: List Form): Option α

private def ParseList.parseForm (pl: ParseList α β) (parse: Form → Option β) (form: Form) : Option α :=
  match form with
    | .list (.name x :: xs) =>
      if pl.parseHead ≠ x then
        none
      else
        pl.parseList parse xs
    | _ => none

private def ParseList.convert(pl: ParseList α β) (f: α → γ): ParseList γ β :=
  {
    parseHead := pl.parseHead,
    parseList (parse: Form → Option β) (list: List Form): Option γ := do
      let a ← pl.parseList parse list
      let c := f a
      c
  }


private def parseAnnot [Irreducible β] : ParseList (Annot (Code β)) (Code β) :=
  {
    parseHead := ":",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Annot (Code β)) := do
      let nameForm ← list[0]?
      let name ← getName nameForm
      let typeForm ← list[1]?
      let type ← parse typeForm
      pure {name := name, type := type}
  }

private def parseBinding  [Irreducible β] : ParseList (Binding (Code β)) (Code β) :=
  {
    parseHead := "let",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Binding (Code β)) := do
      let nameForm ← list[0]?
      let name ← getName nameForm
      let valueForm ← list[1]?
      let value ← parse valueForm
      pure {name := name, value := value}
  }

private def parseTypeof [Irreducible β] : ParseList (Typeof (Code β)) (Code β) :=
  {
    parseHead := "&",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Typeof (Code β)) := do
      let valueForm ← list[0]?
      let value ← parse valueForm
      pure {value := value}
  }

private def parseInh [Irreducible β] : ParseList (Inh (Code β)) (Code β) :=
  {
    parseHead := "*",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Inh (Code β)) := do
    let typeForm ← list[0]?
    let type ← parse typeForm
    pure {type := type}
  }

private def parsePi [Irreducible β] : ParseList (Pi (Code β)) (Code β) :=
  {
    parseHead := "=>",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Pi (Code β)) := do

      let paramForms := list.extract 0 (list.length-1)
      let params ← Util.optionMapAll paramForms (parseAnnot.parseForm parse)

      let bodyForm ← list[list.length-1]?
      let body ← parse bodyForm

      pure {params := params, body := body}
  }

private def parseInd [Irreducible β] : ParseList (Ind (Code β)) (Code β) :=
  {
    parseHead := "ind",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Ind (Code β)) := do
      let nameForm ← list[0]?
      let name ← parseAnnot.parseForm parse nameForm

      let consForm := list.extract 1 list.length
      let cons ← Util.optionMapAll consForm (parsePi.parseForm parse)

      pure {name := name, cons := cons}
  }

private def parseCase [Irreducible β] : ParseList (Case (Code β)) (Code β) :=
  {
    parseHead := "->",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Case (Code β)) := do
      let condForm ← list[0]?
      let cond ← parseBetaFunc parse condForm

      let valueForm ← list[1]?
      let value ← parse valueForm
      pure {cond := cond, value := value}
  }

private def parseMat [Irreducible β] : ParseList (Mat (Code β)) (Code β) :=
  {
    parseHead := "match",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Mat (Code β)) := do
      let compForm ← list[0]?
      let comp ← parse compForm

      let casesForm := list.extract 1 list.length
      let cases ← Util.optionMapAll casesForm (parseCase.parseForm parse)

      pure {comp := comp, cases := cases}
  }

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
    (parseBinding.convert (λ x => (Code.binding x))).parseForm parse,
    (parseTypeof.convert (λ x => (Code.typeof x))).parseForm parse,
    (parseInh.convert (λ x => (Code.inh x))).parseForm parse,


    (parsePi.convert (λ x => (Code.pi x))).parseForm parse,
    (parseInd.convert (λ x => (Code.ind x))).parseForm parse,
    (parseMat.convert (λ x => (Code.mat x))).parseForm parse,
  ]

  -- parse beta (default case)
  ++
  [λ form => ((λ x => Code.beta x) <$> (parseBetaFunc parse form))]

  Util.applyOnce parseFuncList form

end EL
