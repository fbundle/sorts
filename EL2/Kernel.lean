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


partial def reduce? [Repr F] [Frame F] (oldFrame: F) (term: Term): Option (F × InferedTerm) := do
  let inferMany? (frame: F) (termList: List Term): Option (F × List (InferedTerm)) :=
    Util.statefulMap? termList frame (λ frame term => do
      let (frame, iterm) ← reduce? frame term
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
      let (_, iType) ← reduce? frame x.type
      let (_, iArgs) ← inferMany? frame x.args
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
      let (initFrame, _) ← inferMany? frame x.init
      reduce? initFrame x.last

    | bind x =>
      let (_, iValue) ← reduce? frame x.value
      pure (Frame.set oldFrame x.name iValue, iValue)

    | lam x => none
    | app x => none
      -- TODO - for level 0, do reduce if only specified for level > 1 reduce
    | mat x => none



end EL2
