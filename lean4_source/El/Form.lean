namespace Form

inductive Form where
  | name: String → Form
  | list: List Form → Form

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
  let parts := s.split (fun c => c.isWhitespace)
  parts.flatMap (_splitPart sortedSplitTokens)


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

private partial def iterate (iter: α → Option (α × β)) (terminate: α → Bool) (init: α): Option (List β) :=
  let rec loop (acc: List β) (state: α): Option (List β) :=
    match terminate state with
      | true => acc
      | false =>
        match iter state with
          | none => none
          | some ⟨next_state, value⟩ => loop (acc ++ [value]) next_state
  loop [] init

def Parser.parseAll (p: Parser) (tokens: List String) : Option (List Form) :=
  iterate p.parse (λ tokens => tokens.length = 0) tokens

def defaultParser := ({
  openBlockToken := "(",
  closeBlockToken := ")",
  splitTokens := ["(", ")", "+", "-", "*", "/", "=", "==", ":="]
}: Parser).init

#eval (defaultParser.parseAll (defaultParser.tokenize "x:=(3==2)=1 123")).get!





end Form
