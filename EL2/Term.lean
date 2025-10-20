namespace EL2

structure Ann (α: Type) where -- (2: Nat)
  name: String
  type: α
  deriving Repr, BEq

structure BindVal (α: Type) where
  name: String
  value: α
  deriving Repr, BEq

-- BindTyp : type
structure BindTyp (α: Type) where -- List (T: Type)
  name: String
  params: List (Ann α)
  level: Int
  deriving Repr, BEq

-- App : function application
structure App (α: Type) (β: Type) where
  cmd: α
  args: List β
  deriving Repr, BEq

-- BindMk : type constructor
structure BindMk (α: Type) where  -- nil: List T or cons (init: List T) (tail: T): List T
  name: String
  params: List (Ann α)            -- (init: List T) (tail: T)
  type: App String α                 -- (List T)
  deriving Repr, BEq

-- Lam : function abstraction
structure Lam (α: Type) where
  params: List (Ann α)
  body: α
  deriving Repr, BEq

structure Case (α: Type) where
  pattern: App String String
  value: α
  deriving Repr, BEq

-- Mat : match
structure Mat (α: Type) where
  cond: α
  cases: List (Case α)
  deriving Repr, BEq

-- Lst : an non empty list
structure Lst (α: Type) where
  init: List α
  last: α
  deriving Repr, BEq

-- Typ
structure Typ (α: Type) where
  value: α
  deriving Repr, BEq

inductive Term where
  | univ: (level: Int) → Term
  | var: (name: String) → Term
  | lst: Lst Term → Term
  | bind_val: BindVal Term → Term
  | bind_typ: BindTyp Term → Term
  | bind_mk: BindMk Term → Term
  | typ: Typ Term → Term
  | lam: Lam Term → Term
  | app: App Term Term → Term
  | mat: Mat Term → Term
  deriving Repr, BEq -- BEq is computationally equal == DecidableEq is logical equal = and strictly stronger than ==

notation "atom" x => Term.atom x
notation "univ" x => Term.univ x
notation "var" x => Term.var x
notation "lst" x => Term.lst x
notation "bind_typ" x => Term.bind_typ x
notation "bind_val" x => Term.bind_val x
notation "bind_mk" x => Term.bind_mk x
notation "typ" x => Term.typ x
notation "lam" x => Term.lam x
notation "app" x => Term.app x
notation "mat" x => Term.mat x


end EL2
