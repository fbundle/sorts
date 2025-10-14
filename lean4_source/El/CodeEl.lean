import El.Util
import El.Code

namespace Code

inductive SortEl where -- SortEl - basic element of EL
  | int: SortEl
  | univ: Int → SortEl
  | integer: Int → SortEl
  deriving Repr

def SortEl.level (s: SortEl): Int :=
  match s with
    | int => 1
    | univ i => i
    | integer _ => 0

def SortEl.parent (s: SortEl): SortEl :=
  match s with
    | int => (.univ 2)
    | univ i => (.univ (i+1))
    | integer _ => (.int)



private def parseInteger (s: String): Option SortEl := do
  let i ← s.toInt?
  pure (.integer i) -- integer i

private def parseUniverse (s: String): Option SortEl := do
  let s ← s.dropPrefix? "U_"
  let s := s.toString
  let i ← s.toInt?
  pure (.univ i) -- universe level i

private def parseSortEl := Util.applyOnce [
  parseInteger,
  parseUniverse,
  λ _ => none,
]

def _example: List (Code SortEl) :=
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

    Util.optionMap xs (parse parseSortEl ["+"])

#eval _example


end Code
