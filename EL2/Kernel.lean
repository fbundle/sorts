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

def isSubType (type1: Term) (type2: Term): Bool := type1 == type2

mutual
partial def reduceSequential? [Repr F] [Frame F] (frame: F) (terms: List Term): Option (F × List (InferedTerm)) :=
  -- reduce sequentially so that dependent type is captured properly
  Util.statefulMap? terms frame (λ frame term => do
    let (frame, iterm) ← reduce? frame term
    pure (frame, iterm)
  )

partial def reduceParamsType? [Repr F] [Frame F] (frame: F) (params: List (Ann Term)): Option (F × List (Ann InferedTerm)) :=
  Util.statefulMap? params frame (λ frame param => do
    let (frame, itype) ← reduce? frame param.type
    pure (frame, {
      name := param.name,
      type := itype,
    })
  )

partial def bindParamsWithArgs [Repr F] [Frame F] (frame: F) (iParams: List (Ann InferedTerm)) (iArgs: List InferedTerm): Option F := do
  if iParams.length = 0 ∧ iArgs.length = 0 then
    pure frame
  else
    let iParam ← iParams.head?
    let iArg ← iArgs.head?

    if iArg.type != iParam.type.term then
      none
    else
      let frame := Frame.set frame iParam.name iArg
      let iParams := iParams.extract 1
      let iArgs := iArgs.extract 1
      bindParamsWithArgs frame iParams iArgs


partial def reduce? [Repr F] [Frame F] (oldFrame: F) (term: Term): Option (F × InferedTerm) := do
  let frame := oldFrame -- for update
  dbg_trace s!"# 1 {term} with frame \n{repr frame}"
  match term with
    | univ level =>
      pure (oldFrame, {
        term := univ level,
        type := univ level+1,
        level := level + 1, -- univ 1 is at level 2, Nat: univ 1, then Nat is at level 1
      })

    | var name =>
      let iterm ← Frame.get? frame name
      pure (oldFrame, iterm)

    | inh x =>
      let (_, iType) ← reduce? frame x.type
      let (_, iArgs) ← reduceSequential? frame x.args
      pure (oldFrame, {
        term := inh {
          type := iType.term,
          cons := x.cons,
          args := iArgs.map (λ iterm => iterm.term),
        },
        type := iType.term,
        level := iType.level - 1,
      })

    | typ x =>
      let (_, iValue) ← reduce? frame x.value
      let (_, iType) ← reduce? frame iValue.type
      pure (oldFrame, iType)

    | lst x =>
      let (initFrame, _) ← reduceSequential? frame x.init
      reduce? initFrame x.last

    | bind x =>
      let (_, iValue) ← reduce? frame x.value
      pure (Frame.set oldFrame x.name iValue, iValue)

    | lam x =>
      let (paramsFrame, iParams) ← reduceParamsType? frame x.params
      let (_, iType) ← reduce? paramsFrame x.type


      let rParams := iParams.map (λ iparam => {name := iparam.name, type := iparam.type.term : Ann Term})
      let rLevel := (iParams.map (λ iparam => iparam.type.level)).foldl max iType.level

      pure (oldFrame, {
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
      })

    | app x =>
      let (_, iCmd) ← reduce? frame x.cmd
      let (_, iArgs) ← reduceSequential? frame x.args
      let lamCmd ← isLam? iCmd.term

      let (paramsFrame, iParams) ← reduceParamsType? frame lamCmd.params
      let (_, iType) ← reduce? paramsFrame lamCmd.type

      let argsFrame ← bindParamsWithArgs frame iParams iArgs

      let (_, output) ← reduce? argsFrame lamCmd.body

      if output.type != iType.term then
        none
      else
        pure (oldFrame, output)
    | mat x => none

end


end EL2
