import EL2.Term.Term
import EL2.Term.Print
import EL2.Term.Util
import EL2.Term.NameMap
import Std


namespace EL2.Term


-- TODO - decompose into
-- 1. parameter renaming
-- 2. typ unwrapping
-- 3. reduce


def emptyNameMap: Std.HashMap String String := Std.HashMap.emptyWithCapacity
def emptyFrame: Std.HashMap String InferedTerm := Std.HashMap.emptyWithCapacity

partial def renameTerm [Repr M] [NameMap M String] (nameMap: M) (term: Term): Term :=
  -- rename all parameters into _<count> where count = nameNameMap.size save into nameMap
  -- rename all variables according to nameMap
  match term with
    | var name =>
      match NameMap.get? nameMap name with
        | none => term
        | some newName => var newName
    | lam x =>
      let (newNameMap, newParams) := Util.statefulMap x.params nameMap (λ nameMap param =>
        let count := NameMap.size (α := String) nameMap
        let newType := renameTerm nameMap param.type
        let newName := s!"_{count}"
        let newNameMap := NameMap.set nameMap param.name newName
        (newNameMap, {name := newName, type := newType : Ann Term})
      )
      let newBody := renameTerm newNameMap x.body
      lam {
        params := newParams,
        body := newBody,
      }
    | _ => term.map (renameTerm nameMap)

partial def unwrapTyp (term: Term): Term :=
  match term with
    | typ x =>
      match x.value with
      | inh y =>
        let output := y.type
        --dbg_trace s!"[DBG_TRACE] unwrapping {typ x} -> {output}"
        output
      | mat y =>
        let output := mat {
          y with
          cases := y.cases.map (λ case => {
            case with value := typ {value := case.value}
          })
        }
        --dbg_trace s!"[DBG_TRACE] unwrapping {typ x} -> {output}"
        output
      | lam y =>
        let output := lam {
          y with
          body := typ {value := y.body}
        }
        --dbg_trace s!"[DBG_TRACE] unwrapping {typ x} -> {output}"
        output
      | _ =>
        term
    | _ =>
      term.map unwrapTyp

def isSubType? (argType: Term) (paramType: Term): Option Unit := do
  -- TODO - to remove these two unwrapTyp
  let argType := unwrapTyp argType
  let paramType := unwrapTyp paramType
  if argType == paramType then
    pure ()
  else
    dbg_trace s!"[DBG_TRACE] different type"
    dbg_trace s!"argType:\t{argType}"
    dbg_trace s!"paramType:\t{paramType}"
    none

mutual
partial def reduceMany? [Repr F] [NameMap F InferedTerm] (frame: F) (terms: List Term): Option (List (InferedTerm)) :=
  Util.optionMap? terms (reduceTerm? frame)

partial def reduceUniv? [Repr F] [NameMap F InferedTerm] (frame: F) (level: Int): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_univ {univ level}"
  let output := {
    term := univ level,
    type := univ level+1,
  }
  --dbg_trace s!"[DBG_TRACE] reduce_univ_ok {univ level} → {output}"
  pure output

partial def reduceVar? [Repr F] [NameMap F InferedTerm] (frame: F) (name: String): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_var {var name}"
  let output ← NameMap.get? frame name
  --dbg_trace s!"[DBG_TRACE] reduce_var_ok {var name} → {output}"
  pure output

partial def reduceInh? [Repr F] [NameMap F InferedTerm] (frame: F) (x: Inh Term): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_inh {inh x}"
  let iType ← reduceTerm? frame x.type
  let iArgs ← reduceMany? frame x.args
  let output := {
    term := inh {
      type := iType.term,
      cons := x.cons,
      args := iArgs.map (λ iterm => iterm.term),
    },
    type := iType.term,
  }
  --dbg_trace s!"[DBG_TRACE] reduce_inh_ok {inh x} → {output}"
  pure output

partial def reduceTyp? [Repr F] [NameMap F InferedTerm] (frame: F) (x: Typ Term): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_typ {typ x}"
  -- unwrap usual types - put typ as close to leaf as possible - not complete
  let iValue ← reduceTerm? frame x.value
  let iType ← reduceTerm? frame iValue.type
  let unwrappedType := unwrapTyp iType.term
  match unwrappedType with
    | typ _ =>
      -- cannot unwrap further
      --dbg_trace s!"[DBG_TRACE] reduce_typ_ok {typ x} → {iType}"
      pure iType
    | _ =>
      -- unwrapped
      reduceTerm? frame unwrappedType

