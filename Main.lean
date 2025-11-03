import EL2.Parser

open EL2.Parser

private def s := "
inh Nat : Type0
body
"

private def tokens := tokenize s

def main  : IO Unit := do
  IO.println tokens
  IO.println "--------------------------------------"
  match parse tokens with
    | none => return
    | some (rest, e) =>
      IO.println e
