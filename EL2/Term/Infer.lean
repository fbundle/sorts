import EL2.Term.Term
import Std


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



structure InferedType where
  term : Term
  type : Term
  level : Int

mutual
partial def inferType? [Repr Ctx] [Context Ctx InferedType] (ctx: Ctx) (term: Term): Option InferedType := do
  -- rename first -- params, patArgs will be always _0 _1 ...
  let term := renameTerm emptyNameMap term
  -- recursively do WHNF and type infer
  match term with
    | univ level =>
      pure {
        term := term,
        type := univ (level + 1),
        level := level + 1, -- U_1 is at level 2
      }

    | var name =>
      Context.get? ctx name

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
        let iParamValue := {
          term := inh {type := param.type, cons := param.name, args := []},
          type := param.type,
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

      if iArgs.map (λ arg => arg.type) != iLamCmd.params.map (λ param => param.type) then
        none
      else
        -- WHNF
        let (subCtx, _) ← Util.statefulMap? (List.zip iLamCmd.params iArgs) ctx (λ subCtx (param, arg) => do
          let subCtx := Context.set subCtx param.name arg
          pure (subCtx, ())
        )
        inferType? subCtx iLamCmd.body

    | mat x =>
      let iCond ← inferType? ctx x.cond
      let casesTypeLevel ← Util.optionMap? x.cases (λ case => do
        let iCmd: InferedType ← Context.get? ctx case.patCmd
        let iLamCmd ← isLam? iCmd.term
        -- this works since we already renamed everything
        let caseLam := lam {
          params := iLamCmd.params,
          body := case.value,
        }
        let iCaseLam ← inferType? ctx caseLam
        let iCaseLamType ← isLam? iCaseLam.type
        let type := iCaseLamType.body
        pure ({
          patCmd := case.patCmd,
          patArgs := case.patArgs,
          value := type,
          : Case Term
        }, iCaseLam.level)
      )

      let casesType := casesTypeLevel.map (λ (type, level) => type)
      let casesLevel := casesTypeLevel.map (λ (type, level) => level)

      let level ← casesLevel.max?
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
