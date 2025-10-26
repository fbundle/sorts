import EL2.Term.Term
import EL2.Term.TermUtil
import EL2.Term.Util
import Std
import EL2.Term.Print

namespace EL2.Term

-- Ctx -- TODO change this to HashMap
notation "Ctx" α => (List (String × α))
def ctxEmpty {α}: Ctx α := []
def ctxGet? (ctx: Ctx α) (name: String): Option α :=
  match ctx with
    | [] => none
    | (key, value) :: tail =>
      if name = key then
        some value
      else
        ctxGet? tail name

def ctxKeys (ctx: Ctx α): List String := ctx.map (λ (key, value) => key)
def ctxSet (ctx: Ctx α) (name: String) (value: α): Ctx α :=
  (name, value) :: ctx

def ctxSize (ctx: Ctx α): Int := ctx.length

structure InferedType where
  type: Term -- type of term
  level: Nat -- level of term

def typLevel (level: Nat): Nat := level + 1
def inhLevel (level: Nat): Option Nat :=
  match level with
    | 0 => none
    | _ => some (level - 1)


-- TODO change Option to Except String
partial def infer? (ctx: Ctx InferedType) (term: Term) : Option InferedType := do
  -- recursively do WHNF and type infer
  let o: Option InferedType := do
    match term with
      | univ level =>
        pure {
          type := univ typLevel level,
          level := typLevel level,-- U_1 is at level 2
        }

      | var name =>
        ctxGet? ctx name

      | inh x =>
        let iType ← infer? ctx x.type
        pure {
          type := x.type,
          level := ← inhLevel iType.level
        }

      | bnd x =>
        let (subCtx, _) ← Util.statefulMapM x.init ctx (λ subCtx bind => do
          let iValue ← infer? subCtx bind.value
          let subCtx := ctxSet subCtx bind.name iValue
          pure (subCtx, ())
        )

        infer? subCtx x.last

      | lam x =>
        let (subCtx, iValues) ← Util.statefulMapM x.params ctx (λ subCtx param => do
          -- dummy value
          let value := inh {
            type := param.type,
            cons := "",
            args := []
          }
          let iValue ← infer? subCtx value
          let subCtx := ctxSet subCtx param.name iValue

          pure (subCtx, iValue)
        )

        let iBody ← infer? subCtx x.body
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
        let iCmd ← infer? ctx x.cmd
        let iLam ← isLam? iCmd.type
        let iArgs ← x.args.mapM (infer? ctx)
        -- type check
        let paramsType := iLam.params.map (λ param => param.type)
        let argsType := iArgs.map (λ iArg => iArg.type)

        let _ ← (List.zip argsType paramsType).mapM (λ (argType, paramType) => do
          let _ ← isSubType? argType paramType
          pure ()
        )
        -- WHNF
        let (subCtx, _) ← Util.statefulMapM (List.zip iLam.params iArgs) ctx (λ subCtx (param, arg) => do
          let subCtx := ctxSet subCtx param.name arg
          pure (subCtx, ())
        )
        infer? subCtx iLam.body

      | mat x =>
        let iCases: List (Case InferedType) ← x.cases.mapM (λ case => do
          let iCmd ← ctxGet? ctx case.patCmd
          match isLam? iCmd.type with
            | none => -- case is not a lambda - resolve directly
              let iValue ← infer? ctx case.value
              pure {
                patCmd := case.patCmd,
                patArgs := case.patArgs,
                value := iValue,
              }

            | some iLam => -- case is lambda
              -- rename case to match iLam
              let newCase ← renameCase? iLam case
              -- convert case to lambda to reuse infer?
              let iValueLam1 := lam {
                params := iLam.params,
                body := newCase.value,
              }
              let iValueLam2 ← infer? ctx iValueLam1
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
        let iCond ← infer? ctx x.cond
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
      dbg_trace s!"[DBG_TRACE] failed at {term} with ctx {ctxKeys ctx}"
      none
    | some v => pure v



end EL2.Term
