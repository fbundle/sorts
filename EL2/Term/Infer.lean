import EL2.Term.TermUtil

namespace EL2.Term.Infer

def inhEmpty (type: Term) : Term :=
  inh {
    type := type,
    cons := "",
    args := []
  }

structure InferedType where
  type: Term -- type of term
  level: Int -- level of term
  deriving Repr

partial def isSubType (iType1: InferedType) (iType2: InferedType): Bool :=
  let type1 := iType1.type.normalizeName
  let type2 := iType2.type.normalizeName
  match type2 with
    | univ _ =>
      if iType1.level ≤ iType2.level then
        true
      else
        false
    | _ =>
      if type1 == type2 then
        True
      else
        dbg_trace s!"[DBG_TRACE] different type"
        dbg_trace s!"type1:\t{type1}"
        dbg_trace s!"type2:\t{type2}"
        false

mutual

partial def inferTypeParams? [Repr Ctx] [Map Ctx InferedType] (ctx: Ctx) (params: List (Ann Term)): Option (Ctx × List (Bind InferedType)) :=
  Util.statefulMapM params ctx (λ subCtx param => do
    let paramValue := inhEmpty param.type -- dummy value
    let iParamValue ← inferType? subCtx paramValue
    let subCtx := Map.set subCtx param.name iParamValue
    pure (subCtx, {
      name := param.name,
      value := iParamValue,
    })
  )
partial def inferTypeCase? [Repr Ctx] [Map Ctx InferedType] (ctx: Ctx) (case: Case Term): Option (Case InferedType) := do
  let iCons: InferedType ← Map.get? ctx case.patCons
  match isLam? iCons.type with
    | none => -- case is not a lambda - resolve directly
      let iValue ← inferType? ctx case.value
      pure {
        patCons := case.patCons,
        patArgs := case.patArgs,
        value := iValue,
      }

    | some iCons => -- case is lambda
      -- make new set of params according to patArgs
      let newParams := renameParamsWithCase iCons.params case.patArgs
      -- convert case to lambda to reuse inferType?
      let matLam := lam {
        params := newParams,
        body := case.value,
      }
      let iMatLam ← inferType? ctx matLam
      -- iMatLam is a Pi type (newParams) -> typeof case.value
      let valueType := (← isLam? iMatLam.type).body

      pure {
        patCons := case.patCons,
        patArgs := case.patArgs,
        value := {
          type := valueType,
          level := iMatLam.level,
        },
      }

partial def inferType? [Repr Ctx] [Map Ctx InferedType] (ctx: Ctx) (term: Term) : Option InferedType := do
  -- recursively type infer (probably will do WHNF)
  let o: Option InferedType := do
    match term with
      | univ level =>
        pure {
          type := univ level + 1,
          level := level + 1,-- U_1 is at level 2
        }

      | var name =>
        Map.get? ctx name

      | inh x =>
        let iType ← inferType? ctx x.type
        pure {
          type := x.type,
          level := iType.level - 1,
        }

      | bnd x =>
        let (subCtx, _) ← Util.statefulMapM x.init ctx (λ subCtx bind => do
          let iValue ← inferType? subCtx bind.value
          let subCtx := Map.set subCtx bind.name iValue
          pure (subCtx, ())
        )

        inferType? subCtx x.last

      | lam x =>
        let (subCtx, iParams) ← inferTypeParams? ctx x.params
        let iBody ← inferType? subCtx x.body
        let lamLevel := (iParams.map (λ iParam => iParam.value.level)).foldl max (iBody.level)

        let newParams: List (Ann Term) := (List.zip x.params iParams).map (λ (param, iParam) => {
          name := param.name,
          type := iParam.value.type,
        })

        pure {
          type := lam {
            params := newParams,
            body := iBody.type,
          },
          level := lamLevel,
        }

      | app x =>
        -- infer
        let iCmd ← inferType? ctx x.cmd
        let iCmd ← isLam? iCmd.type
        let iArgs ← x.args.mapM (inferType? ctx)
        -- type check
        let (_, iParams) ← inferTypeParams? ctx iCmd.params
        let _ ← (List.zip iArgs iParams).mapM (λ (iArg, iParam) => do
          if isSubType iArg iParam.value then pure () else
            none
        )

        inferType? ctx (inhEmpty iCmd.body)

      | mat x =>
        let iCases: List (Case InferedType) ← x.cases.mapM (inferTypeCase? ctx)
        let iCond ← inferType? ctx x.cond
        let level := (iCases.map (λ iCase => iCase.value.level)).foldl max iCond.level

        pure {
          type := mat {
            cond := x.cond,
            cases := iCases.map (λ iCase =>
              {
                patCons := iCase.patCons,
                patArgs := iCase.patArgs,
                value := iCase.value.type,
              }
            ),
          }
          level := level,
        }

  match o with
    | none =>
      dbg_trace s!"[DBG_TRACE] failed at {term} with ctx {repr ctx}"
      none
    | some v =>
      pure v
end

-- TODO think of some way to reduce type and reduce in general
-- because we currently can only compare two types if they're reduced


end EL2.Term.Infer
