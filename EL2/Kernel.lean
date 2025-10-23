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
  
mutual
partial def reduceList? [Repr F] [Frame F] (frame: F) (terms: List Term): Option (F × List (InferedTerm)) :=
  -- reuse frame so that dependent type is capture
  Util.statefulMap? terms frame (λ frame term => do
    let (frame, iterm) ← reduce? frame term
    pure (frame, iterm)
  )

partial def reduceParams? [Repr F] [Frame F] (frame: F) (params: List (Ann Term)): Option (F × List (Ann InferedTerm)) := do
  let counter: F × Int := (frame, 0)
  let ((frame, count), iParams) ← Util.statefulMap? params counter ((λ (frame, count) param => do
    let (frame, itype) ← reduce? frame param.type
    let dummyTerm := inh {
      type := itype.term,
      cons := dummyName count,
      args := [],
    }
    let (frame, dummyITerm) ← reduce? frame dummyTerm
    let frame := Frame.set frame param.name dummyITerm
    pure ((frame, count + 1), {
      name := param.name,
      type := itype,
    })
  ): F × Int → Ann Term → Option ((F × Int) × Ann InferedTerm))

  pure (frame, iParams)


partial def reduce? [Repr F] [Frame F] (oldFrame: F) (term: Term): Option (F × InferedTerm) := do
  let frame := oldFrame -- for update
  dbg_trace s!"#1 {term}"
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
      let (_, iArgs) ← reduceList? frame x.args
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
      let (initFrame, _) ← reduceList? frame x.init
      reduce? initFrame x.last

    | bind x =>
      let (_, iValue) ← reduce? frame x.value
      pure (Frame.set oldFrame x.name iValue, iValue)

    | lam x =>
      let (paramsFrame, iParams) ← reduceParams? frame x.params
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
      let (_, iArgs) ← reduceList? frame x.args
      let lamCmd ← isLam? iCmd.term

      let (paramsFrame, iParams) ← reduceParams? frame lamCmd.params
      let (_, iType) ← reduce? paramsFrame lamCmd.type

      let argsFrame ← bindParamsWithArgs frame iParams iArgs

      let (_, output) ← reduce? argsFrame lamCmd.body

      if ¬ isSubType output.type iType.term then
        none
      else
        pure (oldFrame, output)
    | mat x =>
      let (_, iCond) ← reduce? frame x.cond

      none

end


end EL2
