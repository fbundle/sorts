namespace EL

-- Context Ctx α is any persistent map Ctx of key String and val α
class Context (Ctx: Type) (α: Type) where
  get?: Ctx → String → Option α
  set: Ctx → String → α → Ctx

-- Irreducible β is any type β
class Irreducible (β: Type) where
  infer: β → β

-- Reducible α β is any type α that can be reduced into β
class Reducible (α: Type) (Ctx: Type) [Context Ctx α] where
  infer: α → Ctx → Option (α × Ctx)
  normalize: α → Ctx → Option (α × Ctx)

end EL
