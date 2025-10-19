import EL2.Term
import EL2.Util

namespace EL2

partial def infer [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- infer: infer type
  match term with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (ctx, .atom p)
    | .t t => match t with
      | .var n =>
        let term : Term β ← Context.get? ctx n
        infer ctx term
      | _ => sorry

def reduceParams (params: List (Ann α)) (f: Ctx → α → Option (Ctx × β)) (ctx: Ctx): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMapAll params ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β))) ctx

partial def normalize [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- normalize
  match term with
    | .atom a =>
      pure (ctx, term)

    | .t t => match t with
      | .var n =>
        let term: Term β ← Context.get? ctx n
        pure (ctx, term)

      | .lst {init, tail} =>
        let (ctx, _) ← Util.optionCtxMapAll init normalize ctx
        normalize ctx tail

      | .bind_val {name, value} =>
        let (ctx, value) ← normalize ctx value
        pure (Context.set ctx name value)

      | .bind_typ {name, params, parent} =>
        let (ctx, params) ← reduceParams params normalize ctx
        let (ctx, parent) ← normalize ctx parent
        let value: Term β := bind_typ {
          name := name,
          params := params,
          parent := parent,
        }

        pure (Context.set ctx name value)

      | .bind_mk {name, params, type} =>
        let (ctx, params) ← reduceParams params normalize ctx

        let {cmd, args} := type
        let (ctx, args) ← Util.optionCtxMapAll args normalize ctx

        let value := bind_mk {
          name := name,
          params := params,
          type := {
            cmd := cmd,
            args := args,
          },
        }
        pure (Context.set ctx name value)

      | .lam {params, body} =>
        let (ctx, params) ← reduceParams params normalize ctx
        let value := lam {
          params := params,
          body := body,
        }
        pure (ctx, value)

      | .app {cmd, args} => sorry
      | .mat {cond, cases} => sorry



end EL2
