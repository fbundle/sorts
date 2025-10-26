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

structure InferedTerm where
  -- term? : resolved value of term
  term?: Option Term
  -- typeCtx and typeTerm : recipe to resolve typeTerm
  typeCtx: Ctx InferedTerm
  typeTerm: Term
  -- level : level of term
  level: Int -- change to Nat

-- TODO change Option to Except String
partial def infer? (ctx: Ctx InferedTerm) (term: Term) : Option InferedTerm := do
  -- recursively do WHNF and type infer
  let o: Option InferedTerm := do
    match term with
      | univ level =>
        pure {
          term? := term,
          typeCtx := ctx,
          typeTerm := univ (level + 1),
          level := level + 1, -- U_1 is at level 2
        }

      | var name =>
        ctxGet? ctx name

      | inh x =>
        let iType ← infer? ctx x.type
        let iArgs ← x.args.mapM (infer? ctx)
        let term ← iType.term?
        let args ← iArgs.mapM (λ iArg => iArg.term?)
        pure {
          term? := inh {
            type := term,
            cons := x.cons,
            args := args,
          },
          typeCtx := ctx,
          typeTerm := x.type,
          level := iType.level - 1,
        }

      | bnd x =>
        let (subCtx, _) ← Util.statefulMapM x.init ctx (λ subCtx bind => do
          let iValue ← infer? subCtx bind.value
          let subCtx := ctxSet subCtx bind.name iValue
          pure (subCtx, ())
        )

        infer? subCtx x.last

      | lam x =>
        let (subCtx, iNamedParamsType) ← Util.statefulMapM x.params ctx (λ subCtx param => do
          let iParamType ← infer? subCtx param.type
          let iParamValue := {
            term? := none, -- dummy param
            -- TODO - for other functions
            -- if resolved value have none term?
            -- return none term? as well
            typeCtx := subCtx,
            typeTerm := param.type,
            level := iParamType.level - 1,
            : InferedTerm
          }
          let subCtx := ctxSet subCtx param.name iParamValue

          pure (subCtx, (param.name, iParamType))
        )

        let iBody ← infer? subCtx x.body
        let lamLevel := (iNamedParamsType.map (λ (name, iParamType) => iParamType.level)).foldl max (iBody.level + 1)

        let newParams ← iNamedParamsType.mapM (λ (name, iType) => do
          pure {
            name := name,
            type := ← iType.term?,
            : Param Term
          }
        )
        pure {
          term? := lam {
            params := newParams,
            body := x.body,
          },
          -- TODO - think
          -- with subCtx we can resolve typeTerm.body immediately
          -- TODO change to ctx, iBody.typeTerm
          typeCtx := subCtx,
          typeTerm := lam {
            params := newParams,
            body := iBody.typeTerm,
          },
          level := lamLevel,
        }

      | app x =>
        -- infer
        let iCmd ← infer? ctx x.cmd
        let iLam ← isLam? (← iCmd.term?)
        let iArgs ← x.args.mapM (infer? ctx)
        -- type check
        let iParamsType ← iLam.params.mapM (λ param => infer? ctx param.type)
        let paramsType ← iParamsType.mapM (λ iParam => iParam.term?)
        let paramsType := iLam.params.map (λ param => param.type) -- TODO remove if the code doesn't work

        let iArgsType ← iArgs.mapM (λ iArg => infer? iArg.typeCtx iArg.typeTerm)
        let argsType ← iArgsType.mapM (λ iArg => iArg.term?)

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
        let caseTypeLevels: List (Case Term × Case Term × Int) ← x.cases.mapM (λ case => do
          let iCmd: InferedTerm ← ctxGet? ctx case.patCmd
          let iCmdTerm ← iCmd.term?
          match isLam? iCmdTerm with
            | none => -- case is not a lambda - resolve directly
              let iValue ← infer? ctx case.value
              let iType ← infer? iValue.typeCtx iValue.typeTerm
              pure (
                {
                  patCmd := case.patCmd,
                  patArgs := case.patArgs,
                  value := ← iValue.term?,
                },
                {
                  patCmd := case.patCmd,
                  patArgs := case.patArgs,
                  value := ← iType.term?,
                },
                iValue.level,
              )

            | some iLam => -- case is lambda
              -- rename case to match iLam
              let case ← renameCase? iLam case
              -- convert case to lambda to reuse infer?
              let iCaseLam ← infer? ctx (lam {
                params := iLam.params,
                body := case.value,
              })
              -- iCaseLamType is a Pi type (iLam.params) -> typeof case.value
              let iCaseLamPi ← infer? iCaseLam.typeCtx iCaseLam.typeTerm
              let valueType := (← isLam? (← iCaseLamPi.term?)).body

              pure (
                case,
                {
                  patCmd := case.patCmd,
                  patArgs := case.patArgs,
                  value := valueType,
                },
                iCaseLam.level,
                )
        )

        -- list of pairs to pair of lists
        let (cases, types, levels) := caseTypeLevels.foldr (λ (x, y, z) (xs, ys, zs)  => (x :: xs, y :: ys, z ::  zs)) ([], [], [])
        let level ← levels.max? -- fail if no case
        let iCond ← infer? ctx x.cond

        pure {
          term? := mat {
            cond := ← iCond.term?,
            cases := cases,
          },
          typeCtx := ctx, -- types is resolved, this is not important
          typeTerm := mat {
            cond := ← iCond.term?,
            cases := types,
          },
          level := level,
        }

  match o with
    | none =>
      dbg_trace s!"[DBG_TRACE] failed at {term} with ctx {ctxKeys ctx}"
      none
    | some v => pure v



end EL2.Term
