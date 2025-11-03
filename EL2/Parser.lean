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

infixr: 60 " ++ " => Parser.concat

def Parser.either (p1: Parser α) (p2: Parser α): Parser α := λ tokens => do
  match p1 tokens with
    | some (rest, a) => some (rest, a)
    | none => p2 tokens

infixr: 50 " || " => Parser.either -- lower precedence than concat

partial def Parser.many (p: Parser α): Parser (List α) :=
  let rec loop (acc: Array α) (tokens: List String): Option (List String × List α) :=
    match p tokens with
      | none => some (tokens, acc.toList)
      | some (rest, e) => loop (acc.push e) rest
  loop #[]

open EL2.Core

def parseFail: Parser α := λ _ => none

def parseString: Parser String := λ tokens =>
  match tokens with
    | [] => none
    | head :: tokens =>
      (tokens, head)

def parseExact (pattern: String): Parser String := parseString.mapPartial $ (λ string =>
  if pattern = string then
    some pattern
  else
    none
)

def parseExactMany (patterns: List String): Parser String :=
  (patterns.map parseExact).foldl
  Parser.either parseFail

-- parse Exp

def parseTyp: Parser Exp := parseString.mapPartial (λ head => do
  if ¬ "Type".isPrefixOf head then none else
  let levelStr := head.stripPrefix "Type"
  let level ← levelStr.toNat?
  dbg_trace s!"parsed type {level}"
  pure (Exp.typ level)
)

def parseVar (specialTokens: List String): Parser Exp :=  λ tokens => do
  match parseExactMany specialTokens tokens with
    | some _ => none
    | none =>
      let (tokens, name) ← parseString tokens
      some (tokens, Exp.var name)

-- def parseApp:

def parseAnn (parseExp: Parser Exp): Parser (String × Exp) :=
  (
    parseExact "(" ++
    parseString ++
    (parseExact ":") ++
    parseExp ++
    parseExact ")"
  ).map (λ (_, name, _, type, _) => (name, type))

def parsePi (parseExp: Parser Exp): Parser Exp :=
  -- named Pi or unnamed Pi
  (
    (parseExact "Π" || parseExact "∀" || parseExact "forall") ++
    parseAnn parseExp ++
    parseExact "->" ++
    parseExp
  ).map (λ (_, (name, typeA), _, typeB) => Exp.pi name typeA typeB)

  ||

  (
    (parseExact "Π" || parseExact "∀" || parseExact "forall") ++
    parseExp ++
    parseExact "->" ++
    parseExp
  ).map (λ (_, typeA, _, typeB) => Exp.pi "_" typeA typeB)

def parseLam (parseExp: Parser Exp): Parser Exp :=
  (
    (parseExact "λ" || parseExact "lam") ++
    parseString ++
    parseExact "=>" ++
    parseExp
  ).map (λ (_, name, _, body) => Exp.lam name body)

def parseBnd (parseExp: Parser Exp): Parser Exp :=
  (
    parseExact "let" ++
    parseString ++ -- name
    parseExact ":=" ++
    parseExp ++ -- value
    parseExact "in" ++
    parseExp -- body
  ).map (λ (_, name, _, value, _, body) =>
    Exp.app (Exp.lam name body) value
  )

  ||

  (
    parseExact "let" ++
    parseString ++ -- name
    parseExact ":" ++
    parseExp ++ -- type
    parseExact ":=" ++
    parseExp ++ -- value
    parseExact "in" ++
    parseExp -- body
  ).map (λ (_, name, _, type, _, value, _, body) => Exp.bnd name value type body)

def parseInh (parseExp: Parser Exp): Parser Exp :=
  (
    parseExact "inh" ++
    parseString ++ -- name
    parseExact ":"++
    parseExp ++ -- type
    parseExact "in" ++
    parseExp -- body
  ).map (λ (_, name, _, type, _, body) => Exp.inh name type body)

def specialTokens: List String := [
  ":", "->", "=>", "let", ":=", "in", "inh", "(", ")", "λ", "lam", "Π", "∀", "forall",
]

def parseApp (parseExp: Parser Exp): Parser Exp :=
  (
    parseExact "(" ++
    parseExp.many ++
    parseExact ")"
  ).mapPartial (λ (_, es, _) =>
    match es with
      | [] => none
      | cmd :: args =>
        some $ args.foldl (λ cmd arg =>
          Exp.app cmd arg
        ) cmd
  )

partial def parseExp: Parser Exp :=
  parseTyp ||
  parseVar specialTokens ||
  parsePi parseExp ||
  parseLam parseExp ||
  parseBnd parseExp ||
  parseInh parseExp ||
  parseApp parseExp

#eval parseExp ["inh", "name", ":", "type", "in", "hehe"]
#eval parseExp ["forall", "type1", "->", "type2"]
#eval parseExp ["forall", "(", "name1", ":", "type1", ")", "->", "type2"]
#eval parseExp ["lam", "name", "=>", "body"]
#eval parseExp ["let", "x", ":=", "3", "in", "x+y"]
#eval parseExp ["let", "x", ":", "type", ":=", "3", "in", "x+y"]
#eval parseExp ["(", "cmd", "arg1", "arg2", ")"]
#eval parseExp ["(", "cmd", ")"]


end EL2.Parser.Internal




namespace EL2.Parser

def tokenize := (Internal.newTokenizer ("(" :: ")" :: Internal.specialTokens))

def parse := Internal.parseExp

end EL2.Parser
