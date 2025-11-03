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
  let (rest, a) ← p tokens
  let b ← f a
  pure (rest, b)

def Parser.map (p: Parser α) (f: α → β): Parser β := p.mapPartial (λ a => some (f a))

open EL2.Core

def parseHead: Parser String := λ tokens =>
  match tokens with
    | [] => none
    | head :: rest =>
      (rest, head)

def parseSingle (convert?: String → Option α): Parser α := parseHead.mapPartial convert?

def predToOption (f: α → Bool): α → Option α := λ a =>
  if f a then some a else none

def parseExact (pattern: String): Parser String := parseSingle $ predToOption (· = pattern)

def parseTyp: Parser Exp := parseSingle (λ head => do
  if ¬ "Type".isPrefixOf head then none else
  let levelStr := head.stripPrefix "Type"
  let level ← levelStr.toNat?
  pure (Exp.typ level)
)

def parseVar: Parser Exp :=







end EL2.Parser.Internal


namespace EL2.Parser

export EL2.Parser.Internal (newTokenizer)

end EL2.Parser
