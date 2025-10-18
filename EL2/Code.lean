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
inductive Code (β: Type) where
  | atom: β → Code β
  | var: String → Code β
  | list: List (Code β) → Code β
  | ann: Ann (Code β) → Code β
  | bind_val: BindVal (Code β) → Code β
  | bind_typ: BindTyp (Code β) → Code β
  | bind_mk: BindMk (Code β) → Code β
  | app: App (Code β) (Code β) → Code β
  | lam: Lam (Code β) → Code β
  deriving Repr



partial def Code.inferCode [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx) : Option (Code β × Ctx) := do
  -- infer: turn everything to type then normalize
  match c with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (.atom p, ctx)
    | .var n =>
      let c : Code β ← Context.get? ctx n
      c.inferCode ctx
    | _ => sorry

partial def Code.normalizeCode [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx): Option (Code β × Ctx) := do
  -- normalize: just normalize
  match c with
    | .atom a =>
      pure (c, ctx) -- return itself
    | .var n =>
      let c: Code β ← Context.get? ctx n
      c.normalizeCode ctx
    | _ => sorry


end EL2
