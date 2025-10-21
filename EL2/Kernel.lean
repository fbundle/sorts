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

structure InferedTerm where
  term: Term
  type: Term
  deriving Repr, BEq

partial def infer? [Repr Ctx] [Context Ctx InferedTerm] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × InferedTerm) := do
  -- (ctx: Ctx) - map name -> type
  match term with
    | .inh type =>
      pure (ctx, {term := term, type := type})

    | .univ level =>
      pure (ctx, {term := term, type := .univ (level+1)})

    | .var name =>
      let inferedTerm: InferedTerm ← Context.get? ctx name
      pure (ctx, inferedTerm)

    | .lst {init, last} =>
      let (ctx, _) ← Util.optionCtxMap? init ctx (infer? reduce)
      infer? reduce ctx last

    | .bind_val {name, value} =>
      let (ctx, inferedValue) ← infer? reduce ctx value
      let ctx := Context.insert ctx name inferedValue
      pure (ctx, inferedValue)

    | .bind_typ {name, params, level} =>
      let (ctx, _) ← reduceParams? params ctx (infer? reduce)
      -- type of type definition is Pi
      let inferedTerm: InferedTerm := {
        term := term,
        type := .lam {
          params := params,
          body := .univ level,
        },
      }

      let ctx := Context.insert ctx name inferedTerm
      pure (ctx, inferedTerm)

    | .bind_mk {name, params, type} =>

      dbg_trace s!"1 checking bind_mk {name} {repr params}"

      let (ctx, _) ← reduceParams? params ctx (infer? reduce)

      let {cmd := typeName, args := typeArgs} := type
      dbg_trace s!"2 checking bind_mk {name} {repr params}"

      match Context.get? ctx typeName with
        | some {term := typeTerm, type := .lam typeType: InferedTerm} =>
          -- type of type is Pi/lam

          let {params := typeParams, body := _} := typeType
          dbg_trace s!"4 checking bind_mk {name} {repr typeParams}"

          -- set dummy args
          let (ctx, _) ← reduceParamsWithName? typeParams ctx ((λ ctx paramName paramType => do
            let (ctx, inferedParamType) ← infer? true ctx paramType -- for type level, do reduce
            let ctx := Context.insert ctx paramName {
              term := .inh inferedParamType.term,
              type := inferedParamType.term,
              : InferedTerm
            }
            some (ctx, inferedParamType)
          ))
          dbg_trace s!"5 checking bind_mk {name} {repr ctx}"

          -- resolve return type
          let (ctx, typeArgsInferedTerm) ← Util.optionCtxMap? typeArgs ctx (infer? reduce)

          let typeArgsType := typeArgsInferedTerm.map (λ inferedTerm => inferedTerm.type)
          dbg_trace s!"6 checking bind_mk {name} {repr typeArgsType}"

          let _ ← matchParamsArgs? typeParams typeArgsType


          -- type of a type constructor is Pi
          let inferedTerm: InferedTerm := {
            term := term,
            type := .lam {
              params := params,
              body := .lam typeType,
            },
          }
          let ctx := Context.insert ctx name inferedTerm
          pure (ctx, inferedTerm)
        | _ => none

    | .typ {value} =>
      let (ctx, inferedValue) ← infer? reduce ctx value
      infer? reduce ctx inferedValue.type

    | .lam {params, body} =>
      -- type of parent is Pi
      let inferedTerm: InferedTerm := {
        term := term,
        type := .lam {
          params := params,
          body := .typ {value := body},
          -- we use typ to create a future type infer object
          -- normalizing typ will invoke inferType?
        },
      }

      pure (ctx, inferedTerm)

    | .app {cmd, args} =>
      dbg_trace s!"1 checking app {repr cmd} {repr args}"

      let (ctx, cmdType) ← infer? reduce ctx cmd
      match cmdType with
        | {term := cmdTerm, type := .lam {params := cmdTypeParams, body := cmdTypeBody}} =>
          -- type of bind_typ, bind_mk, lam is lam/Pi
          let (ctx, argsInferedTerm) ← Util.optionCtxMap? args ctx (infer? reduce)
          let argsType := argsInferedTerm.map (λ inferedTerm => inferedTerm.type)
          let _ ← matchParamsArgs? cmdTypeParams argsType

          if reduce then
            sorry
          else
            -- set dummy args
            let (ctx, _) ← reduceParamsWithName? cmdTypeParams ctx ((λ ctx paramName paramType => do
              let (ctx, inferedParamType) ← infer? true ctx paramType -- for type level, do reduce
              let ctx := Context.insert ctx paramName {
                term := .inh inferedParamType.term,
                type := inferedParamType.term,
                : InferedTerm
              }
              some (ctx, inferedParamType)
            ))
            -- return the type of body given the context
            pure (ctx, {
              term := .inh cmdTypeBody,
              type := cmdTypeBody,
            })

        | _ => none
    | .mat {cond, cases} =>
      pure (ctx, {
        term := .univ 1,
        type := .univ 2,
      })
      -- TODO change it



end EL2
