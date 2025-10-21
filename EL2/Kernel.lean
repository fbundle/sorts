import EL2.Term
import EL2.Util

namespace EL2

namespace Util

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

end Util

structure InferedTerm where
  term: Term
  type: Term
  deriving Repr, BEq

class Context Ctx where
  set: Ctx → String → InferedTerm → Ctx
  get?: Ctx → String → Option InferedTerm

partial def infer? [Repr Ctx] [Context Ctx] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × InferedTerm) := do
  let inferMany? (reduce: Bool) (ctx: Ctx) (lst: List Term): Option (Ctx × List InferedTerm) :=
    Util.optionCtxMap? lst ctx (infer? reduce)

  let isBindTyp? (term1: Term): Option (BindTyp Term) :=
    match term with
      | .bind_typ bind_typ => some bind_typ
      | _ => none


  let rec bindParams? (ctx: Ctx) (params: List (Ann Term)) (args: List InferedTerm): Option Ctx := do
    if params.length = 0 ∧ args.length = 0 then
      pure ctx
    else
      let headParam ← params.head?
      let headArg ← args.head?
      if headParam.type != headArg.type then
        none
      else

      let newCtx := Context.set ctx headParam.name headArg

      let tailParams := params.extract 1
      let tailArgs := args.extract 1
      bindParams? newCtx tailParams tailArgs

  -- (ctx: Ctx) - map name -> type
  match term with
    | .inh type =>
      pure (ctx, {term := term, type := type})

    | .univ level =>
      pure (ctx, {term := term, type := .univ (level+1)})

    | .var name =>
      let term1: InferedTerm ← Context.get? ctx name
      pure (ctx, term1)

    | .lst {init, last} =>
      let (ctx, _) ← inferMany? reduce ctx init
      infer? reduce ctx last

    | .bind_val {name, value} =>
      let (_, value1) ← infer? reduce ctx value
      pure (Context.set ctx name value1, value1)

    | .bind_typ {name, params, level} =>
      let (_, _) ← inferMany? reduce ctx (params.map (λ ann => ann.type))
      -- type of type definition is Pi
      let term1: InferedTerm := {
        term := term,
        type := .lam {
          params := params,
          body := .univ level,
        },
      }

      pure (Context.set ctx name term1, term1)

    | .bind_mk {name, params, type := {cmd := typeName, args := typeArgs}} =>
      -- resolve params
      let (paramCtx, paramsType1) ← inferMany? true ctx (params.map (λ {name, type} => type)) -- always reduce at type level
      -- bind dummy params
      let dummyArgs := paramsType1.map (λ {term, type} => {
        term := .inh term,
        type := term,
        : InferedTerm
      })
      let dummyCtx ← bindParams? ctx params dummyArgs -- always ok
      -- type check type application
      let type1 ← Context.get? dummyCtx typeName
      let {name := typeName, params := typeParams, level := typeLevel} ← isBindTyp? type1.term -- type of bind_typ is Pi/lam
      let (_, typeArgs1) ← inferMany? true dummyCtx typeArgs -- always reduce at type level

      let argsCtx ← bindParams? ctx params typeArgs1 -- type check for args

      -- make output
      let term1: InferedTerm := {
        term := term,
        type := .lam {
          params := params,
          body := .bind_typ {
            name := typeName,
            params := typeParams,
            level := typeLevel,
          },
        }
      }

      pure (Context.set ctx name term1, term1)

    | .typ {value} =>
      let (ctx, value1) ← infer? reduce ctx value
      infer? reduce ctx value1.type

    | .lam {params, body} =>
      -- type of parent is Pi
      let term1: InferedTerm := {
        term := term,
        type := .lam {
          params := params,
          body := .typ {value := body},
          -- we use typ to create a future type infer object
          -- normalizing typ will invoke inferType?
        },
      }

      pure (ctx, term1)

    | .app {cmd, args} =>
      pure (ctx, {
        term := .univ 1,
        type := .univ 2,
      })
      -- TODO change it
    | .mat {cond, cases} =>
      pure (ctx, {
        term := .univ 1,
        type := .univ 2,
      })
      -- TODO change it



end EL2
