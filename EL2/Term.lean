import EL2.Class

namespace EL2

structure Ann (α: Type) where -- (2: Nat)
  name: String
  type: α

structure BindVal (α: Type) where
  name: String
  value: α

-- BindTyp : type
structure BindTyp (α: Type) where -- List (T: Type)
  name: String
  params: List (Ann α)
  parent: α

-- App : function application
structure App (α: Type) (β: Type) where
  cmd: α
  args: List β

-- BindMk : type constructor
structure BindMk (α: Type) where  -- nil: List T or cons (init: List T) (tail: T): List T
  name: String
  params: List (Ann α)            -- (init: List T) (tail: T)
  type: App String α                 -- (List T)

-- Lam : function abstraction
structure Lam (α: Type) where
  params: List (Ann α)
  body: α

structure Case (α: Type) where
  pattern: App String String
  value: α

-- Mat : match
structure Mat (α: Type) where
  cond: α
  cases: List (Case α)

-- Lst : an non empty list
structure Lst (α: Type) where
  init: List α
  tail: α

inductive T (α: Type) where
  | var: (name: String) → T α
  | lst: Lst α → T α
  | bind_val: BindVal α → T α
  | bind_typ: BindTyp α → T α
  | bind_mk: BindMk α → T α
  | lam: Lam α → T α
  | app: App α α → T α
  | mat: Mat α → T α

inductive Term (β: Type) where
  | atom: (a: β) → Term β
  | t: T (Term β) → Term β

inductive InhTerm (β: Type) where
  | atom: (a: β) → InhTerm β
  | t: (value: T (InhTerm β)) → (type: T (InhTerm β)) → InhTerm β

instance: Coe (T (Term β)) (Term β) where
  coe := (.t ·)

end EL2
