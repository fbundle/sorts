import EL.Form
import EL.Code
import EL.Util

namespace EL

abbrev Form := Form.Form

def parseName (form: Form): Option String :=
  match form with
    | .name n => some n
    | _ => none

partial def parseBetaFunc (parseCmd: Form → Option α) (parseArg: Form → Option β) (form: Form): Option (Beta α β) :=
  match form with
    | .list (cmdForm :: argForms) => do
      let cmd ← parseCmd cmdForm
      let args ← Util.optionMapAll argForms parseArg
      pure {cmd := cmd, args := args}
    | _ => none

partial def parsePatternBeta (form: Form): Option (Beta String String) :=
  match form with
    | .name n => some {cmd := n, args := []}
    | .list _ => parseBetaFunc parseName parseName form

structure ParseList γ where
  parseHead: List String
  parseList (list: List Form): Option γ

def ParseList.parseForm (pl: ParseList γ) (form: Form) : Option γ :=
  match form with
    | .list (.name x :: xs) =>
      if pl.parseHead.contains x then
        pl.parseList xs
      else
        none
    | _ => none

def ParseList.convert (pl: ParseList γ) (f: γ → δ): ParseList δ :=
  {
    parseHead := pl.parseHead,
    parseList (list: List Form): Option δ := do
      let b ← pl.parseList list
      let c := f b
      c
  }

def parseAnnot (parseLeft: (Form → Option α)) (parseRight: (Form → Option β)): ParseList (Annot α β) :=
  {
    parseHead := [":"],
    parseList (list: List Form): Option (Annot α β) := do
      let leftForm ← list[0]?
      let left ← parseLeft leftForm
      let rightForm ← list[1]?
      let right ← parseRight rightForm
      pure {left := left, right := right}
  }

def parseBinding(parse: (Form → Option α))  : ParseList (Binding α) :=
  {
    parseHead := ["let", ":="],
    parseList (list: List Form): Option (Binding α) := do
      let nameForm ← list[0]?
      let name ← parseName nameForm
      let valueForm ← list[1]?
      let value ← parse valueForm
      pure {name := name, value := value}
  }

def parseInfer(parse: (Form → Option α)) : ParseList (Infer α) :=
  {
    parseHead := ["infer", "&"],
    parseList (list: List Form): Option (Infer α) := do
      let valueForm ← list[0]?
      let value ← parse valueForm
      pure {value := value}
  }


def parsePi (parseAnnotType: Form → Option α) (parseBody: Form → Option β) : ParseList (Pi α β) :=
  {
    parseHead := ["lambda", "=>"],
    parseList (list: List Form): Option (Pi α β) := do
      let paramForms := list.dropLast
      let bodyForm ← list.getLast?

      let params ← Util.optionMapAll paramForms (parseAnnot parseName parseAnnotType).parseForm
      let body ← parseBody bodyForm

      pure {params := params, body := body}
  }

def parsePatternPiAlphaBetaStringString (parseAnnotType: Form → Option α) (form: Form): Option (Pi α (Beta String String)) :=
  match form with
    | .name n =>
      some {params := [], body := {cmd := n, args := []}}
    | .list _ =>
      (parsePi parseAnnotType parsePatternBeta).parseForm form

def parseInd (parse: Form → Option α): ParseList (Ind α) :=
  {
    parseHead := ["inductive"],
    parseList (list: List Form): Option (Ind α) := do
      let nameForm ← list[0]?
      let name ← (parseAnnot (parsePatternPiAlphaBetaStringString parse) parse).parseForm nameForm

      let consForm := list.extract 1 list.length
      let cons ← Util.optionMapAll consForm (parseAnnot parseName (parsePatternPiAlphaBetaStringString parse)).parseForm

      pure {name := name, cons := cons}
  }

def parseCase (parse: Form → Option α): ParseList (Case α) :=
  {
    parseHead := ["case", "->"],
    parseList (list: List Form): Option (Case α) := do
      let patternForm ← list[0]?
      let pattern ← parsePatternBeta patternForm

      let valueForm ← list[1]?
      let value ← parse valueForm

      pure {pattern := pattern, value := value}
  }

def parseMat(parse: Form → Option α) : ParseList (Mat α) :=
  {
    parseHead := ["match"],
    parseList (list: List Form): Option (Mat α) := do
      let condForm ← list[0]?
      let cond ← parse condForm

      let casesForm := list.extract 1 list.length
      let cases ← Util.optionMapAll casesForm (parseCase parse).parseForm

      pure {cond := cond, cases := cases}
  }

partial def parseCode
  (parseAtom: String → Option β)
  (form: Form): Option (Code β) := do

  let parseAtomFunc (form: Form): Option (Code β) := do
    let n ← parseName form
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
    ((parseBinding parse).convert (Code.binding ·)).parseForm,
    ((parseInfer parse).convert (Code.infer ·)).parseForm,
    ((parsePi parse parse).convert (Code.pi ·)).parseForm,
    ((parseInd parse).convert (Code.ind ·)).parseForm,
    ((parseMat parse).convert (Code.mat ·)).parseForm,
  ]
  ++
  -- parse beta (default case)
  [λ (form: Form) => do
    let b ← parseBetaFunc parse parse form
    pure (Code.beta b)
  ]


  Util.applyAtmostOnce parseFuncList form

end EL
