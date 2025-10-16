import EL.EL
import Std


def main (args : List String) : IO UInt32 := do
  match args with
  | [fileName] => do
      let content ← IO.FS.readFile fileName
      let tokens := EL.tokenize content
      let result := Util.parseAll EL.parse tokens
      if result.remaining.length ≠ 0 then
        let remaining := String.join (result.remaining.intersperse " ")
        IO.println s!"{repr result.items}\nerror at {remaining}"
        return 1
      else
        IO.println s!"{repr result.items}"
        return 0
  | _ => do
      IO.eprintln "Usage: el <file>"
      return 1
