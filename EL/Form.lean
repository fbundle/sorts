import Std.Data


namespace Form

inductive Form where
  | name: String → Form
  | list: List Form → Form
  deriving Repr

def toString (form: Form) : String :=
  match form with
    | Form.name s => s
    | Form.list fs => "(" ++ String.join ((fs.map toString).intersperse " ") ++ ")"

instance : ToString Form where
 toString := toString


def _sortSplitTokens (splitTokens : List String) : List String :=
  -- sort tokens so that if s2 is a prefix of s1, s1 should come first
  let lessEqual (s1: String) (s2: String): Bool :=
    if (s2.length < s1.length) && (s2.isPrefixOf s1) then true else
    if (s1.length < s2.length) && (s1.isPrefixOf s2) then false else
    s1 < s2

  splitTokens.mergeSort lessEqual


partial def _stringIndexOf (s: String) (substring: String): Option Nat :=
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


partial def _splitPart (sortedSplitTokens : List String) (part : String) : List String :=
  match sortedSplitTokens with
    | [] => [part]
    | s :: ss =>
      match _stringIndexOf part s with
        | some n =>
          let before := part.take n
          let after := part.drop (n + s.length)
          let beforeParts := if before.isEmpty then [] else _splitPart sortedSplitTokens before
          let afterParts := if after.isEmpty then [] else _splitPart sortedSplitTokens after
          beforeParts ++ [s] ++ afterParts
        | none => _splitPart ss part


partial def _parseUntilClose (parse: List String → Option (List String × Form)) (closeBlockToken: String) (forms: Array Form) (tokens : List String) : Option (List String × Array Form) :=
  match tokens with
    | [] => none
    | t :: ts =>
      if t = closeBlockToken then
        some (ts, forms)
      else
        match parse tokens with
          | some (tokens, form) =>
            _parseUntilClose parse closeBlockToken (forms.push form) tokens
          | none => none




structure BlockParser where
  openBlockToken: String
  closeBlockToken: String
  postProcess: List Form → Option (List Form)

structure Parser where
  blockParsers: Std.HashMap String BlockParser
  splitTokens: List String

def newParser (splitTokens: List String): Parser :=
  {
    blockParsers := Std.HashMap.emptyWithCapacity,
    splitTokens := _sortSplitTokens splitTokens
  }

def Parser.addBlockParser (p: Parser) (bp: BlockParser): Parser :=
  {p with
    blockParsers := p.blockParsers.insert bp.openBlockToken bp,
    splitTokens := _sortSplitTokens (p.splitTokens ++ [bp.openBlockToken, bp.closeBlockToken]),
  }

def Parser.tokenize (p: Parser) (s : String) : List String :=
  let parts := s.split (·.isWhitespace)
  let output := parts.flatMap (_splitPart p.splitTokens)
  let output := output.filter (·.length > 0)
  output

partial def Parser.parse (p: Parser) (tokens: List String): Option (List String × Form) :=
 match tokens with
    | [] => none
    | t :: ts =>
      match p.blockParsers.get? t with
        | none => some (ts, Form.name t)
        | some bp => do
          let (ts, forms) ← _parseUntilClose p.parse bp.closeBlockToken #[] ts
          let forms ← bp.postProcess forms.toList
          pure (ts, Form.list forms)

-- default parser

partial def infixProcess (forms: List Form): Option (List Form) := do
  if forms.length ≤ 2 then
    forms
  else
    let op ← forms[forms.length-2]?
    let last ← forms[forms.length-1]?
    let init := (forms.extract 0 (forms.length-2))
    let init ← infixProcess init
    pure [op, Form.list init, last]


def defaultParser := ((
    newParser [
      "+", "-", "*", "/", "%", "=", "==",
      ",", ";", ":", ":=", "=>", "->",
      ]
  ).addBlockParser {
    openBlockToken := "(",
    closeBlockToken := ")",
    postProcess := (some ·),
  }).addBlockParser {
    openBlockToken := "{",
    closeBlockToken := "}",
    postProcess := infixProcess,
  }

end Form
