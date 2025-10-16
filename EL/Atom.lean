import EL.Util
import EL.Parse

namespace EL

inductive Atom where -- Atom - basic element of EL
  | int: Atom
  | univ: Int → Atom
  | integer: Int → Atom
  deriving Repr

def Atom.inferAtom (s: Atom): Atom :=
  match s with
    | int => .univ 2
    | univ i => .univ (i+1)
    | integer _ => .int

def parseInteger (s: String): Option Atom := do
  let i ← s.toInt?
  pure (.integer i) -- integer i

def parseUniverse (s: String): Option Atom := do
  let s ← s.dropPrefix? "U_"
  let s := s.toString
  let i ← s.toInt?
  pure (.univ i) -- universe level i

def parseAtom := Util.applyOnce [
  parseInteger,
  parseUniverse,
  λ _ => none,
]


end EL
