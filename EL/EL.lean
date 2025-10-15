import EL.Form
import EL.Atom

namespace EL

partial def defaultParseAll (source: String): List (EL.Code EL.Atom) × String :=
  let tokens := Form.defaultParser.tokenize source

  let rec loop (lines: Array (EL.Code EL.Atom)) (tokens: List String) : Array (EL.Code EL.Atom) × List String :=
    match Form.defaultParser.parse tokens with
      | none => (lines, tokens)
      | some (tokens, form) =>
        match EL.parse EL.parseAtom form with
          | none => (lines, tokens)
          | some code => loop (lines.push code) tokens

  let (lines, tokens) := loop #[] tokens
  (lines.toList, String.join (tokens.intersperse " "))

end EL
