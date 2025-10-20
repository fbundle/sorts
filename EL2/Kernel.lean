import EL2.Term
import EL2.Util

namespace EL2



def reduceParams (params: List (Ann α)) (f: Ctx → α → Option (Ctx × β)) (ctx: Ctx): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMapAll params ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β))) ctx

partial def inferType [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- (ctx: Ctx) - map name -> type
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
        let ctx := Context.set ctx name value
        pure (ctx, value)

      | .bind_typ {name, params, parent} =>
        let (ctx, params) ← reduceParams params normalize ctx
        let (ctx, parent) ← normalize ctx parent
        let value: Term β := bind_typ {
          name := name,
          params := params,
          parent := parent,
        }

        pure (Context.set ctx name value, value)

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
        let ctx := Context.set ctx name value
        pure (ctx, value)

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
