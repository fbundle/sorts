import EL2.Term

namespace EL2


partial def Term.infer [Irreducible β] [Context Ctx (Term β)] (c: Term β) (ctx: Ctx) : Option (Term β × Ctx) := do
  -- infer: turn everything to type then normalize
  match c with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (.atom p, ctx)
    | .var n =>
      let c : Term β ← Context.get? ctx n
      c.infer ctx
    | _ => sorry

partial def Term.normalize [Irreducible β] [Context Ctx (Term β)] (c: Term β) (ctx: Ctx): Option (Term β × Ctx) := do
  -- normalize: just normalize
  match c with
    | .atom a =>
      pure (c, ctx) -- return itself
    | .var n =>
      let c: Term β ← Context.get? ctx n
      c.normalize ctx
    | _ => sorry



end EL2
