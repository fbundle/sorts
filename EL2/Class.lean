namespace EL2

class Irreducible β where
  inferAtom: β → β
  inhabited: (type: α) → (level: Int) → β -- useful to create dummy variable for lambda type checking


class Context Ctx α where
  set: Ctx → String → α → (Ctx × α)
  get?: Ctx → String → Option α


end EL2
