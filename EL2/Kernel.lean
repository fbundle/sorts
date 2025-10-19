import EL2.Term
import EL2.Util

namespace EL2

partial def Term.infer [Irreducible β] [Context Ctx (Term β)] (c: Term β) (ctx: Ctx) : Option (Ctx × Term β) := do
  -- infer: infer type
  match c with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (ctx, .atom p)
    | .var n =>
      let c : Term β ← Context.get? ctx n
      c.infer ctx
    | _ => sorry

partial def Term.normalize [Irreducible β] [Context Ctx (Term β)] (c: Term β) (ctx: Ctx): Option (Ctx × Term β) := do
  -- normalize
  match c with
    | .atom a =>
      pure (ctx, c)

    | .var n =>
      let c: Term β ← Context.get? ctx n
      c.normalize ctx

    | .list init tail =>
      let (ctx, _) ← Util.optionCtxMapAll init ((λ ctx term =>
        term.normalize ctx
      ): Ctx → Term β → Option (Ctx × Term β)) ctx
      tail.normalize ctx

    | .bind_val {name, value} =>
      let (ctx, value) ← value.normalize ctx
      let ctx := Context.set ctx name value
      pure (ctx, value)

    | .bind_typ {name, params, parent} =>
      let (ctx, params) ← Util.optionCtxMapAll params ((λ ctx ann => do
        let (ctx, type) ← ann.type.normalize ctx
        pure (ctx, {name := ann.name, type := type : Ann (Term β)})
      ): Ctx → Ann (Term β) → Option (Ctx × Ann (Term β))) ctx
      let (ctx, parent) ← parent.normalize ctx
      let value: Term β := .bind_typ {
        name := name,
        params := params,
        parent := parent,
      }
      let ctx := Context.set ctx name value
      pure (ctx, value)

    | _ => sorry



end EL2
