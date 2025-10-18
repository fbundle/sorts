import EL2.Class
import EL2.Util

namespace EL

inductive Atom where -- Atom - basic element of EL
  | int_type: Atom
  | univ: Int → Atom
  | integer: Int → Atom
  deriving Repr

def Atom.inferAtom (s: Atom): Atom :=
  match s with
    | int_type => .univ 2
    | univ i => .univ (i+1)
    | integer _ => .int_type

instance: Irreducible Atom where
  inferAtom := Atom.inferAtom


end EL
