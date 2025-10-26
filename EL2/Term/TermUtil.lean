import EL2.Term.Term
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

def T.map (t: T α) (f: α → β): T β :=
  Id.run (t.mapM (λ x => Id.run (f x)))

def Term.mapM [Monad m] (term: Term) (f: Term → m Term): m Term := do
  match term with
    | .univ _ => pure term
    | .var _ => pure term
    | .t x =>
      let y ← x.mapM f
      pure (Term.t y)


def Term.map (term: Term) (f: Term → Term): Term :=
  Id.run (term.mapM (λ x => pure (f x)))

partial def Term.toReducedTerm? (term: Term): Option ReducedTerm := do
  match term with
    | .univ level => pure (ReducedTerm.univ level)
    | .var _ => none
    | .t x =>
      let y ← x.mapM Term.toReducedTerm?
      pure (ReducedTerm.t y)

partial def ReducedTerm.toTerm (reducedTerm: ReducedTerm): Term :=
  match reducedTerm with
    | .univ level => Term.univ level
    | .t x =>
      let y := x.map ReducedTerm.toTerm
      Term.t y

class Context M α where
  size: M → Nat
  set: M → String → α → M
  get?: M → String → Option α

instance : Context (Std.HashMap String α) α where
  size := Std.HashMap.size
  set := Std.HashMap.insert
  get? := Std.HashMap.get?

def emptyNameMap: Std.HashMap String String := Std.HashMap.emptyWithCapacity

def dummyName [Context N String] (nameMap: N): String :=
  let count := Context.size (α := String) nameMap
  s!"_{count}"

partial def renameTerm? [Repr N] [Context N String] (nameMap: N) (term: Term): Option Term := do
  -- nameMap holds a mapping oldName -> newName
  -- rename all parameters into _<count> where count = nameNameMap.size save into nameMap
  -- rename all match parameters into _<count>
  -- rename all variables according to nameMap
  match term with
    | var oldName =>
      let newName ← Context.get? nameMap oldName
      pure (var newName)

    | lam x =>
      let (newNameMap, newParams) ← Util.statefulMapM x.params nameMap (λ oldNameMap param => do
        let newType ← renameTerm? oldNameMap param.type

        let newName := dummyName oldNameMap
        let newNameMap := Context.set oldNameMap param.name newName

        (newNameMap, {name := newName, type := newType : Param Term})
      )
      let newBody ← renameTerm? newNameMap x.body
      pure (lam {
        params := newParams,
        body := newBody,
      })
    | mat x =>
      let newCond ← renameTerm? nameMap x.cond
      let newCases ← x.cases.mapM (λ case => do
        let (newNameMap, newPatArgs) := Util.statefulMap case.patArgs nameMap (λ oldNameMap patArg =>
          let newName := dummyName oldNameMap
          let newNameMap := Context.set oldNameMap patArg newName

          (newNameMap, newName)
        )
        let newValue ← renameTerm? newNameMap case.value
        pure {
          patCmd := case.patCmd,
          patArgs := newPatArgs,
          value := newValue,
          : Case Term
        }
      )
      pure (mat {
        cond := newCond,
        cases := newCases,
      })

    | _ => term.mapM (renameTerm? nameMap)



namespace EL2.Term
