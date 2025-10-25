import EL2.Term.Util

namespace EL2.Term

structure Inh (α: Type) where
  type: α
  cons: String -- cons and args make sure object is constructed uniquely
  args: List α -- i.e. (inh Nat succ zero) = (inh Nat succ zero) by definition
  deriving Repr, BEq

structure Typ (α: Type) where
  value: α
  deriving Repr, BEq

structure Bind (α: Type) where
  name: String
  value: α
  deriving Repr, BEq

structure Bnd (α: Type) where
  init: List (Bind α)
  last: α
  deriving Repr, BEq

structure Ann (α: Type) where
  name: String
  type: α
  deriving Repr, BEq

structure Lam (α: Type) where
  params: List (Ann α)
  body: α
  deriving Repr, BEq

structure App (α: Type) where
  cmd: α -- either Var or Lam
  args: List α
  deriving Repr, BEq

structure Case (α: Type) where
  patCmd: String
  patArgs: List String
  value: α
  deriving Repr, BEq

structure Mat (α: Type) where
  cond: α
  cases: List (Case α)
  deriving Repr, BEq

structure Hole (α: Type) where
  type: α
  deriving Repr, BEq

inductive T (α: Type) where
  -- | hole: T α -- like sorry, just to fill in the blank
  -- | trace: α → α -- print alpha
  | inh: Inh α → T α
  | typ: Typ α → T α -- typ is hole - TODO rename
  | bnd: Bnd α → T α
  | lam: Lam α → T α
  | app: App α → T α
  | mat: Mat α → T α
  deriving Repr, BEq

def T.optionMap? (t: T α) (f: α → Option β) : Option (T β) := do
  match t with
    | T.inh x =>
      let type ← f x.type
      let args ← Util.optionMap? x.args f
      T.inh {
        type := type,
        cons := x.cons,
        args := args
      }
    | T.typ x =>
      let value ← f x.value
      T.typ {
        value := value,
      }
    | T.bnd x =>
      let init ← Util.optionMap? x.init (λ bind => do
        let value ← f bind.value
        pure {
          name := bind.name,
          value := value,
          : Bind β
        }
      )
      let last ← f x.last
      T.bnd {
        init := init,
        last := last,
      }
    | T.lam x =>
      let params ← Util.optionMap? x.params (λ param => do
        let type ← f param.type
        pure {
          name := param.name,
          type := type,
          : Ann β
        }
      )
      let body ← f x.body
      T.lam {
        params := params,
        body := body,
      }
    | T.app x =>
      let cmd ← f x.cmd
      let args ← Util.optionMap? x.args f
      T.app {
        cmd := cmd,
        args := args,
      }
    | T.mat x =>
      let cond ← f x.cond
      let cases ← Util.optionMap? x.cases (λ case => do
        let value ← f case.value
        pure {
          patCmd := case.patCmd,
          patArgs := case.patArgs,
          value := value
          : Case β
        }
      )
      T.mat {
        cond := cond,
        cases := cases,
      }

inductive Term where
  | univ: (level: Int) → Term
  | var: (name: String) → Term
  | t: T Term → Term
  deriving Repr, BEq -- BEq is computationally equal == DecidableEq is logical equal = and strictly stronger than ==

def Term.map (term: Term) (f: Term → Term): Term :=
  match term with
    | .univ _ => term
    | .var _ => term
    | .t x =>
      match x.optionMap? (λ term => some (f term)) with
        | none =>
          dbg_trace "[UNREACHABLE]"
          term -- unreachable
        | some y => Term.t y

def Term.optionMap? (term: Term) (f: Term → Option Term): Option Term := do
  match term with
    | .univ _ => term
    | .var _ => term
    | .t x =>
      let y ← x.optionMap? f
      Term.t y


notation "univ" x => Term.univ x
notation "var" x => Term.var x
notation "inh" x => Term.t (T.inh x)
notation "typ" x => Term.t (T.typ x)
notation "bnd" x => Term.t (T.bnd x)
notation "lam" x => Term.t (T.lam x)
notation "app" x => Term.t (T.app x)
notation "mat" x => Term.t (T.mat x)

structure InferedTerm where
  term: Term
  type: Term
  deriving Repr

-- util
def isLam? (term: Term): Option (Lam Term) :=
  match term with
    | lam l => some l
    | _ => none

def isInh? (term: Term): Option (Inh Term) :=
  match term with
    | inh i => some i
    | _ => none

end EL2.Term
