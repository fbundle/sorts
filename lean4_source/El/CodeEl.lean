import El.Util
import El.Code

namespace Code


inductive Atom where
  | univ: Int → Atom
  | integer: Int → Atom
  deriving Repr

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
