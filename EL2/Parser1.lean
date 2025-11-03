import Parser.Combinator
import EL2.Core

namespace EL2.Parser1
open Parser.Combinator
open EL2.Core

abbrev StringParser := Parser Char

mutual

partial def parse: StringParser Exp :=
  sorry

partial def parseUniv: StringParser Exp := λ xs => do
  let (name, rest) ← parseName xs
  if "Type".isPrefixOf name then
    let levelStr := name.stripPrefix "Type"
    let level ← levelStr.toNat?
    some (Exp.typ level, rest)
  else
    none
end
#eval parseUniv "Type123".toList

partial def parseVar: StringParser Exp := parseName.map (λ name => Exp.var name)

partial def parseType: StringParser Exp :=
  -- : X (-> X)^n

  let parseAnn: StringParser (String × Exp) :=
    (
      parseExactString "(" ++
      parseWhiteSpaceWeak ++
      parseName ++
      parseWhiteSpaceWeak ++
      parseExactString ":" ++
      parseWhiteSpaceWeak ++
      parse ++
      parseWhiteSpaceWeak ++
      parseExactString ")"
    ).map (λ (_, _, name, _, _, _, type, _, _) => (name, type))


  -- let parseX: Parser (Exp ⊕ (String × Exp)) := λ xs =>
  sorry





end EL2.Parser1
