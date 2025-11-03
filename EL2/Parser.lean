import EL2.Core

namespace EL2.Parser.Internal
-- TOKENIZER
def sortSplitTokens (splitTokens : List String) : List String :=
  -- sort tokens so that if s2 is a prefix of s1, s1 should come first
  let lessEqual (s1: String) (s2: String): Bool :=
    if (s2.length < s1.length) && (s2.isPrefixOf s1) then true else
    if (s1.length < s2.length) && (s1.isPrefixOf s2) then false else
    s1 < s2

  splitTokens.mergeSort lessEqual

partial def stringIndexOf? (s: String) (substring: String): Option Nat :=
  -- return the starting position of substring in s if exists
  if substring.isEmpty then
    some 0
  else
    let rec helper (s: String) (substring: String) (startIdx: Nat) : Option Nat :=
      if startIdx + substring.length > s.length then
        none
      else if s.extract {byteIdx := startIdx} {byteIdx := startIdx + substring.length} = substring then
        some startIdx
      else
        helper s substring (startIdx + 1)
    helper s substring 0

partial def splitPart (sortedSplitTokens : List String) (part : String) : List String :=
  match sortedSplitTokens with
    | [] => [part]
    | s :: ss =>
      match stringIndexOf? part s with
        | some n =>
          let before := part.take n
          let after := part.drop (n + s.length)
          let beforeParts := if before.isEmpty then [] else splitPart sortedSplitTokens before
          let afterParts := if after.isEmpty then [] else splitPart sortedSplitTokens after
          beforeParts ++ [s] ++ afterParts
        | none => splitPart ss part

def tokenize (sortedSplitTokens: List String) (s: String) : List String :=
   -- remove comments
  let lines := s.split (· = '\n')
  let lines := lines.map (λ line =>
    match line.splitOn "--" with
      | head :: _ => head -- take everything before --
      | _ => line
  )
  let s := String.join (lines.intersperse "\n")

  -- basic tokenize
  let parts := s.split (·.isWhitespace)

  -- further tokenize by splitTokens
  let output := parts.flatMap (splitPart sortedSplitTokens)

  -- drop empty tokens
  let output := output.filter (·.length > 0)

  output

def newTokenizer (splitTokens: List String): String → List String :=
  tokenize (sortSplitTokens splitTokens)

-- PARSER COMBINATOR

def Parser α := List String → Option ((List String) × α)

def Parser.mapPartial  (p: Parser α) (f: α → Option β): Parser β := λ tokens => do
  let (tokens, a) ← p tokens
  let b ← f a
  pure (tokens, b)

def Parser.map (p: Parser α) (f: α → β): Parser β := p.mapPartial (λ a => some (f a))

def Parser.concat (p1: Parser α) (p2: Parser β): Parser (α × β) := λ tokens => do
  let (tokens, a) ← p1 tokens
  let (tokens, b) ← p2 tokens
  (tokens, (a, b))

infixr: 60 " * " => Parser.concat

def Parser.either (p1: Parser α) (p2: Parser α): Parser α := λ tokens => do
  match p1 tokens with
    | some (tokens, a) => some (tokens, a)
    | none => p2 tokens

infixr: 50 " || " => Parser.either -- lower precedence than concat

open EL2.Core

def parseString: Parser String := λ tokens =>
  match tokens with
    | [] => none
    | head :: tokens =>
      (tokens, head)

def parseSingle (convert?: String → Option α): Parser α := parseString.mapPartial convert?

def predToOption (f: α → Bool): α → Option α := λ a =>
  if f a then some a else none

def parseExact (pattern: String): Parser String := parseSingle $ predToOption (· = pattern)

-- parse Exp
def parseTyp: Parser Exp := parseSingle (λ head => do
  if ¬ "Type".isPrefixOf head then none else
  let levelStr := head.stripPrefix "Type"
  let level ← levelStr.toNat?
  pure (Exp.typ level)
)

def parseVar: Parser Exp := parseString.map (Exp.var ·)

-- def parseApp:

def parseAnn (parseExp: Parser Exp): Parser (String × Exp) :=
  (
    parseString *
    (parseExact ":") *
    parseExp
  ).map (λ (name, _, type) => (name, type))

def parsePi (parseExp: Parser Exp): Parser Exp :=
  -- named Pi or unnamed Pi

  (
    parseAnn parseExp *
    parseExact "->" *
    parseExp
  ).map (λ ((name, typeA), _, typeB) => Exp.pi name typeA typeB)

  ||

  (
    parseExp *
    parseExact "->" *
    parseExp
  ).map (λ (typeA, _, typeB) => Exp.pi "_" typeA typeB)



def parseLam (parseExp: Parser Exp): Parser Exp :=
  (
    parseString *
    parseExact "=>" *
    parseExp
  ).map (λ (name, _, body) => Exp.lam name body)

def parseBnd (parseExp: Parser Exp): Parser Exp :=
  (
    parseExact "let" *
    parseString * -- name
    parseExact ":" *
    parseExp * -- type
    parseExact ":=" *
    parseExp * -- value
    parseExact "in" *
    parseExp -- body
  ).map (λ (_, name, _, type, _, value, _, body) => Exp.bnd name value type body)

def parseInh (parseExp: Parser Exp): Parser Exp :=
  (
    parseExact "inh" *
    parseString * -- name
    parseExact ":" *
    parseExp * -- type
    parseExact "in" *
    parseExp -- body
  ).map (λ (_, name, _, type, _, body) => Exp.inh name type body)

def parseApp (parseExp: Parser Exp): Parser Exp :=
  sorry

def parseExp: Parser Exp :=
  sorry







end EL2.Parser.Internal


namespace EL2.Parser

export EL2.Parser.Internal (newTokenizer)

end EL2.Parser
