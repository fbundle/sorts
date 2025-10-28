import EL2.Term.TermUtil


namespace EL2.Term.Infer


structure ReducedTerm where
  term?: Option Term -- reduced term
  type: Term -- type of term
  level: Int -- level of term
  deriving Repr

def lift (o: Option β) (e: α): Except α β :=
  match o with
    | some v => Except.ok v
    | none => Except.error e

mutual

partial def reduce? [Repr Ctx] [Map Ctx ReducedTerm] (ctx: Ctx) (term: Term) : Except String ReducedTerm := do
  match term with
    | univ level =>
      pure {
        term? := some (univ level),
        type := univ level + 1,
        level := level + 1,-- U_1 is at level 2
        : ReducedTerm
      }
    | var name =>
      lift (Map.get? ctx name) s!"name {name} not found"

    | _ => sorry

  sorry
end


end EL2.Term.Infer
