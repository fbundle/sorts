import EL2.Class

namespace EL2

structure Ann (α: Type) where -- (2: Nat)
  name: String
  type: α
  deriving Repr

structure BindVal (α: Type) where
  name: String
  value: α
  deriving Repr

-- BindTyp : type
structure BindTyp (α: Type) where -- List (T: Type)
  name: String
  params: List (Ann α)
  parent: α
  deriving Repr

-- App : function application
structure App (α: Type) (β: Type) where
  cmd: α
  args: List β
  deriving Repr

-- BindMk : type constructor
structure BindMk (α: Type) where  -- nil: List T or cons (init: List T) (tail: T): List T
  name: String
  params: List (Ann α)            -- (init: List T) (tail: T)
  type: App String α                 -- (List T)
  deriving Repr

-- Lam : function abstraction
structure Lam (α: Type) where
  params: List (Ann α)
  body: α
  deriving Repr

-- TODO - add Mat - match

-- β is an atomic type which is reduced into itself, e.g. integer
-- it instantiates Reducible β β
-- Code β is any type which can be reduced into β
-- it instantiates Reducible (Code β) β
-- it is usually denoted by α
inductive Term (β: Type) where
  | atom: β → Term β
  | var: String → Term β
  | list: List (Term β) → Term β
  | ann: Ann (Term β) → Term β
  | bind_val: BindVal (Term β) → Term β
  | bind_typ: BindTyp (Term β) → Term β
  | bind_mk: BindMk (Term β) → Term β
  | app: App (Term β) (Term β) → Term β
  | lam: Lam (Term β) → Term β
  deriving Repr



end EL2
