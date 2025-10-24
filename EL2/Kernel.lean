import EL2.Term
import EL2.TermNot
import EL2.Util
import EL2.Print -- for debugging

namespace EL2

class Frame F α where
  set: F → String → α → F
  get?: F → String → Option α

structure InferedTerm where
  term: Term
  type: Term
  deriving Repr

def isLam? (term: Term): Option (Lam Term) :=
  match term with
    | lam l => some l
    | _ => none

def isInh? (term: Term): Option (Inh Term) :=
  match term with
    | inh i => some i
    | _ => none

def isSubType (type1: Term) (type2: Term): Bool :=
  type1 == type2

partial def matchCases? (inhCond: Inh Term) (cases: List (Case Term)): Option (Bnd Term) := do
  -- TODO make dummy name match every -> return a list of terms
  -- TODO make reduce returns a list of terms in typecheck mode
  -- TODO use body instead of type in Lam
  let headCase ← cases.head?
  if inhCond.cons = headCase.patCmd ∧ inhCond.args.length = headCase.patArgs.length then
    let init: List (Bind Term) := (List.zip headCase.patArgs inhCond.args).map (λ (name, value) => {
      name := name,
      value := value,
    })
    pure {
      init := init,
      last := headCase.value,
    }
  else
    let cases := cases.extract 1
    matchCases? inhCond cases


mutual
partial def reduceMany? [Repr F] [Frame F InferedTerm] (frame: F) (terms: List Term): Option (List (InferedTerm)) :=
  Util.optionMap? terms (reduceTerm? frame)

partial def reduceUniv? [Repr F] [Frame F InferedTerm] (frame: F) (level: Int): Option InferedTerm := do
  dbg_trace s!"reduce univ {univ level}"
  pure {
    term := univ level,
    type := univ level+1,
  }

partial def reduceVar? [Repr F] [Frame F InferedTerm] (frame: F) (name: String): Option InferedTerm := do
  dbg_trace s!"reduce var {var name}"
  let iterm ← Frame.get? frame name
  pure iterm

partial def reduceInh? [Repr F] [Frame F InferedTerm] (frame: F) (x: Inh Term): Option InferedTerm := do
  dbg_trace s!"reduce inh {inh x}"
  let iType ← reduceTerm? frame x.type
  let iArgs ← reduceMany? frame x.args
  pure {
    term := inh {
      type := iType.term,
      cons := x.cons,
      args := iArgs.map (λ iterm => iterm.term),
    },
    type := iType.term,
  }

partial def reduceTyp? [Repr F] [Frame F InferedTerm] (frame: F) (x: Typ Term): Option InferedTerm := do
  dbg_trace s!"reduce typ {typ x}"
  let iValue ← reduceTerm? frame x.value
  let iType ← reduceTerm? frame iValue.type
  pure iType

partial def reduceBnd? [Repr F] [Frame F InferedTerm] (frame: F) (x: Bnd Term): Option InferedTerm := do
  dbg_trace s!"reduce bnd {bnd x}"
  let (frame, _) ← Util.statefulMap? x.init frame (λ frame {name, value} => do
    let iValue ← reduceTerm? frame value
    let frame := Frame.set frame name iValue
    some (frame, iValue)
  )
  reduceTerm? frame x.last

partial def reduceLam? [Repr F] [Frame F InferedTerm] (frame: F) (x: Lam Term): Option InferedTerm := do
  dbg_trace s!"reduce lam {lam x}"
  pure {
    term := lam x,
    type := lam {
      params := x.params,
      body := typ {
        value := x.body,
      }
    },
  }

partial def bindParamsWithArgs? [Repr F] [Frame F InferedTerm] (frame: F) (params: List (Ann Term)) (args: List Term): Option F := do
  if params.length = 0 ∧ args.length = 0 then
    pure frame
  else
    let param ← params.head?
    let arg ← args.head?

    let iParamType ← reduceTerm? frame param.type
    let iArg ← reduceTerm? frame arg

    if ¬ isSubType iArg.type iParamType.term then
      none
    else
      let frame := Frame.set frame param.name iArg
      let params := params.extract 1
      let args := args.extract 1
      bindParamsWithArgs? frame params args

partial def reduceApp? [Repr F] [Frame F InferedTerm] (frame: F) (x: App Term): Option InferedTerm := do
  dbg_trace s!"reduce app {app x}"
  let iCmd ← reduceTerm? frame x.cmd
  let lamCmd ← isLam? iCmd.term
  let argFrame ← bindParamsWithArgs? frame lamCmd.params x.args
  reduceTerm? argFrame lamCmd.body

partial def reduceMat? [Repr F] [Frame F InferedTerm] (frame: F) (x: Mat Term): Option InferedTerm := do
  dbg_trace s!"reduce mat {mat x}"
  let iCond ← reduceTerm? frame x.cond
  let inhCond ← isInh? iCond.term

  let terms ← matchCases? inhCond x.cases

  reduceTerm? frame (bnd terms)

partial def reduceTerm? [Repr F] [Frame F InferedTerm] (frame: F) (term: Term): Option InferedTerm := do
  match term with
    | univ level => reduceUniv? frame level
    | var name => reduceVar? frame name
    | inh x => reduceInh? frame x
    | typ x => reduceTyp? frame x
    | bnd x => reduceBnd? frame x
    | lam x => reduceLam? frame x
    | app x => reduceApp? frame x
    | mat x => reduceMat? frame x

end

partial def fill? [Repr F] [Frame F InferedTerm] (frame: F) (term: Term): Option (F × Term) :=
  -- fill in all the holes
  -- e.g return type
  none


-- type check


class ListFrame F α where
  set: F → String → α → F
  get: F → String → List α

structure InferedType where
  type: Term
  level: Int


mutual



partial def inferType? [Repr F] [ListFrame F (List InferedTerm)] (frame: F) (term: Term): List InferedTerm :=
  -- TODO decompose inferType and reduceTerm
  []

end

end EL2
