import EL2.Term.Term
import EL2.Term.TermUtil
import EL2.Term.Util
import Std
import EL2.Term.Print

namespace EL2.Term

def renameCase? (cons: Lam Term) (case: Case Term): Option (Case Term) := do
  -- rename case patArgs according to constructor
  -- return renamed value
  let (newNameMap, newPatArgs) := Util.statefulMap (List.zip case.patArgs cons.params) emptyNameMap (λ nameMap (patArg, param) =>
    let newNameMap := Context.set nameMap patArg param.name -- rename patArg to paramNam
    (newNameMap, param.name)
  )
  let newValue ← renameTerm? newNameMap case.value
  pure {
    patCmd := case.patCmd,
    patArgs := newPatArgs,
    value := newValue,
  }

partial def isSubType? (type1: Term) (type2: Term): Option Unit := do
  if type1 != type2 then
    dbg_trace s!"[DBG_TRACE] different type"
    dbg_trace s!"type1:\t{type1}"
    dbg_trace s!"type2:\t{type2}"
    none
  else
    pure ()

structure FutureTerm where
  ctx [Context Ctx Term]: Ctx
  term: Term

structure InferedTerm where
  term?: Option Term
  type: FutureTerm
  level: Int

instance: ToString InferedTerm where
  toString (iterm: InferedTerm) :=
    match iterm.term? with
      | none =>
        s!"term: none type: {iterm.type.term} level: {iterm.level}"
      | some term =>
        s!"term: {term} type: {iterm.type.term} level: {iterm.level}"

instance: Repr InferedTerm where
  reprPrec (iterm: InferedTerm) (_: Nat): Std.Format := toString iterm


mutual
partial def inferType? [Repr Ctx] [Context Ctx InferedTerm] (ctx: Ctx) (term: Term): Option InferedTerm := do
  -- recursively do WHNF and type infer
  match term with
    | univ level =>
      pure {
        term? := term,
        type := univ (level + 1),
        level := level + 1, -- U_1 is at level 2
      }

    | var name =>
      Context.get? ctx name

    | inh x =>
      let iType ← inferType? ctx x.type
      let iArgs ← x.args.mapM (inferType? ctx)

      pure {
        term? := inh {
          type := ← iType.term?,
          cons := x.cons,
          args := ← iArgs.mapM (λ iArg => iArg.term?),
        },
        type := ← iType.term?,
        level := iType.level - 1,
      }

    | bnd x =>
      let (subCtx, _) ← Util.statefulMapM x.init ctx (λ subCtx bind => do
        let iValue ← inferType? subCtx bind.value
        let subCtx := Context.set subCtx bind.name iValue
        pure (subCtx, ())
      )

      inferType? subCtx x.last

    | lam x =>
      let (subCtx, iNamedTypes) ← Util.statefulMapM x.params ctx (λ subCtx param => do
        let iType ← inferType? subCtx param.type
        let iParamValue := {
          term? := none, -- dummy
          type := ← iType.term?,
          level := iType.level - 1,
          : InferedTerm
        }
        let subCtx := Context.set subCtx param.name iParamValue

        pure (subCtx, (param.name, iType))
      )

      let iBody ← inferType? subCtx x.body
      let lamLevel := (iNamedTypes.map (λ (name, iType) => iType.level)).foldl max (iBody.level + 1)

      let newParams ← iNamedTypes.mapM (λ (name, iType) => do
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
        type := lam {
          params := newParams,
          body := iBody.type,
        },
        level := lamLevel,
      }

    | app x =>
      -- infer
      let iCmd ← inferType? ctx x.cmd
      let iLam ← isLam? (← iCmd.term?)
      let iArgs ← x.args.mapM (inferType? ctx)
      -- type check
      let iArgsType ← iArgs.mapM (λ iArg => inferType? ctx iArg.type)
      let iParamsType ← iLam.params.mapM (λ param => inferType? ctx param.type)
      let _ ← (List.zip iArgsType iParamsType).mapM (λ (iArgType, iParamType) => do
        let _ ← isSubType? (← iArgType.term?) (← iParamType.term?)
        pure ()
      )
      -- WHNF
      let (subCtx, _) ← Util.statefulMapM (List.zip iLam.params iArgs) ctx (λ subCtx (param, arg) => do
        let subCtx := Context.set subCtx param.name arg
        pure (subCtx, ())
      )
      inferType? subCtx iLam.body

    | mat x =>
      let casesTypeLevel ← x.cases.mapM (λ case => do
        let iCmd: InferedTerm ← Context.get? ctx case.patCmd
        let iCmdTerm ← iCmd.term?
        match isLam? iCmdTerm with
          | none => -- case is not a lambda
            let iValue ← inferType? ctx case.value
            let iType ← inferType? ctx iValue.type
            pure ({
              patCmd := case.patCmd,
              patArgs := case.patArgs,
              value := ← iType.term?,
              : Case Term
            }, iValue.level)

          | some iLam => -- case is lambda
            -- rename case
            let case ← renameCase? iLam case
            -- convert case to lambda to reuse inferType?
            let iCase ← inferType? ctx (lam {
              params := iLam.params,
              body := case.value,
            })
            let iCaseLamType ← isLam? iCase.type
            -- let iType ← inferType? ctx iCaseLamType.body
            let iTypeTerm := iCaseLamType.body
            pure ({
              patCmd := case.patCmd,
              patArgs := case.patArgs,
              value := iTypeTerm,
              : Case Term
            }, iCaseLam.level)
      )

      let casesType := casesTypeLevel.map (λ (type, level) => type)
      let casesLevel := casesTypeLevel.map (λ (type, level) => level)

      let level ← casesLevel.max?

      let iCond ← inferType? ctx x.cond
      dbg_trace s!"[DBG_TRACE] cond {x.cond} → {iCond}"
      pure {
        term := term,
        type := mat {
          cond := iCond.term,
          cases := casesType,
        },
        level := level,
      }

end


end EL2.Term
