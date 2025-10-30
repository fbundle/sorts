import EL2.Term.Term
import EL2.Term.Print

import Std

namespace EL2.Util

def statefulMapM [Monad m] (xs: List α) (state: State) (f: State → α → m (State × β)) : m (State × List β) :=
  let rec loop (ys: Array β) (xs: List α) (state: State): m (State × List β) := do
    match xs with
      | [] => pure (state, ys.toList)
      | x :: xs =>
        let (state, y) ← f state x
        loop (ys.push y) xs state

  loop #[] xs state

def statefulMap (xs: List α) (state: State) (f: State → α → State × β): State × List β :=
  Id.run (statefulMapM xs state (λ s x => pure (f s x)))

structure Ctx α where
  list: List (String × α)
  deriving Repr

partial def Ctx.get? (ctx: Ctx α) (name: String): Option α :=
  match ctx.list with
    | [] => none
    | (key, val) :: list =>
      if name = key then
        some val
      else
        {list := list: Ctx α}.get? name

partial def Ctx.set (ctx: Ctx α) (name: String) (val: α): Ctx α :=
  {list := (name, val) :: ctx.list}

def emptyCtx: Ctx α := {list := []}

end EL2.Util

namespace EL2.Term

class Map M α where
  size: M → Nat
  set: M → String → α → M
  get?: M → String → Option α

-- TODO - change Monad to Applicative
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
          : Ann β
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
          patCons := case.patCons,
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

partial def Term.normalizeName (term: Term) (nameMap: Std.HashMap String String := emptyNameMap) : Term :=
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
        let newType := param.type.normalizeName oldNameMap

        let newName := dummyName oldNameMap
        let newNameMap := oldNameMap.insert param.name newName

        (newNameMap, {name := newName, type := newType : Ann Term})
      )
      let newBody := x.body.normalizeName newNameMap
      lam {
        params := newParams,
        body := newBody,
      }

    | mat x =>
      let newCond := x.cond.normalizeName nameMap
      let newCases := x.cases.map (λ case =>
        let (newNameMap, newPatArgs) := Util.statefulMap case.patArgs nameMap (λ oldNameMap patArg =>
          let newName := dummyName oldNameMap
          let newNameMap := oldNameMap.insert patArg newName

          (newNameMap, newName)
        )
        let newValue := case.value.normalizeName newNameMap
        {
          patCons := case.patCons,
          patArgs := newPatArgs,
          value := newValue,
          : Case Term
        }
      )
      mat {
        cond := newCond,
        cases := newCases,
      }

    | _ => term.map (·.normalizeName nameMap)

def renameParamsWithCase (params: List (Ann Term)) (patArgs: List String): List (Ann Term) :=
  -- given patArgs and constructor
  -- rename the param and type of constructor to match patArgs
  let (_, newParams) := Util.statefulMap (List.zip patArgs params) emptyNameMap (λ oldNameMap (patArg, param) =>
    let newType := param.type.normalizeName oldNameMap
    let newName := patArg
    let newNameMap := oldNameMap.insert param.name newName
    (newNameMap, {
      name := newName,
      type := newType,
      : Ann Term
    })
  )
  newParams


-- util
def isLam? (term: Term): Option (Lam Term) :=
  match term with
    | lam x => some x
    | _ => none

namespace EL2.Term
