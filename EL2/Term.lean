namespace EL2

structure Inh (α: Type) where
  type: α
  cons: String
  args: List α
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
  type: α
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
  -- | hole: T α -- sorry
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


end EL2
