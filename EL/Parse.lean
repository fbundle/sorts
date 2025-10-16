import EL.Form
import EL.Code

namespace EL

abbrev getName := Form.getName
abbrev getList := Form.getList
abbrev Form := Form.Form

def parseName (form: Form): Option String :=
  match form with
    | .name n => some n
    | _ => none

def parseBetaOfSomething (parse: Form → Option α) (form: Form): Option (Beta α) := do
  match form with
    | .list (x :: xs) =>
      let cmd ← parse x
      let args ← Util.optionMapAll xs parse
      pure {cmd := cmd, args := args}
    | _ => none

def parseBetaOfStringFunc : (Form → Option String) → Form → Option (Beta String) := parseBetaOfSomething

def parseBetaFunc: (Form → Option (Code β)) → Form → Option (Beta (Code β)) := parseBetaOfSomething

structure ParseList α β where
  parseHead: String
  parseList (parse: Form → Option α) (list: List Form): Option β

def ParseList.parseForm (pl: ParseList α β) (parse: Form → Option α) (form: Form) : Option β :=
  match form with
    | .list (.name x :: xs) =>
      if pl.parseHead ≠ x then
        none
      else
        pl.parseList parse xs
    | _ => none

def ParseList.convert(pl: ParseList α β) (f: β → γ): ParseList α γ :=
  {
    parseHead := pl.parseHead,
    parseList (parse: Form → Option α) (list: List Form): Option γ := do
      let b ← pl.parseList parse list
      let c := f b
      c
  }

def parseAnnotOfSomething α : ParseList α (Annot α) :=
  {
    parseHead := ":",
    parseList (parse: Form → Option α) (list: List Form): Option (Annot α) := do
      let nameForm ← list[0]?
      let name ← getName nameForm
      let typeForm ← list[1]?
      let type ← parse typeForm
      pure {name := name, type := type}
  }

def parseAnnotOfPi : ParseList (Pi (Code β)) (Annot (Pi (Code β))) := parseAnnotOfSomething (Pi (Code β))

def parseAnnot {β} := parseAnnotOfSomething (Code β)




def parseBinding  : ParseList (Code β) (Binding (Code β)) :=
  {
    parseHead := "let",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Binding (Code β)) := do
      let nameForm ← list[0]?
      let name ← getName nameForm
      let valueForm ← list[1]?
      let value ← parse valueForm
      pure {name := name, value := value}
  }

def parseInfer : ParseList (Code β) (Infer (Code β)) :=
  {
    parseHead := "type",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Infer (Code β)) := do
      let valueForm ← list[0]?
      let value ← parse valueForm
      pure {value := value}
  }

def parsePi : ParseList (Code β) (Pi (Code β)) :=
  {
    parseHead := "lambda",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Pi (Code β)) := do
      let paramForms := list.extract 0 (list.length-1)
      let params ← Util.optionMapAll paramForms (parseAnnot.parseForm parse)

      let bodyForm ← list[list.length-1]?
      let body ← parse bodyForm

      pure {params := params, body := body}
  }


def parseInd : ParseList (Code β) (Ind (Code β)) :=
  {
    parseHead := "inductive",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Ind (Code β)) := do
      let nameForm ← list[0]?
      let name ← parseAnnot.parseForm parse nameForm

      let consForm := list.extract 1 list.length
      let cons ← Util.optionMapAll consForm (parseAnnotOfPi.parseForm (parsePi.parseForm parse))

      pure {name := name, cons := cons}
  }

def parseCase : ParseList (Code β) (Case (Code β)) :=
  {
    parseHead := "case",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Case (Code β)) := do
      let patternForm ← list[0]?
      let pattern ← parseBetaOfStringFunc parseName patternForm

      let valueForm ← list[1]?
      let value ← parse valueForm

      pure {pattern := pattern, value := value}
  }

def parseMat : ParseList (Code β) (Mat (Code β)) :=
  {
    parseHead := "match",
    parseList (parse: Form → Option (Code β)) (list: List Form): Option (Mat (Code β)) := do
      let condForm ← list[0]?
      let cond ← parse condForm

      let casesForm := list.extract 1 list.length
      let cases ← Util.optionMapAll casesForm (parseCase.parseForm parse)

      pure {cond := cond, cases := cases}
  }

partial def parseCode
  (parseAtom: String → Option β)
  (form: Form): Option (Code β) := do

  let parseAtomFunc (form: Form): Option (Code β) := do
    let n ← getName form
    let a ← parseAtom n
    pure (.atom a)

  let parseNameFunc (form: Form): Option (Code β) := do
    let n ← parseName form
    pure (.name n)


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
    (parseInfer.convert (λ x => (Code.infer x))).parseForm parse,
    (parsePi.convert (λ x => (Code.pi x))).parseForm parse,
    (parseInd.convert (λ x => (Code.ind x))).parseForm parse,
    (parseMat.convert (λ x => (Code.mat x))).parseForm parse,
  ]
  ++
  -- parse beta (default case)
  [λ form => ((λ x => Code.beta x) <$> (parseBetaFunc parse form))]

  Util.applyOnce parseFuncList form

end EL
