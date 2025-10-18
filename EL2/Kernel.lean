import EL2.Code

namespace EL2


partial def Code.inferCode [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx) : Option (Code β × Ctx) := do
  -- infer: turn everything to type then normalize
  match c with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (.atom p, ctx)
    | .var n =>
      let c : Code β ← Context.get? ctx n
      c.inferCode ctx
    | _ => sorry

partial def Code.normalizeCode [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx): Option (Code β × Ctx) := do
  -- normalize: just normalize
  match c with
    | .atom a =>
      pure (c, ctx) -- return itself
    | .var n =>
      let c: Code β ← Context.get? ctx n
      c.normalizeCode ctx
    | _ => sorry



end EL2
