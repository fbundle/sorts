import EL2.Term
import EL2.Util

namespace EL2

def reduceParams (params: List (Ann α)) (f: Ctx → α → Option (Ctx × β)) (ctx: Ctx): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMapAll params ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β))) ctx

partial def matchParamsArgs (params: List (Ann α)) (argsType: List α) (le: α → α → Option Unit): Option Unit := do
  if params.length = 0 ∧ argsType.length = 0 then
    ()
  else
    let headParam ← params.head?
    let headArgsType ← argsType.head?
    let _ ← le headArgsType headParam.type

    let tailParams := params.extract 1
    let tailArgsType := argsType.extract 1
    matchParamsArgs tailParams tailArgsType le

partial def equal [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (x: Term β) (y: Term β): Option Unit := do
  sorry

partial def inferType [Irreducible β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- (ctx: Ctx) - map name -> type
  match term with
    | .atom a =>
      pure (ctx, atom Irreducible.inferType a)

    | .t t => match t with
      | .var n =>
        let type: Term β ← Context.get? ctx n
        pure (ctx, type)

      | .lst {init, last} =>
        let (ctx, _) ← Util.optionCtxMapAll init inferType ctx
        inferType ctx last

      | .bind_val {name, value} =>
        let (ctx, valueType) ← inferType ctx value
        let ctx := Context.set ctx name valueType
        pure (ctx, valueType)

      | .bind_typ {name, params, parent} =>
        let (ctx, _) ← reduceParams params inferType ctx
        let (ctx, _) ← inferType ctx parent
        pure (Context.set ctx name term, parent)

      | .bind_mk {name, params, type} =>
        let (ctx, _) ← reduceParams params inferType ctx

        let (typeName, typeArgs) := (type.cmd, type.args)
        let (ctx, typeArgsType) ← Util.optionCtxMapAll typeArgs inferType ctx

        match Context.get? ctx typeName with
          | some (Term.t (T.bind_typ {name, params, parent})) =>
            let _ ← matchParamsArgs params typeArgsType (equal ctx)
            let ctx := Context.set ctx name (sorry: Term β) -- type constructor

            pure (ctx, sorry)
          | _ => none

      | .lam {params, body} =>
        sorry

      | .app {cmd, args} => sorry
      | .mat {cond, cases} => sorry



end EL2
