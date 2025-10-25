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
  | inh: Inh α → T α
  | typ: Typ α → T α -- TODO drop
  | bnd: Bnd α → T α
  | lam: Lam α → T α
  | app: App α → T α
  | mat: Mat α → T α
  deriving Repr, BEq

inductive Term where
  | univ: (level: Int) → Term
  | var: (name: String) → Term
  | t: T Term → Term
  deriving Repr, BEq -- BEq is computationally equal == DecidableEq is logical equal = and strictly stronger than ==

notation "univ" x => Term.univ x
notation "var" x => Term.var x
notation "inh" x => Term.t (T.inh x)
notation "typ" x => Term.t (T.typ x)
notation "bnd" x => Term.t (T.bnd x)
notation "lam" x => Term.t (T.lam x)
notation "app" x => Term.t (T.app x)
notation "mat" x => Term.t (T.mat x)

inductive ReducedTerm where
  | univ: (level: Int) → ReducedTerm
  | t: T ReducedTerm → ReducedTerm
  deriving Repr, BEq

notation "r_univ" x => ReducedTerm.univ x
notation "r_inh" x => ReducedTerm.t (T.inh x)
notation "r_typ" x => ReducedTerm.t (T.typ x)
notation "r_bnd" x => ReducedTerm.t (T.bnd x)
notation "r_lam" x => ReducedTerm.t (T.lam x)
notation "r_app" x => ReducedTerm.t (T.app x)
notation "r_mat" x => ReducedTerm.t (T.mat x)

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
