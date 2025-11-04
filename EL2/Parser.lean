import Parser.Combinator
import EL2.Core

namespace EL2.Parser.Internal
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

def parseName: Parser Char String :=
  String.toStringParser $ (pred ("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_.".contains ·)).many1


mutual

partial def parseApp: Parser Char Exp :=
  -- parse any thing starts with (
  (
    String.exact "(" ++
    String.whitespaceWeak ++
    parse ++
    (String.whitespace ++ parse).many ++
    String.whitespaceWeak ++
    String.exact ")"
  ).filterMap (λ (_, _, cmd, args, _, _) =>
    chainCmd cmd (args.map Prod.snd)
  )

partial def parseHom: Parser Char Exp :=
  let parseAnn: Parser Char (String × Exp) :=
      (
        String.exact "(" ++
        String.whitespaceWeak ++
        parseName ++
        String.whitespaceWeak ++
        String.exact ":" ++
        String.whitespaceWeak ++
        parse ++
        String.whitespaceWeak ++
        String.exact ")"
      ).map (λ (_, _, name, _, _, _, type, _, _) => (name, type))

      || parse.map (λ e => ("_", e))

  (
    String.exact "hom" ++
    (String.whitespaceWeak ++ parseAnn).many ++
    String.whitespaceWeak ++
    String.exact "->" ++
    String.whitespaceWeak ++
    parse
  ).map (λ (_, params, _, _, _, typeB) =>
    chainPi (params.map (λ (_, name, typeA) => (name, typeA))) typeB
  )

partial def parseLam: Parser Char Exp :=
  -- parse anything starts with lam
  -- lam name [ name]^n => body
  (
    String.exact "lam" ++
    (String.whitespace ++ parseName).many ++
    String.whitespace ++
    String.exact "=>" ++
    String.whitespace ++
    parse
  ).map (λ (_, names, _, _, _, body) =>
    chainLam (names.map Prod.snd) body
  )

partial def parseUniv: Parser Char Exp := λ xs => do
  let (name, rest) ← parseName xs
  if "Type".isPrefixOf name then
    let levelStr := name.stripPrefix "Type"
    match levelStr.toNat? with
      | none => err xs
      | some level =>
        Except.ok (Exp.typ level, rest)
  else
    err xs

partial def parseVar: Parser Char Exp := parseName.filterMap (λ name =>
  let specialNames := [
    "lam", "let", "inh", "hom"
  ]

  if specialNames.contains name then
    none
  else
    some (Exp.var name)
)


partial def parseBnd: Parser Char Exp :=
  -- parse anything starts with let
  -- typed let
  (
    String.exact "let" ++
    String.whitespaceWeak ++
    parseName ++
    String.whitespaceWeak ++
    String.exact ":" ++
    String.whitespaceWeak ++
    parse ++
    String.whitespaceWeak ++
    String.exact ":=" ++
    String.whitespaceWeak ++
    parse ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, _, _, type, _, _, _, value, _, body) =>
    Exp.bnd name value type body
  )
  ||

  -- untyped let
  (
    String.exact "let" ++
    String.whitespaceWeak ++
    parseName ++
    String.whitespaceWeak ++
    String.exact ":=" ++
    String.whitespaceWeak ++
    parse ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, _, _, value, _, body) =>
    Exp.app (Exp.lam name body) value
  )

partial def parseInh: Parser Char Exp :=
  -- parse anything starts with inh
  (
    String.exact "inh" ++
    String.whitespaceWeak ++
    parseName ++
    String.whitespaceWeak ++
    String.exact ":" ++
    String.whitespaceWeak ++
    parse ++
    parseLineBreak ++
    parse
  ).map (λ (_, _, name, _, _, _, type, _, body) =>
    Exp.inh name type body
  )

partial def parse: Parser Char Exp :=
  (
    parseUniv ||-- starts with Type
    parseApp || -- starts with (
    parseLam || -- starts with lam
    parseHom || -- starts with hom
    parseBnd || -- starts with let
    parseInh || -- starts with inh
    parseVar    -- everything else
  )

end


end EL2.Parser.Internal

namespace EL2.Parser
def parse: Parser.Combinator.Parser Char EL2.Core.Exp :=
  (
    Parser.Combinator.String.whitespaceWeak ++
    Internal.parse ++
    Parser.Combinator.String.whitespaceWeak
  ).map (λ (_, e, _) => e)




#eval parse "
  inh Nat_rec :
    (P : Nat -> Type0) ->
    (P zero) ->
    ((n : Nat) -> (P n) -> (P (succ n))) ->
    (n : Nat) -> (P n)
body
".toList



end EL2.Parser
