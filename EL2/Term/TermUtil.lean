import EL2.Term.Term
import EL2.Term.Util
import Std


namespace EL2.Term

def T.mapM [Monad m] (t: T α β) (fa: α → m γ) (fb: β → m δ) : m (T γ δ) := do
  match t with
    | T.inh x =>
      let type ← fa x.type
      let args ← x.args.mapM fa
      pure (T.inh {
        type := type,
        cons := x.cons,
        args := args
      })
    | T.bnd x =>
      let init ← x.init.mapM (λ bind => do
        let value ← fa bind.value
        pure {
          name := bind.name,
          value := value,
          : Bind γ
        }
      )
      let last ← fa x.last
      pure (T.bnd {
        init := init,
        last := last,
      })
    | T.lam x =>
      let params ← x.params.mapM (λ param => do
        let type ← fa param.type
        pure {
          name := param.name,
          type := type,
          : Param γ
        }
      )
      let body ← fb x.body
      pure (T.lam {
        params := params,
        body := body,
      })
    | T.app x =>
      let cmd ← fa x.cmd
      let args ← x.args.mapM fa
      pure (T.app {
        cmd := cmd,
        args := args,
      })
    | T.mat x =>
      let cond ← fa x.cond
      let cases ← x.cases.mapM (λ case => do
        let value ← fa case.value
        pure {
          patCmd := case.patCmd,
          patArgs := case.patArgs,
          value := value
          : Case γ
        }
      )
      pure (T.mat {
        cond := cond,
        cases := cases,
      })

def T.map (t: T α β) (fa: α → γ) (fb: β → δ): T γ δ :=
  Id.run (t.mapM (λ x => pure (fa x)) (λ y => pure (fb y)))

def Term.mapM [Monad m] (term: Term) (f: Term → m Term): m Term := do
  match term with
    | .univ _ => pure term
    | .var _ => pure term
    | .t x =>
      let y ← x.mapM f f
      pure (Term.t y)


def Term.map (term: Term) (f: Term → Term): Term :=
  Id.run (term.mapM (λ x => pure (f x)))

partial def ReducedTerm.toTerm (reducedTerm: ReducedTerm): Term :=
  match reducedTerm with
    | .univ level => Term.univ level
    | .pi x =>
      let y := (T.lam x).map ReducedTerm.toTerm ReducedTerm.toTerm
      Term.t y
    | .t x =>
      let y := x.map ReducedTerm.toTerm id
      Term.t y

instance: Coe ReducedTerm Term where
  coe rt := rt.toTerm

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

-- util
def isTermLam? (term: Term): Option (Lam Term Term) :=
  match term with
    | .t (.lam x) => some x
    | _ => none

def isReducedTermLam? (term: ReducedTerm): Option (Lam ReducedTerm Term) :=
  match term with
    | .t (.lam x) => some x
    | _ => none


namespace EL2.Term
