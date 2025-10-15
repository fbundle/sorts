import EL.Form
import EL.Atom

namespace EL

def tokenize := Form.defaultParser.tokenize

def parse (tokens: List String): Option (List String × (Code Atom)) := do
  let (tokens, form) ← Form.defaultParser.parse tokens
  let code ← parseCode parseAtom form
  pure (tokens, code)


private def source := "
  (:= Nat (*U_2))
  (:= n0 (*Nat))
  (:= succ (*(-> Nat)))

  (:= n1 (succ n0))
  (:= n2 (succ n0))
  (:= x 3)
  (:= y 4)

  (+ x y)
"

#eval Util.iterateAll parse (tokenize source)

end EL
