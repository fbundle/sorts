import Parser.Combinator
import EL2.Core

namespace EL2.Parser1
open Parser.Combinator
open EL2.Core

mutual

partial def parse: Parser Char Exp :=
  sorry

partial def parseUniv: Parser Char Exp := λ xs => do
  let (name, rest) ← String.name xs
  if "Type".isPrefixOf name then
    let levelStr := name.stripPrefix "Type"
    let level ← levelStr.toNat?
    some (Exp.typ level, rest)
  else
    none
end
#eval parseUniv "Type123".toList

partial def parseVar: Parser Char Exp := String.name.map (λ name => Exp.var name)

partial def parseColonType: Parser Char Exp :=
  -- : X (-> X)^n for some n ≥ 0
  let parseAnn: Parser Char (String × Exp) :=
    (
      String.exact "(" ++
      String.whitespaceWeak ++
      String.name ++
      String.whitespaceWeak ++
      String.exact ":" ++
      String.whitespaceWeak ++
      parse ++
      String.whitespaceWeak ++
      String.exact ")"
    ).map (λ (_, _, name, _, _, _, type, _, _) => (name, type))

    || parse.map (λ e => ("_", e))

  let parseArrowAnn : Parser Char (String × Exp) :=
    -- -> X
    (
      String.exact "->" ++
      String.whitespaceWeak ++
      parseAnn
    ).map (λ (_, _, x) => x)

  -- colon then type
  (
    String.exact ":" ++
    String.whitespaceWeak ++
    parseAnn ++
    (String.whitespaceWeak ++ parseArrowAnn).list
  ).map (λ (_, _, ann1, ann2s) =>
    let ann2s: List (String × Exp) := (ann2s.map (λ ((_, ann2s): String × (String × Exp)) =>
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

def parseLam: Parser Char Exp :=
  (
    (String.exact "λ" || String.exact "lam") ++
    String.whitespaceWeak ++
    String.name ++
    String.whitespaceWeak ++
    String.exact "=>" ++
    parse
  ).map (λ (_, _, name, _, _, body) => Exp.lam name body)

def parseLineBreak :=
  -- <whitespace_without_newline> <newline> <writespace>
  String.whiteSpaceWithoutNewLineWeak ++
  (String.exact "\n" || String.exact ";") ++
  String.whitespaceWeak

def parseBnd: Parser Char Exp :=
  (
    String.exact "let" ++
    String.whitespaceWeak ++
    String.name ++
    String.whitespaceWeak ++
    parseColonType ++
    String.whitespaceWeak ++
    String.exact ":=" ++
    String.whitespaceWeak ++
    parse ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, type, _, _, _, value, _, body) =>
    Exp.bnd name value type body
  )

def parseInh: Parser Char Exp :=
  (
    String.exact "inh" ++
    String.whitespaceWeak ++
    String.name ++
    String.whitespaceWeak ++
    parseColonType ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, type, _, body) =>
    Exp.inh name type body
  )

end EL2.Parser1
