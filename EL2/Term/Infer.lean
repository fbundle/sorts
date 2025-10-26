import EL2.Term.Term
import EL2.Term.TermUtil
import EL2.Term.Util
import Std
import EL2.Term.Print

namespace EL2.Term

class Map M α where
  size: M → Nat
  set: M → String → α → M
  get?: M → String → Option α

structure InferedType where
  type: Term -- type of term
  level: Int -- level of term
  deriving Repr

-- TODO change Option to Except String
partial def inferType? [Repr Ctx] [Map Ctx InferedType] (ctx: Ctx) (term: Term) : Option InferedType := do
  dbg_trace s!"[DBG_TRACE] infering at {term}"
  -- recursively do WHNF and type infer
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
        let (subCtx, iValues) ← Util.statefulMapM x.params ctx (λ subCtx param => do
          -- dummy value
          let value := inh {
            type := param.type,
            cons := "",
            args := []
          }
          let iValue ← inferType? subCtx value
          let subCtx := Map.set subCtx param.name iValue

          pure (subCtx, iValue)
        )

        let iBody ← inferType? subCtx x.body
        let lamLevel := (iValues.map (λ iValue => iValue.level)).foldl max (iBody.level)

        let newParams := (List.zip x.params iValues).map (λ (param, iValue) =>
          {
            name := param.name,
            type := iValue.type,
            : Param Term
          }
        )

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
        let iLam ← isLam? iCmd.type
        let iArgs ← x.args.mapM (inferType? ctx)
        -- type check
        let paramsType := iLam.params.map (λ param => param.type)
        let argsType := iArgs.map (λ iArg => iArg.type)

        let _ ← (List.zip argsType paramsType).mapM (λ (argType, paramType) => do
          let _ ← isSubType? argType paramType
          pure ()
        )

        inferType? ctx (inh {
          type := iLam.body,
          cons := "",
          args := []
        })

      | mat x =>
        let iCases: List (Case InferedType) ← x.cases.mapM (λ case => do
          let iCmd: InferedType ← Map.get? ctx case.patCmd
          match isLam? iCmd.type with
            | none => -- case is not a lambda - resolve directly
              let iValue ← inferType? ctx case.value
              pure {
                patCmd := case.patCmd,
                patArgs := case.patArgs,
                value := iValue,
              }

            | some iLam => -- case is lambda
              -- rename case to match iLam
              let newCase ← renameCase? iLam case
              -- convert case to lambda to reuse inferType?
              let iValueLam1 := lam {
                params := iLam.params,
                body := newCase.value,
              }
              let iValueLam2 ← inferType? ctx iValueLam1
              -- iValueLam2 is a Pi type (iLam.params) -> typeof case.value
              let iValueLam3 ← isLam? iValueLam2.type
              let iValue := iValueLam3.body

              pure {
                patCmd := newCase.patCmd,
                patArgs := newCase.patArgs,
                value := {
                  type := iValue,
                  level := iValueLam2.level,
                }
              }
        )
        let iCond ← inferType? ctx x.cond
        let level := (iCases.map (λ iCase => iCase.value.level)).foldl max iCond.level

        pure {
          type := mat {
            cond := x.cond,
            cases := iCases.map (λ iCase =>
              {
                patCmd := iCase.patCmd,
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
    | some v => pure v



end EL2.Term
