import EL2.Term
import EL2.Util

namespace EL2

class Context Ctx α where
  insert: Ctx → String → α → Ctx
  get?: Ctx → String → Option α

def reduceParamsWithName? (params: List (Ann α)) (ctx: Ctx) (f: Ctx → String → α → Option (Ctx × β)): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMap? params ctx ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx name type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β)))

def reduceParams? (params: List (Ann α)) (ctx: Ctx) (f: Ctx → α → Option (Ctx × β)): Option (Ctx × List (Ann β)) :=
  Util.optionCtxMap? params ctx ((λ ctx {name, type} => do
    let (ctx, type) ← f ctx type
    pure (ctx, {name := name, type := type})
  ): Ctx → Ann α → Option (Ctx × (Ann β)))

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

mutual

partial def inferType? [Repr Ctx] [Context Ctx Term] (ctx: Ctx) (term: Term): Option (Ctx × Term) := do
  -- (ctx: Ctx) - map name -> type
  match term with
    | univ level =>
      pure (ctx, univ (level+1))

    | var n =>
      let parent: Term ← Context.get? ctx n
      pure (ctx, parent)

    | lst {init, last} =>
      let (ctx, _) ← Util.optionCtxMap? init ctx inferType?
      inferType? ctx last

    | bind_val {name, value} =>
      let (ctx, parent) ← inferType? ctx value
      let ctx := Context.insert ctx name parent
      pure (ctx, parent)

    | bind_typ {name, params, level} =>
      let (ctx, _) ← reduceParams? params ctx inferType?
      -- type of type definition is Pi
      let parent := lam {
        params := params,
        body := univ level,
      }
      let ctx := Context.insert ctx name parent
      pure (ctx, parent)

    | bind_mk {name, params, type} =>

      dbg_trace s!"1 checking bind_mk {name} {repr params}"

      let (ctx, _) ← reduceParams? params ctx inferType?

      let {cmd := typeName, args := typeArgs} := type
      dbg_trace s!"2 checking bind_mk {name} {repr params}"

      match Context.get? ctx typeName with
        | some (lam typeType) =>
          -- type of type is Pi/lam

          let {params := typeParams, body := _} := typeType
          dbg_trace s!"4 checking bind_mk {name} {repr typeParams}"

          -- set dummy args
          let (ctx, _) ← reduceParamsWithName? typeParams ctx ((λ ctx name value =>
            let ctx := Context.insert ctx name value
            -- TODO possibly need to normalize/reduce this
            -- reduce = typeInfer? + substitution
            some (ctx, value)
          ))
          dbg_trace s!"5 checking bind_mk {name} {repr ctx}"

          -- resolve return type
          let (ctx, typeArgsType) ← Util.optionCtxMap? typeArgs ctx inferType?



          let _ ← matchParamsArgs? typeParams typeArgsType


          -- type of a type constructor is Pi
          let parent := lam {
            params := params,
            body := lam typeType,
          }
          let ctx := Context.insert ctx name parent
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
        -- normalizing typ will invoke inferType?
      }
      pure (ctx, parent)

    | app {cmd, args} =>
      dbg_trace s!"1 checking app {repr cmd} {repr args}"

      let (ctx, cmdType) ← inferType? ctx cmd
      match cmdType with
        | lam {params := cmdTypeParams, body := cmdTypeBody} =>
          -- type of bind_typ, bind_mk, lam is lam/Pi
          let (ctx, argsType) ← Util.optionCtxMap? args ctx inferType?
          let _ ← matchParamsArgs? cmdTypeParams argsType
          -- set dummy args
          let (ctx, _) ← reduceParamsWithName? cmdTypeParams ctx ((λ ctx name value =>
            let ctx := Context.insert ctx name value
            some (ctx, value)
          ))
          -- return the type of body given the context
          pure (ctx, cmdTypeBody)
        | _ => none
    | mat {cond, cases} =>
      pure (ctx, univ 1)
      -- TODO change it

partial def reduceTerm? [Repr Ctx] [Context Ctx Term] (ctx: Ctx) (term: Term): Option (Ctx × Term) :=
  sorry

end


end EL2
