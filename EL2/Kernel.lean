import EL2.Term
import EL2.Util

namespace EL2


class Context Ctx where
  set: Ctx → String → Term × Term × Int → Ctx
  get?: Ctx → String → Option (Term × Term × Int)

structure InferedTerm where
  term: Term
  type: Term
  level: Int

def dummyName (i: Int): String := s!"dummy_{i}"

partial def infer? [Repr Ctx] [Context Ctx] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × InferedTerm) := do
  let isLam? (term: Term): Option (Lam Term) :=
    match term with
      | lam l => some l
      | _ => none
  none



end EL2
