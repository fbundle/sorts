import EL.Class

namespace EL

structure Beta (α: Type) where
  cmd: α
  args: List α
  deriving Repr

structure Annot (α: Type) (β: Type) where
  left: α
  right: β
  deriving Repr

structure Binding (α: Type) where
  name: String
  value: α
  deriving Repr

structure Infer (α: Type) where -- Type of
  value: α
  deriving Repr

structure Pi (α: Type) (β: Type) where -- Pi or Lambda
  params: List (Annot String α)
  body: β
  deriving Repr

structure Ind (α: Type) where -- Inductive
  name: Annot (String ⊕ Pi α (Beta String)) α
  cons: List (Annot String (String ⊕ Pi α (String ⊕ Beta String)))
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
  | infer: Infer (Code β) → Code β
  | pi: Pi (Code β) (Code β) → Code β
  | ind: Ind (Code β) → Code β
  | ind_dep: IndDep (Code β)→ Code β
  | mat: Mat (Code β) → Code β
  deriving Repr

partial def Code.inferCode [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx) : Option (Code β × Ctx) := do
  -- infer: turn everything to type then normalize
  match c with
    | .atom a =>
      let p := Irreducible.inferAtom a
      pure (.atom p, ctx)
    | .name n =>
      let c : Code β ← Context.get? ctx n
      c.inferCode ctx
    | _ => sorry

partial def Code.normalizeCode [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx): Option (Code β × Ctx) := do
  -- normalize: just normalize
  match c with
    | .atom a =>
      pure (c, ctx) -- return itself
    | .name n =>
      let c: Code β ← Context.get? ctx n
      c.normalizeCode ctx
    | _ => sorry


end EL