partial def reduceBnd? [Repr F] [NameMap F InferedTerm] (frame: F) (x: Bnd Term): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_bnd {bnd x}"
  let (frame, _) ← Util.statefulMap? x.init frame (λ frame {name, value} => do
    let iValue ← reduceTerm? frame value
    let frame := NameMap.set frame name iValue
    some (frame, iValue)
  )
  let iLast ← reduceTerm? frame x.last
  --dbg_trace s!"[DBG_TRACE] reduce_bnd_ok {bnd x} → {iLast}"
  pure iLast

partial def reduceLam? [Repr F] [NameMap F InferedTerm] (frame: F) (x: Lam Term): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_lam {lam x}"
  let output := {
    term := lam x,
    type := lam {
      params := x.params,
      body := typ {
        value := x.body,
      }
    },
  }
  --dbg_trace s!"[DBG_TRACE] reduce_lam_ok {lam x} → {output}"
  pure output

partial def bindParamsWithArgs? [Repr F] [NameMap F InferedTerm] (frame: F) (params: List (Ann Term)) (args: List Term): Option F := do
  if params.length = 0 ∧ args.length = 0 then
    pure frame
  else
    let param ← params.head?
    let arg ← args.head?

    let iParamType ← reduceTerm? frame param.type
    let iArg ← reduceTerm? frame arg
    let iArgType ← reduceTerm? frame iArg.type

    match isSubType? iArgType.term iParamType.term with
      | none =>
        --dbg_trace s!"[DBG_TRACE] bind_params_with_args type check failed {iArgType} -> {iParamType}"
        none
      | some _ =>
        let frame := NameMap.set frame param.name iArg
        let params := params.extract 1
        let args := args.extract 1
        bindParamsWithArgs? frame params args

partial def reduceApp? [Repr F] [NameMap F InferedTerm] (frame: F) (x: App Term): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_app {app x}"
  let iCmd ← reduceTerm? frame x.cmd
  let lamCmd ← isLam? iCmd.term
  let argFrame ← bindParamsWithArgs? frame lamCmd.params x.args
  let iBody ← reduceTerm? argFrame lamCmd.body
  --dbg_trace s!"[DBG_TRACE] reduce_app_ok {app x} → {iBody}"
  pure iBody

partial def matchCases? (inhCond: Inh Term) (cases: List (Case Term)): Option (Bnd Term) := do
  let headCase ← cases.head?
  if inhCond.cons = headCase.patCmd ∧ inhCond.args.length = headCase.patArgs.length then
    let init: List (Bind Term) := (List.zip headCase.patArgs inhCond.args).map (λ (name, value) => {
      name := name,
      value := value,
    })
    pure {
      init := init,
      last := headCase.value,
    }
  else
    let cases := cases.extract 1
    matchCases? inhCond cases

partial def reduceMat? [Repr F] [NameMap F InferedTerm] (frame: F) (x: Mat Term): Option InferedTerm := do
  --dbg_trace s!"[DBG_TRACE] reduce_mat {mat x}"
  let iCond ← reduceTerm? frame x.cond
  let inhCond ← isInh? iCond.term
  let terms ← matchCases? inhCond x.cases
  let output ← reduceTerm? frame (bnd terms)
  --dbg_trace s!"[DBG_TRACE] reduce_mat_ok {mat x} → {output}"
  pure output

partial def reduceTerm? [Repr F] [NameMap F InferedTerm] (frame: F) (term: Term): Option InferedTerm := do
    -- TODO - currently normalize term doesn't have access to frame
  -- so it cannot normalize further things like Nat into (inh U_1 Nat)
  match term with
    | univ level => reduceUniv? frame level
    | var name => reduceVar? frame name
    | inh x => reduceInh? frame x
    | typ x => reduceTyp? frame x
    | bnd x => reduceBnd? frame x
    | lam x => reduceLam? frame x
    | app x => reduceApp? frame x
    | mat x => reduceMat? frame x





end

end EL2.Term
