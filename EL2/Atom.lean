import EL2.Class
import EL2.Util

namespace EL2

inductive Atom where -- Atom - basic element of EL
  | int_type: Atom
  | univ: (level: Int) → Atom
  | integer: (value: Int) → Atom
  | inh : (type: α) → (level: Int) → Atom

def Atom.inferAtom (s: Atom): Atom :=
  match s with
    | int_type => .univ 2
    | univ i => .univ (i+1)
    | integer _ => .int_type
    | inh type level => .inh type (level + 1)

def Atom.inhabited (type: α) (level: Int): Atom :=
  .inh type level


instance: Irreducible Atom where
  inferAtom := Atom.inferAtom
  inhabited := Atom.inhabited

instance: ToString Atom where
  toString (a: Atom): String :=
    match a with
      | .int_type => "int"
      | .univ i => s!"U_{i}"
      | .integer i => s!"{i}"
      | .inh type level => s!"inh"


end EL2
