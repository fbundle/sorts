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

-- Reducible α β is any type α that can be reduced into β
class Reducible (α: Type) (β: Type) where
  level: α → Int
  parent: α → β
  reduce: α → β

-- β is an atomic type which is reduced into itself, e.g. integer
-- it instantiates Reducible β β
-- Code β is any type which can be reduced into β
-- it instantiates Reducible (Code β) β
-- it is usually denoted by α
inductive Code (β: Type) [Reducible β β] where
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

def Code.level [Reducible β β] (c: Code β): Int :=
  match c with
    | .atom a => Reducible.level (α := β) (β := β) a -- somehow, just a.level does not work
    | _ => sorry -- TODO

instance [Reducible β β]: Reducible (Code β) β where
  level (s: Code β): Int := sorry
  parent (s: Code β): β := sorry -- equivalent to typecheck
  reduce (s: Code β): β := sorry -- equivalent to execute



end Code
