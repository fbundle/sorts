import El.Form
import El.Util

abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form

namespace Code

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

structure Arrow (α: Type) where
  a: α
  b: α
  deriving Repr

class AtomClass (β: Type) where
  level: β → Int
  parent: β → β

inductive Code (β: Type) [AtomClass β] where
  | atom: β → Code β
  | name: String → Code β
  | beta: Beta (Code β) → Code β
  | other: Other (Code β) → Code β
  | annot: Annot (Code β) → Code β
  | binding: Binding (Code β) → Code β
  | typeof: Typeof (Code β) → Code β
  | inh: Inh (Code β) → Code β
  | pi: Pi (Code β) → Code β
  | arrow: Arrow (Code β) → Code β
  deriving Repr

instance [AtomClass β]: AtomClass (Code β) where
  level (s: Code β): Int := sorry
  parent (s: Code β): Code β := sorry

end Code
