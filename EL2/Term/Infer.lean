import EL2.Term.Term
import Std
import EL2.Term.Print

namespace EL2.Term

class Context M α where
  size: M → Nat
  set: M → String → α → M
  get?: M → String → Option α

instance : Context (Std.HashMap String α) α where
  size := Std.HashMap.size
  set := Std.HashMap.insert
  get? := Std.HashMap.get?


def emptyNameMap: Std.HashMap String String := Std.HashMap.emptyWithCapacity

partial def renameTerm [Repr M] [Context M String] (nameMap: M) (term: Term): Term :=
  -- rename all parameters into _<count> where count = nameNameMap.size save into nameMap
  -- rename all variables according to nameMap
  match term with
    | var name =>
      match Context.get? nameMap name with
        | none => term
        | some newName => var newName
    | lam x =>
      let (newNameMap, newParams) := Util.statefulMap x.params nameMap (λ nameMap param =>
        let count := Context.size (α := String) nameMap
        let newName := s!"_{count}"
        let newType := renameTerm nameMap param.type
        let newNameMap := Context.set nameMap param.name newName
        (newNameMap, {name := newName, type := newType : Ann Term})
      )
      let newBody := renameTerm newNameMap x.body
      lam {
        params := newParams,
        body := newBody,
      }
    | mat x =>
      let newCond := renameTerm nameMap x.cond
      let newCases := x.cases.map (λ case =>
        let (newNameMap, newPatArgs) := Util.statefulMap case.patArgs nameMap (λ nameMap patArg =>
          let count := Context.size (α := String) nameMap
          let newName := s!"_{count}"
          let newNameMap := Context.set nameMap patArg newName
          (newNameMap, newName)
        )
        let newValue := renameTerm newNameMap case.value
        {
          patCmd := case.patCmd,
          patArgs := newPatArgs,
          value := newValue,
          : Case Term
        }
      )
      mat {
        cond := newCond,
        cases := newCases,
      }
    | _ => term.map (renameTerm nameMap)

def renameCaseFromParams (params: List (Ann Term)) (case: Case Term): Term :=
  let (newNameMap, _) := Util.statefulMap (List.zip case.patArgs params) emptyNameMap (λ nameMap (patArg, param) =>
    let newNameMap := Context.set nameMap patArg param.name
    (newNameMap, ())
  )
  let newValue := renameTerm newNameMap case.value
  newValue

partial def isSubTypeMany? (type1List: List Term) (type2List: List Term): Option Unit := do
  if type1List.length = 0 ∧ type2List.length = 0 then
    pure ()
  else
    let type1 ← type1List.head?
    let type2 ← type2List.head?
    if type1 != type2 then
      dbg_trace s!"[DBG_TRACE] different type"
      dbg_trace s!"type1:\t{type1}"
      dbg_trace s!"type2:\t{type2}"
      none
    else
      isSubTypeMany? (type1List.extract 1) (type2List.extract 1)

structure InferedType where
  -- for type safety, we may want to change term: Term into Option ReducedTerm where ReducedTerm doesn't contain any variable
  -- currently, term was not reduced fully
  term : Term -- currently we are using Typ as a hole, next we can replace by Option Term
  type : Term
  level : Int
  -- ctx [Context Ctx InferedType] : Ctx -- Step 2 - TODO add context when term is not reduced



instance: ToString InferedType where
  toString (iterm: InferedType) :=
    s!"term: {iterm.term} type: {iterm.type} level: {iterm.level}"

instance: Repr InferedType where
  reprPrec (iterm: InferedType) (prec: Nat): Std.Format := toString iterm


mutual
partial def inferType? [Repr Ctx] [Context Ctx InferedType] (ctx: Ctx) (term: Term): Option InferedType := do
  -- recursively do WHNF and type infer
  match term with
    | univ level =>
      pure {
        term := term,
        type := univ (level + 1),
        level := level + 1, -- U_1 is at level 2
      }

    | var name =>
      let iX: InferedType ← Context.get? ctx name
      let iTerm := match iX.term with
          | typ _ => term -- if param then return itself -- TODO change to None
          | _ => iX.term
      pure {
        term := iTerm,
        type := iX.type,
        level := iX.level,
      }

    | inh x =>
      let iX ← inferType? ctx x.type
      pure {
        term := term,
        type := x.type,
        level := iX.level - 1,
      }

    | typ x => none -- typ is deprecated

    | bnd x =>
      let (subCtx, _) ← Util.statefulMap? x.init ctx (λ subCtx bind => do
        let iValue ← inferType? subCtx bind.value
        let subCtx := Context.set subCtx bind.name iValue
        pure (subCtx, ())
      )

      inferType? subCtx x.last

    | lam x =>
      let (subCtx, iNamedParams) ← Util.statefulMap? x.params ctx (λ subCtx param => do
        let iType ← inferType? subCtx param.type
        dbg_trace s!"[DBG_TRACE] iType {param.type} -> {iType.term}"
        let iParamValue := {
          term := typ {value := iType.term},
          type := iType.term,
          level := iType.level - 1,
          : InferedType
        }
        let subCtx := Context.set subCtx param.name iParamValue
        pure (subCtx, (param.name, iParamValue))
      )

      let newParams := iNamedParams.map (λ (name, iParamValue) => {
        name := name,
        type := iParamValue.type,
        : Ann Term
      })

      let iBody ← inferType? subCtx x.body
      let lamLevel := (iNamedParams.map (λ (name, iParamValue) => iParamValue.level)).foldl max iBody.level

      pure {
        term := lam {
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
      let iCmd ← inferType? ctx x.cmd
      let iLamCmd ← isLam? iCmd.term
      let iArgs ← Util.optionMap? x.args (inferType? ctx)

      let iArgsType ← Util.optionMap? iArgs (λ iArg => inferType? ctx iArg.type)
      let iParamsType ← Util.optionMap? iLamCmd.params (λ param => inferType? ctx param.type)

      let _ ← isSubTypeMany? (iArgsType.map (λ iArgType => iArgType.term)) (iParamsType.map (λ iParamType => iParamType.term))

      -- WHNF
      let (subCtx, _) ← Util.statefulMap? (List.zip iLamCmd.params iArgs) ctx (λ subCtx (param, arg) => do
        let subCtx := Context.set subCtx param.name arg
        pure (subCtx, ())
      )
      inferType? subCtx iLamCmd.body

    | mat x =>
      let casesTypeLevel ← Util.optionMap? x.cases (λ case => do
        let iCmd: InferedType ← Context.get? ctx case.patCmd
        match isLam? iCmd.term with
          | none =>
            let iValue ← inferType? ctx case.value
            let iType ← inferType? ctx iValue.type
            pure ({
              patCmd := case.patCmd,
              patArgs := case.patArgs,
              value := iType.term,
              : Case Term
            }, iValue.level)

          | some iLamCmd =>
            -- assume term is already renamed
            let value := renameCaseFromParams iLamCmd.params case
            let caseLam := lam {
              params := iLamCmd.params,
              body := value,
            }
            let iCaseLam ← inferType? ctx caseLam
            let iCaseLamType ← isLam? iCaseLam.type
            let iType ← inferType? ctx iCaseLamType.body
            pure ({
              patCmd := case.patCmd,
              patArgs := case.patArgs,
              value := iType.term,
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
