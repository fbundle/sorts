import EL2.Term.Term
import EL2.Term.Print
import EL2.Term.Util
import Std


namespace EL2.Term

def T.mapM [Monad m] (t: T α) (f: α → m β) : m (T β) := do
  match t with
    | T.inh x =>
      let type ← f x.type
      let args ← x.args.mapM f
      pure (T.inh {
        type := type,
        cons := x.cons,
        args := args
      })
    | T.bnd x =>
      let init ← x.init.mapM (λ bind => do
        let value ← f bind.value
        pure {
          name := bind.name,
          value := value,
          : Bind β
        }
      )
      let last ← f x.last
      pure (T.bnd {
        init := init,
        last := last,
      })
    | T.lam x =>
      let params ← x.params.mapM (λ param => do
        let type ← f param.type
        pure {
          name := param.name,
          type := type,
          : Param β
        }
      )
      let body ← f x.body
      pure (T.lam {
        params := params,
        body := body,
      })
    | T.app x =>
      let cmd ← f x.cmd
      let args ← x.args.mapM f
      pure (T.app {
        cmd := cmd,
        args := args,
      })
    | T.mat x =>
      let cond ← f x.cond
      let cases ← x.cases.mapM (λ case => do
        let value ← f case.value
        pure {
          patCmd := case.patCmd,
          patArgs := case.patArgs,
          value := value
          : Case β
        }
      )
      pure (T.mat {
        cond := cond,
        cases := cases,
      })

def T.map (t: T α) (f: α → β) : T β :=
  Id.run (t.mapM (λ x => pure (f x)))

def Term.mapM [Monad m] (term: Term) (f: Term → m Term): m Term := do
  match term with
    | .univ _ => pure term
    | .var _ => pure term
    | .t x =>
      let y ← x.mapM f
      pure (Term.t y)


def Term.map (term: Term) (f: Term → Term): Term :=
  Id.run (term.mapM (λ x => pure (f x)))


def emptyNameMap: Std.HashMap String String := Std.HashMap.emptyWithCapacity

def dummyName (nameMap: Std.HashMap String String): String :=
  s!"_{nameMap.size}"

partial def renameTerm (nameMap: Std.HashMap String String) (term: Term): Term :=
  -- nameMap holds a mapping oldName -> newName
  -- rename all parameters into _<count> where count = nameNameMap.size save into nameMap
  -- rename all match parameters into _<count>
  -- rename all variables according to nameMap
  match term with
    | var oldName =>
      match nameMap.get? oldName with
        | some newName => var newName
        | none => var oldName

    | lam x =>
      let (newNameMap, newParams) := Util.statefulMap x.params nameMap (λ oldNameMap param =>
        let newType := renameTerm oldNameMap param.type

        let newName := dummyName oldNameMap
        let newNameMap := oldNameMap.insert param.name newName

        (newNameMap, {name := newName, type := newType : Param Term})
      )
      let newBody := renameTerm newNameMap x.body
      lam {
        params := newParams,
        body := newBody,
      }
    | mat x =>
      let newCond := renameTerm nameMap x.cond
      let newCases := x.cases.map (λ case =>
        let (newNameMap, newPatArgs) := Util.statefulMap case.patArgs nameMap (λ oldNameMap patArg =>
          let newName := dummyName oldNameMap
          let newNameMap := oldNameMap.insert patArg newName

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

def renameParamsWithCase (params: List (Param Term)) (patArgs: List String): List (Param Term) :=
  -- given patArgs and constructor
  -- rename the param and type of constructor to match patArgs
  let (newNameMap, newParams) := Util.statefulMap (List.zip patArgs params) emptyNameMap (λ oldNameMap (patArg, param) =>
    let newType := renameTerm oldNameMap param.type
    let newName := patArg
    let newNameMap := oldNameMap.insert param.name newName
    (newNameMap, {
      name := newName,
      type := newType,
      : Param Term
    })
  )
  newParams

partial def isSubType? (type1: Term) (type2: Term): Option Unit := do
  if type1 != type2 then
    dbg_trace s!"[DBG_TRACE] different type"
    dbg_trace s!"type1:\t{type1}"
    dbg_trace s!"type2:\t{type2}"
    none
  else
    pure ()


-- util
def isLam? (term: Term): Option (Lam Term) :=
  match term with
    | lam x => some x
    | _ => none

namespace EL2.Term
