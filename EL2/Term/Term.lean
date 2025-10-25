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
  cmd: α
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

inductive T (α: Type) where
  -- | hole: T α -- like sorry, just to fill in the blank
  -- | trace: α → α -- print alpha
  | univ: (level: Int) → T α
  | var: (name: String) → T α
  | inh: Inh α → T α
  | typ: Typ α → T α
  | bnd: Bnd α → T α
  | lam: Lam α → T α
  | app: App α → T α
  | mat: Mat α → T α
  deriving Repr, BEq

inductive Term where
  | t: T Term → Term
  deriving Repr, BEq -- BEq is computationally equal == DecidableEq is logical equal = and strictly stronger than ==

notation "univ" x => Term.t (T.univ x)
notation "var" x => Term.t (T.var x)
notation "inh" x => Term.t (T.inh x)
notation "typ" x => Term.t (T.typ x)
notation "bnd" x => Term.t (T.bnd x)
notation "lam" x => Term.t (T.lam x)
notation "app" x => Term.t (T.app x)
notation "mat" x => Term.t (T.mat x)


class NameMap M α where
  size: M → Nat
  set: M → String → α → M
  get?: M → String → Option α

structure InferedTerm where
  term: Term
  type: Term
  deriving Repr

end EL2.Term
