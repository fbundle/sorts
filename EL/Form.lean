import EL.Util

namespace Form

inductive Form where
  | name: String → Form
  | list: List Form → Form
  deriving Repr

def getName (form: Form): Option String :=
  match form with
    | .name name => some name
    | _ => none

def getList (form: Form): Option (List Form) :=
  match form with
    | .list list => some list
    | _ => none

def toString (form: Form) : String :=
  match form with
    | Form.name s => s
    | Form.list fs => "(" ++ String.join ((fs.map toString).intersperse " ") ++ ")"

instance : ToString Form := ⟨toString⟩


private def _sortSplitTokens (splitTokens : List String) : List String :=
  -- sort tokens so that if s2 is a prefix of s1, s1 should come first
  let lessEqual (s1: String) (s2: String): Bool :=
    if (s2.length < s1.length) && (s2.isPrefixOf s1) then true else
    if (s1.length < s2.length) && (s1.isPrefixOf s2) then false else
    s1 < s2

  splitTokens.mergeSort lessEqual


private partial def _stringIndexOf (s: String) (substring: String): Option Nat :=
  -- return the starting position of substring in s if exists
  if substring.isEmpty then
    some 0
  else
    let rec helper (s: String) (substring: String) (startIdx: Nat) : Option Nat :=
      if startIdx + substring.length > s.length then
        none
      else if s.extract ⟨startIdx⟩ ⟨startIdx + substring.length⟩ = substring then
        some startIdx
      else
        helper s substring (startIdx + 1)
    helper s substring 0


private partial def _splitPart (sortedSplitTokens : List String) (part : String) : List String :=
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

private def _tokenize (sortedSplitTokens : List String) (s : String) : List String :=
  let parts := s.split (λ c => c.isWhitespace)
  let output := parts.flatMap (_splitPart sortedSplitTokens)
  let output := output.filter (λ s => s.length > 0)
  output



def parser := List String → Option (List String × Form)

private partial def _parseUntil (p: parser)  (closeBlockToken: String) (acc: List Form) (tokens : List String) : Option (List String × List Form) :=
  match tokens with
    | [] => none
    | t :: ts =>
      if t = closeBlockToken then
        some ⟨ts, acc⟩
      else
        match p tokens with
          | some ⟨remainingTokens, form⟩ =>
            _parseUntil p closeBlockToken (acc ++ [form]) remainingTokens
          | none => none

private partial def _parse (openBlockToken: String) (closeBlockToken: String) (tokens : List String) : Option (List String × Form) :=
  match tokens with
    | [] => none
    | t :: ts =>
      if t = openBlockToken then
        match _parseUntil (_parse openBlockToken closeBlockToken) closeBlockToken [] ts with
          | some ⟨ts, forms⟩ => some ⟨ts, Form.list forms⟩
          | none => none
      else
        some ⟨ts, Form.name t⟩

structure Parser where
  openBlockToken: String
  closeBlockToken: String
  splitTokens: List String

def Parser.init (p: Parser) : Parser :=
  {p with splitTokens := _sortSplitTokens p.splitTokens}

def Parser.tokenize (p: Parser) (s: String): List String :=
  _tokenize p.splitTokens s

def Parser.parse (p: Parser) (tokens: List String): Option (List String × Form) :=
  _parse p.openBlockToken p.closeBlockToken tokens


-- default parser

def defaultParser := ({
  openBlockToken := "(",
  closeBlockToken := ")",
  splitTokens := ["(", ")", "+", "-", "*", "/", "=", "==", ":=", "=>", "->"]
}: Parser).init

private def _example := "x:=(3==2)=1 123"

#eval Util.parseAll defaultParser.parse (defaultParser.tokenize _example)

end Form
