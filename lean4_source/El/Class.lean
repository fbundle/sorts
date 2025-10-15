namespace EL

-- Context Ctx α is any persistent map Ctx of key String and val α
class Context (Ctx: Type) (α: Type) where
  get?: Ctx → String → Option α
  set: Ctx → String → α → Ctx

-- Irreducible β is any type β
class Irreducible (β: Type) where
  level: β → Option Int
  parent: β → Option β

-- Reducible α β is any type α that can be reduced into β
class Reducible (α: Type) (β: Type) (Ctx: Type) [Irreducible β] [Context Ctx α] where
  level: α → Ctx → Option Int
  parent: α → Ctx → Option β
  reduce: α → Ctx → Option β

end EL
