import EL2.Term
import EL2.Util

namespace EL2

partial def infer [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- infer: infer type
  match term with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (ctx, .atom p)
    | .var n =>
      let term : Term β ← Context.get? ctx n
      infer ctx term
    | _ => sorry

partial def normalize [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- normalize
  match term with
    | .atom a =>
      pure (ctx, term)

    | .var n =>
      let term: Term β ← Context.get? ctx n
      normalize ctx term

    | .list init tail =>
      let (ctx, _) ← Util.optionCtxMapAll init normalize ctx
      normalize ctx tail

    | .bind_val {name, value} =>
      let (ctx, value) ← normalize ctx value
      let ctx := Context.set ctx name value
      pure (ctx, value)

    | .bind_typ {name, params, parent} =>
      let (ctx, params) ← Util.optionCtxMapAll params ((λ ctx {name, type} => do
        let (ctx, type) ← normalize ctx type
        pure (ctx, {name := name, type := type : Ann (Term β)})
      ): Ctx → Ann (Term β) → Option (Ctx × (Ann (Term β)))) ctx

      let (ctx, parent) ← normalize ctx parent
      let value: Term β := .bind_typ {
        name := name,
        params := params,
        parent := parent,
      }
      let ctx := Context.set ctx name value
      pure (ctx, value)

    | .bind_mk {name, params, type} =>
      let (ctx, params) ← Util.optionCtxMapAll params ((λ ctx {name, type} => do
        let (ctx, type) ← normalize ctx type
        pure (ctx, {name := name, type := type : Ann (Term β)})
      ): Ctx → Ann (Term β) → Option (Ctx × (Ann (Term β)))) ctx

      let {cmd, args} := type
      let (ctx, args) ← Util.optionCtxMapAll args normalize ctx

      let value: Term β := .bind_mk {
        name := name,
        params := params,
        type := {
          cmd := cmd,
          args := args,
        },
      }
      let ctx := Context.set ctx name value
      pure (ctx, value)

    | _ => sorry



end EL2
