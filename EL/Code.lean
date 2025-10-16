namespace EL

structure Beta (α: Type) where
  cmd: α
  args: List α
  deriving Repr

structure Annot (α: Type) where
  name: String
  type: α
  deriving Repr

structure Binding (α: Type) where
  name: String
  value: α
  deriving Repr

structure Typeof (α: Type) where
  value: α
  deriving Repr

structure Pi (α: Type) where -- Pi or Lambda
  params: List (Annot α)
  body: α
  deriving Repr

structure Ind (α: Type) where -- Inductive
  name: Annot α
  cons: List (Annot (Pi α))
  deriving Repr

structure Case (α: Type) where
  pattern: Beta String
  value: α
  deriving Repr

structure Mat (α: Type) where
  cond: α
  cases: List (Case α)
  deriving Repr


-- β is an atomic type which is reduced into itself, e.g. integer
-- it instantiates Reducible β β
-- Code β is any type which can be reduced into β
-- it instantiates Reducible (Code β) β
-- it is usually denoted by α
inductive Code (β: Type) where
  | atom: β → Code β
  | name: String → Code β
  | beta: Beta (Code β) → Code β
  | binding: Binding (Code β) → Code β
  | typeof: Typeof (Code β) → Code β
  | pi: Pi (Code β) → Code β
  | ind: Ind (Code β) → Code β
  | mat: Mat (Code β) → Code β
  deriving Repr

structure Context α where
  set: String → α → Context α
  get?: String → Option α

partial def Code.infer (c: Code β) (ctx: Context (Code β)) (inferAtom: β → β): Option (Code β × Context (Code β)) := do
  -- infer: turn everything to type then normalize
  match c with
    | .atom a =>
      let p ← inferAtom a
      pure (.atom p, ctx)
    | .name n =>
      let c ← ctx.get? n
      c.infer ctx inferAtom
    | _ => sorry

partial def Code.normalize (c: Code β) (ctx: Context (Code β)): Option (Code β × Context (Code β)) := do
  -- normalize: just normalize
  match c with
    | .atom a =>
      pure (c, ctx) -- return itself
    | .name n =>
      let c ← ctx.get? n
      c.normalize ctx
    | _ => sorry


end EL
