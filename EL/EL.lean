import EL.Form
import EL.Atom

namespace EL

def tokenize := Form.defaultParser.tokenize

def parse (tokens: List String): Option (List String × (Code Atom)) := do
  let (tokens, form) ← Form.defaultParser.parse tokens
  let code ← parseCode parseAtom form
  pure (tokens, code)

end EL
