import EL2.Term
import EL2.Util

namespace EL2


class Context Ctx where
  set: Ctx → String → Term × Term → Ctx
  get?: Ctx → String → Option (Term × Term)

partial def infer? [Repr Ctx] [Context Ctx] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × Term × Term) := do
  -- return (ctx, term, type)
  match term with
    | .univ level =>
      pure (ctx, term, Term.univ (level+1))
    | .var name =>
      let (term, type) ← Context.get? ctx name
      pure (ctx, term, type)
    | .inh type _ _ =>
      pure (ctx, term, type)
    | .infer value =>
      let (_, valueTerm, valueType) ← infer? reduce ctx value
      pure (ctx, valueTerm, valueType)
    | .list init last =>
      none
    | _ => none



end EL2
