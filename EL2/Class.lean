namespace EL2

class Irreducible β where
  inferAtom: β → β
  inhabited: α → Int → β


class Context Ctx α where
  set: Ctx → String → α → (Ctx × α)
  get?: Ctx → String → Option α


end EL2
