import EL2.Term
import EL2.Util
import EL2.Print -- for debugging

namespace EL2

structure InferedTerm where
  term: Term
  type: Term
  level: Int
  deriving Repr

class Frame F where
  set: F → String → InferedTerm → F
  get?: F → String → Option InferedTerm


def dummyName (i: Int): String := s!"dummy_{i}"

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

partial def bindParamsWithArgs [Repr F] [Frame F] (frame: F) (iParams: List (Ann InferedTerm)) (iArgs: List InferedTerm): Option F := do
  if iParams.length = 0 ∧ iArgs.length = 0 then
    pure frame
  else
    let iParam ← iParams.head?
    let iArg ← iArgs.head?

    if ¬ isSubType iArg.type iParam.type.term then
      none
    else
      let frame := Frame.set frame iParam.name iArg
      let iParams := iParams.extract 1
      let iArgs := iArgs.extract 1
      bindParamsWithArgs frame iParams iArgs

partial def matchCases (inhCond: Inh Term) (cases: List (Case Term)): Option (Lst Term) := do
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
    matchCases inhCond cases


mutual
partial def reduceMany? [Repr F] [Frame F] (frame: F) (terms: List Term): Option (List (InferedTerm)) :=
  Util.optionMap? terms (reduce? frame)

partial def reduceLst? [Repr F] [Frame F] (frame: F) (l: Lst Term): Option InferedTerm := do
  let (frame, _) ← Util.statefulMap? l.init frame (λ frame {name, value} => do
    let iValue ← reduce? frame value
    let frame := Frame.set frame name iValue
    some (frame, iValue)
  )
  reduce? frame l.last

partial def reduceParams? [Repr F] [Frame F] (frame: F) (params: List (Ann Term)): Option (F × List (Ann InferedTerm)) := do
  -- reduce and bind params with dummy values
  -- reuse frame so that dependent type (Pi, Sigma) is captured
  let counter: F × Int := (frame, 0)
  let ((frame, count), iParams) ← Util.statefulMap? params counter ((λ (frame, count) param => do
    let itype ← reduce? frame param.type
    let dummyTerm := inh {
      type := itype.term,
      cons := dummyName count,
      args := [],
    }
    let dummyITerm ← reduce? frame dummyTerm
    let frame := Frame.set frame param.name dummyITerm
    pure ((frame, count + 1), {
      name := param.name,
      type := itype,
    })
  ): F × Int → Ann Term → Option ((F × Int) × Ann InferedTerm))

  pure (frame, iParams)

partial def reduceCases? [Repr F] [Frame F] (frame: F) (cases: List (Case Term)): Option (F × List (Case InferedTerm)) := do
  let oldFrame := frame
  let iCases ← Util.optionMap? cases ((λ {patCmd, patArgs, value} => do
    let iValue ← reduce? frame value
    pure {patCmd := patCmd, patArgs := patArgs, value := iValue}
  ): Case Term → Option (Case InferedTerm))

  pure (oldFrame, iCases)

partial def reduce? [Repr F] [Frame F] (frame: F) (term: Term): Option InferedTerm := do
  dbg_trace s!"#1 {term}"
  match term with
    | univ level =>
      pure {
        term := univ level,
        type := univ level+1,
        level := level + 1, -- univ 1 is at level 2, Nat: univ 1, then Nat is at level 1
      }

    | var name =>
      let iterm ← Frame.get? frame name
      pure iterm

    | inh x =>
      let iType ← reduce? frame x.type
      let iArgs ← reduceMany? frame x.args
      pure {
        term := inh {
          type := iType.term,
          cons := x.cons,
          args := iArgs.map (λ iterm => iterm.term),
        },
        type := iType.term,
        level := iType.level - 1,
      }

    | typ x =>
      let iValue ← reduce? frame x.value
      let iType ← reduce? frame iValue.type
      pure iType

    | lst x =>
      let (frame, _) ← Util.statefulMap? x.init frame (λ frame {name, value} => do
        let iValue ← reduce? frame value
        let frame := Frame.set frame name iValue
        some (frame, iValue)
      )
      reduce? frame x.last

    | lam x =>
      let (paramsFrame, iParams) ← reduceParams? frame x.params
      let iType ← reduce? paramsFrame x.type


      let rParams := iParams.map (λ iparam => {name := iparam.name, type := iparam.type.term : Ann Term})
      let rLevel := (iParams.map (λ iparam => iparam.type.level)).foldl max iType.level

      pure {
        term := lam {
          params := rParams,
          type := iType.term,
          body := x.body,
        },
        type := lam {
          params := rParams,
          type := iType.type,
          body := typ {
            value := x.body,
          }
        },
        level := rLevel,
      }

    | app x =>
      let iCmd ← reduce? frame x.cmd
      let iArgs ← reduceMany? frame x.args
      let lamCmd ← isLam? iCmd.term

      let (paramsFrame, iParams) ← reduceParams? frame lamCmd.params
      let iType ← reduce? paramsFrame lamCmd.type

      let argsFrame ← bindParamsWithArgs frame iParams iArgs

      let output ← reduce? argsFrame lamCmd.body

      if ¬ isSubType output.type iType.term then
        none
      else
        pure output
    | mat x =>
      let iCond ← reduce? frame x.cond
      let inhCond ← isInh? iCond.term

      let terms ← matchCases inhCond x.cases

      reduce? frame (lst terms)
end

partial def fill? [Repr F] [Frame F] (frame: F) (term: Term): Option (F × Term) :=
  -- fill in all the holes
  -- e.g return type
  none


end EL2
