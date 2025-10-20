namespace EL2

class Irreducible β where
  inferType: β → β


class Context Ctx α where
  set: Ctx → String → α → Ctx
  get?: Ctx → String → Option α


end EL2
