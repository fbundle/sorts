namespace EL2


class Context Ctx α where
  set: Ctx → String → α → Ctx
  get?: Ctx → String → Option α


end EL2
