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
  Util.applyAtmostOnce [
    parseBetaFunc parseName parseName,
    parseName >=> (some {cmd := ·, args := []}),
  ] form


def parseAnnotFunc (parseLeft: (Form → Option α)) (parseRight: (Form → Option β)): Form → Option (Annot α β) :=
  ({
    parseHead := [":"],
    parseList (list: List Form): Option (Annot α β) := do
      let leftForm ← list[0]?
      let left ← parseLeft leftForm
      let rightForm ← list[1]?
      let right ← parseRight rightForm
      pure {left := left, right := right}
  }: Form.ParseList (Annot α β)).parseForm

def parseBindingFunc(parse: (Form → Option α))  : Form → Option  (Binding α) :=
  ({
    parseHead := ["let", ":="],
    parseList (list: List Form): Option (Binding α) := do
      let nameForm ← list[0]?
      let name ← parseName nameForm
      let valueForm ← list[1]?
      let value ← parse valueForm
      pure {name := name, value := value}
  }: Form.ParseList (Binding α)).parseForm


def parseInferFunc (parse: (Form → Option α)) : Form → Option (Infer α) :=
  ({
    parseHead := ["infer", "&"],
    parseList (list: List Form): Option (Infer α) := do
      let valueForm ← list[0]?
      let value ← parse valueForm
      pure {value := value}
  }: Form.ParseList (Infer α)).parseForm


def parsePiFunc (parseAnnotType: Form → Option α) (parseBody: Form → Option β) : Form → Option (Pi α β) :=
  ({
    parseHead := ["lambda", "=>"],
    parseList (list: List Form): Option (Pi α β) := do
      let paramForms := list.dropLast
      let bodyForm ← list.getLast?

      let params ← Util.optionMapAll paramForms (parseAnnotFunc parseName parseAnnotType)
      let body ← parseBody bodyForm

      pure {params := params, body := body}
  }: Form.ParseList (Pi α β)).parseForm

def parsePatternPiAlphaBetaStringString (parseAnnotType: Form → Option α) (form: Form): Option (Pi α (Beta String String)) :=
  Util.applyAtmostOnce [
    parsePiFunc parseAnnotType parsePatternBeta,
    parsePatternBeta >=> (some {params := [], body := ·}),
  ] form


def parseIndFunc (parse: Form → Option α): Form → Option (Ind α) :=
  ({
    parseHead := ["inductive"],
    parseList (list: List Form): Option (Ind α) := do
      let nameForm ← list[0]?
      let name ← (parseAnnotFunc (parsePatternPiAlphaBetaStringString parse) parse) nameForm

      let consForm := list.extract 1 list.length
      let cons ← Util.optionMapAll consForm (parseAnnotFunc parseName (parsePatternPiAlphaBetaStringString parse))

      pure {name := name, cons := cons}
  }: Form.ParseList (Ind α)).parseForm

def parseCaseFunc (parse: Form → Option α): Form → Option (Case α) :=
  ({
    parseHead := ["case", "->"],
    parseList (list: List Form): Option (Case α) := do
      let patternForm ← list[0]?
      let pattern ← parsePatternBeta patternForm

      let valueForm ← list[1]?
      let value ← parse valueForm

      pure {pattern := pattern, value := value}
  }: Form.ParseList (Case α)).parseForm

def parseMatFunc(parse: Form → Option α) : Form → Option (Mat α) :=
  ({
    parseHead := ["match"],
    parseList (list: List Form): Option (Mat α) := do
      let condForm ← list[0]?
      let cond ← parse condForm

      let casesForm := list.extract 1 list.length
      let cases ← Util.optionMapAll casesForm (parseCaseFunc parse)

      pure {cond := cond, cases := cases}
  }: Form.ParseList (Mat α)).parseForm

partial def parseCode
  (parseAtom: String → Option β)
  (form: Form): Option (Code β) := do

  let parse := parseCode parseAtom
  let parseFuncList: List (Form → Option (Code β)) :=
  -- parse name
  [
    parseName >=> parseAtom >=> (λ x => some (Code.atom x)),
    parseName >=> (λ x => some (Code.name x)),
  ]
  ++
  -- parse basic
  [
    (parseBindingFunc parse) >=> (Code.binding ·),
    (parseInferFunc parse) >=> (Code.infer ·),
    (parsePiFunc parse parse) >=> (Code.pi ·),
    (parseIndFunc parse) >=> (Code.ind ·),
    (parseMatFunc parse) >=> (Code.mat ·),
  ]
  ++
  -- parse beta (default case)
  [parseBetaFunc parse parse >=> (Code.beta ·)]

  Util.applyAtmostOnce parseFuncList form

end EL
