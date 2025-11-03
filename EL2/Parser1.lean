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

    || parse.map (λ e => ("_", e))

  let parseArrowAnn : StringParser (String × Exp) :=
    -- -> X
    (
      parseExactString "->" ++
      parseWhiteSpaceWeak ++
      parseAnn
    ).map (λ (_, _, x) => x)

  let aaa := (
    parseExactString ":" ++
    parseWhiteSpaceWeak ++
    parseAnn ++
    (parseWhiteSpaceWeak ++ parseArrowAnn ++ parseWhiteSpaceWeak).list
  ).map (λ (_, _, ann1, ann2s) =>
    let ann2s: List (String × Exp) := (ann2s.map (λ ((_, ann2s, _): String × (String × Exp) × String) =>
      ann2s
    ))
    let anns := ann1 :: ann2s

    let init := anns.extract 0 (anns.length - 1)
    -- get last elem using construction of anns
    let last := anns.getLast (List.cons_ne_nil ann1 ann2s)

    let rec loop (lastExp: Exp) (anns: List (String × Exp)): Exp :=
      match anns with
        | [] => lastExp
        | (name, type) :: rest =>
          loop (Exp.pi name type lastExp) rest

    loop last.snd init

  )




  sorry





end EL2.Parser1
