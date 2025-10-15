import EL.Form
import EL.Class
import EL.Util


namespace EL

abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form
abbrev Frame := Util.Frame

structure Beta (α: Type) where
  cmd: α
  args: List α
  deriving Repr

structure Other (α: Type) where -- Other - any form
  head: String
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

structure Inh (α: Type) where -- Inhabited
  type: α
  deriving Repr

structure Pi (α: Type) where -- Pi or Lambda
  params: List (Annot α)
  body: α
  deriving Repr

-- β is an atomic type which is reduced into itself, e.g. integer
-- it instantiates Reducible β β
-- Code β is any type which can be reduced into β
-- it instantiates Reducible (Code β) β
-- it is usually denoted by α
inductive Code (β: Type) [Irreducible β] where
  | atom: β → Code β
  | name: String → Code β
  | beta: Beta (Code β) → Code β
  | other: Other (Code β) → Code β
  | annot: Annot (Code β) → Code β
  | binding: Binding (Code β) → Code β
  | typeof: Typeof (Code β) → Code β
  | inh: Inh (Code β) → Code β
  | pi: Pi (Code β) → Code β
  deriving Repr

partial def Code.parent [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx): Option (Code β × Ctx) := do
  match c with
    | .atom a =>
      let p ← Irreducible.parent (β := β) a
      pure (.atom p, ctx)
    | .name n =>
      let c ← Context.get? (α := Code β) ctx n
      c.parent ctx
    | _ => sorry

partial def Code.reduce [Irreducible β] [Context Ctx (Code β)] (c: Code β) (ctx: Ctx): Option (Code β × Ctx) := do
    match c with
    | .atom a =>
      pure (c, ctx) -- return itself
    | .name n =>
      let c ← Context.get? (α := Code β) ctx n
      c.reduce ctx
    | _ => sorry

instance [Irreducible β] [Context Ctx (Code β)]: Reducible (Code β) Ctx where
  parent := Code.parent
  reduce := Code.reduce


end EL
