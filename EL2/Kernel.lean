import EL2.Term
import EL2.Util

namespace EL2


class Context Ctx where
  set: Ctx → String → Term × Term → Ctx
  get?: Ctx → String → Option (Term × Term)

partial def infer? [Repr Ctx] [Context Ctx] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × Term × Term) := do
  none



end EL2
