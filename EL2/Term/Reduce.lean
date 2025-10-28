import EL2.Term.Term
import EL2.Term.TermUtil
import EL2.Term.Util
import EL2.Term.Print

namespace EL2.Term.Infer


structure ReducedTerm where
  term?: Option Term -- reduced term
  type: Term -- type of term
  level: Int -- level of term
  deriving Repr

mutual

partial def reduce? [Repr Ctx] [Map Ctx ReducedTerm] (ctx: Ctx) (term: Term) : Option ReducedTerm := do
  let o : Option ReducedTerm := do
    match term with
      | univ level =>
        pure {
          term? := some (univ level),
          type := univ level + 1,
          level := level + 1,-- U_1 is at level 2
          : ReducedTerm
        }
      | var name =>
        Map.get? ctx name

      | _ => none

    sorry
  match o with
    | none =>
      dbg_trace s!"[DBG_TRACE] failed at {term} with ctx {repr ctx}"
      none
    | some v =>
      pure v
end


end EL2.Term.Infer
