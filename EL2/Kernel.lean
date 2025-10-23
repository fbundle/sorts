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

mutual
partial def reduceSequential? [Repr F] [Frame F] (frame: F) (termList: List Term): Option (F × List (InferedTerm)) :=
  -- reduce sequentially so that dependent type is captured properly
  Util.statefulMap? termList frame (λ frame term => do
    let (frame, iterm) ← reduce? frame term
    pure (frame, iterm)
  )

partial def reduceParamsType? [Repr F] [Frame F] (frame: F) (paramList: List (Ann Term)): Option (F × List (Ann InferedTerm)) :=
  Util.statefulMap? paramList frame (λ frame param => do
    let (frame, itype) ← reduce? frame param.type
    pure (frame, {
      name := param.name,
      type := itype,
    })
  )


partial def reduce? [Repr F] [Frame F] (oldFrame: F) (term: Term): Option (F × InferedTerm) := do
  let frame := oldFrame -- for update
  dbg_trace s!"# reducing {term} with frame \n{repr frame}"
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


      none
      -- TODO - for level 0, do reduce if only specified for level > 1 reduce
    | mat x => none

end


end EL2
