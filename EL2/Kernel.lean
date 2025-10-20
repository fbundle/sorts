import EL2.Term
import EL2.Util

namespace EL2

def reduceParamsWithName? (params: List (Ann α))(ctx: Ctx) (f: Ctx → String → α → Option (Ctx × β)): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMap? params ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx name type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β))) ctx

def reduceParams? (params: List (Ann α)) (ctx: Ctx) (f: Ctx → α → Option (Ctx × β)): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMap? params ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β))) ctx

partial def matchParamsArgs? [BEq α] (params: List (Ann α)) (argsType: List α): Option Unit := do
  if params.length = 0 ∧ argsType.length = 0 then
    ()
  else
    let headParam ← params.head?
    let headArgsType ← argsType.head?
    if headParam.type != headArgsType then
      none
    else

    let tailParams := params.extract 1
    let tailArgsType := argsType.extract 1
    matchParamsArgs? tailParams tailArgsType

partial def inferType? [Irreducible β] [BEq β] [Context Ctx (Term β)] (ctx: Ctx) (term: Term β): Option (Ctx × Term β) := do
  -- (ctx: Ctx) - map name -> type
  match term with
    | atom a =>
      pure (ctx, atom Irreducible.inferType a)

    | var n =>
      let parent: Term β ← Context.get? ctx n
      pure (ctx, parent)

    | lst {init, last} =>
      let (ctx, _) ← Util.optionCtxMap? init inferType? ctx
      inferType? ctx last

    | bind_val {name, value} =>
      let (ctx, parent) ← inferType? ctx value
      let ctx := Context.set ctx name parent
      pure (ctx, parent)

    | bind_typ {name, params, parent} =>
      let (ctx, _) ← reduceParams? params ctx inferType?
      let (ctx, _) ← inferType? ctx parent
      pure (Context.set ctx name term, parent)

    | bind_mk {name, params, type} =>
      let (ctx, _) ← reduceParams? params ctx inferType?

      let (typeName, typeArgs) := (type.cmd, type.args)
      let (ctx, typeArgsType) ← Util.optionCtxMap? typeArgs inferType? ctx

      match Context.get? ctx typeName with
        | some (bind_typ type) =>
          let {name := typeName, params := typeParams, parent := typeParent} := type
          let _ ← matchParamsArgs? typeParams typeArgsType
          -- type of a type constructor is Pi
          let parent := lam {
            params := params,
            body := bind_typ type,
          }
          let ctx := Context.set ctx name parent
          pure (ctx, parent)
        | _ => none

    | typ {value} =>
      let (ctx, type) ← inferType? ctx value
      let (ctx, parent) ← inferType? ctx type
      pure (ctx, parent)

    | lam {params, body} =>
      -- type of parent is Pi
      let parent := lam {
        params := params,
        body := typ {value := body},
        -- we use typ to create a future type infer object
        -- normalizing typ will invoke inferType
      }
      pure (ctx, parent)

    | app {cmd, args} => sorry
    | mat {cond, cases} => sorry



end EL2
