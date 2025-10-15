import El.Util
import El.Code
import El.CodeParse

namespace Code

inductive Atom where -- Atom - basic element of EL
  | int: Atom
  | univ: Int → Atom
  | integer: Int → Atom
  deriving Repr

def Atom.level (s: Atom) (frame: Frame Atom): Option Int :=
  match s with
    | int => some 1
    | univ i => some i
    | integer _ => some 0

def Atom.parent (s: Atom) (frame: Frame Atom): Option Atom :=
  match s with
    | int => some (.univ 2)
    | univ i => some (.univ (i+1))
    | integer _ => some (.int)

instance: Reducible Atom Atom where
  level (s: Atom) (frame: Frame Atom): Option Int := s.level frame
  parent (s: Atom) (frame: Frame Atom): Option Atom := s.parent frame
  reduce (s: Atom) (frame: Frame Atom): Option Atom := s

private def parseInteger (s: String): Option Atom := do
  let i ← s.toInt?
  pure (.integer i) -- integer i

private def parseUniverse (s: String): Option Atom := do
  let s ← s.dropPrefix? "U_"
  let s := s.toString
  let i ← s.toInt?
  pure (.univ i) -- universe level i

private def parseAtom := Util.applyOnce [
  parseInteger,
  parseUniverse,
  λ _ => none,
]

def _example: List (Code Atom) :=
  let source := "
    (:= Nat (*U_2))
    (:= n0 (*Nat))
    (:= succ (*(-> Nat)))

    (:= n1 (succ n0))
    (:= n2 (succ n0))
    (:= x 3)
    (:= y 4)

    (+ x y)
  "
  match Form.defaultParseAll source with
    | none => []
    | some xs =>

    Util.optionMap xs (parse parseAtom ["+"])

#eval _example


end Code
