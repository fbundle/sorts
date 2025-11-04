import Parser.Combinator
import EL2.Core

namespace EL2.Parser1
open Parser.Combinator
open EL2.Core


def parseLineBreak :=
  -- <whitespace_without_newline> <newline> <writespace>
  String.whiteSpaceWithoutNewLineWeak ++
  (String.exact "\n" || String.exact ";") ++
  String.whitespaceWeak

def chainCmd (cmd: Exp) (args: List Exp): Exp :=
  match args with
    | [] => cmd
    | arg :: args =>
      chainCmd (Exp.app cmd arg) args

def chainPi (anns: List (String × Exp)) (last: Exp): Exp :=
  match anns with
    | [] => last
    | (name, type) :: anns =>
      Exp.pi name type (chainPi anns last)

def chainLam (names: List String) (body: Exp): Exp :=
  match names with
    | [] => body
    | name :: names =>
      Exp.lam name (chainLam names body)

mutual

partial def parse: Parser Char Exp :=
  parseUniv ||
  parseVar

partial def parseApp: Parser Char Exp :=
  (
    String.exact "(" ++
    String.whitespaceWeak ++
    parse.list ++
    String.whitespaceWeak ++
    String.exact ")"
  ).filterMap (λ (_, _, es, _, _) =>
    match es with
      | [] => none
      | cmd :: args => some (chainCmd cmd args)
  )


partial def parseUniv: Parser Char Exp := λ xs => do
  let (name, rest) ← String.name xs
  if "Type".isPrefixOf name then
    let levelStr := name.stripPrefix "Type"
    let level ← levelStr.toNat?
    some (Exp.typ level, rest)
  else
    none

partial def parseVar: Parser Char Exp := String.name.map (λ name => Exp.var name)

partial def parseColonArrow: Parser Char Exp :=
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

    chainPi init last.snd
  )

partial def parseLam: Parser Char Exp :=
  (
    (String.exact "λ" || String.exact "lam") ++
    String.whitespaceWeak ++
    String.name.list ++
    String.whitespaceWeak ++
    String.exact "=>" ++
    parse
  ).map (λ (_, _, names, _, _, body) => chainLam names body)



partial def parseBnd: Parser Char Exp :=
  (
    String.exact "let" ++
    String.whitespaceWeak ++
    String.name ++
    String.whitespaceWeak ++
    parseColonArrow ++
    String.whitespaceWeak ++
    String.exact ":=" ++
    String.whitespaceWeak ++
    parse ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, type, _, _, _, value, _, body) =>
    Exp.bnd name value type body
  )

partial def parseInh: Parser Char Exp :=
  (
    String.exact "inh" ++
    String.whitespaceWeak ++
    String.name ++
    String.whitespaceWeak ++
    parseColonArrow ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, type, _, body) =>
    Exp.inh name type body
  )
end



end EL2.Parser1
