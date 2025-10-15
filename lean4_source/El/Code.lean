import El.Form
import El.Util

abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form
abbrev Frame := Util.Frame

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

-- Irreducible β is any type β
class Irreducible (β: Type) where
  level: β → Option Int
  parent: β → Option β

-- Reducible α β is any type α that can be reduced into β
class Reducible (α: Type) (β: Type) [Irreducible β] where
  level: α → Frame α → Option Int
  parent: α → Frame α → Option β
  reduce: α → Frame α → Option β

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
  | arrow: Arrow (Code β) → Code β
  deriving Repr

partial def Code.level [Irreducible β] (c: Code β) (frame: Frame (Code β)): Option Int := do
  match c with
    | .atom a => Irreducible.level a -- somehow, just a.level does not work
    | .name n =>
      let a ← frame.get? n
      a.level frame
    | _ => sorry -- TODO

instance [Irreducible β]: Reducible (Code β) β where
  level (c: Code β) (frame: Frame (Code β)): Option Int := c.level frame
  parent (c: Code β) (frame: Frame (Code β)): Option β := sorry -- equivalent to typecheck
  reduce (c: Code β) (frame: Frame (Code β)): Option β := sorry -- equivalent to execute



end Code
