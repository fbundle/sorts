import EL2.Term
import EL2.Util

namespace EL2

structure InferedTerm where
  term: Term
  type: Term
  level: Int

class Frame F where
  set: F → String → InferedTerm → F
  get?: F → String → Option InferedTerm


def dummyName (i: Int): String := s!"dummy_{i}"


partial def infer? [Repr F] [Frame F] (oldFrame: F) (term: Term) (reduce: Bool := false): Option (F × InferedTerm) := do
  let inferMany? (frame: F) (termList: List Term) (reduce: Bool): Option (F × List (InferedTerm)) :=
    Util.statefulMap? termList frame (λ frame term => do
      let (frame, iterm) ← infer? frame term reduce
      pure (frame, iterm)
    )


  let frame := oldFrame -- for update

  let isLam? (term: Term): Option (Lam Term) :=
    match term with
      | lam l => some l
      | _ => none

  match term with
    | univ level =>
      pure (oldFrame, {
        term := univ level,
        type := univ level+1,
        level := level,
      })

    | var name =>
      let iterm ← Frame.get? frame name
      pure (oldFrame, iterm)

    | inh x =>
      let (_, iType) ← infer? frame x.type reduce
      let (_, iArgs) ← inferMany? frame x.args reduce
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
      let (_, iValue) ← infer? frame x.value reduce
      let (_, iType) ← infer? frame iValue.type reduce
      pure (oldFrame, iType)

    | lst x =>
      let (initFrame, _) ← inferMany? frame x.init reduce
      infer? initFrame x.last reduce

    | bind x =>
      let (_, iValue) ← infer? frame x.value reduce
      pure (Frame.set oldFrame x.name iValue, iValue)

    | lam x => none
    | app x => none
    | mat x => none



end EL2
