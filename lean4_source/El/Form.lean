namespace Form

universe u

inductive Form where
  | name: String → Form
  | list: List Form → Form

def toString (form: Form) : String :=
  match form with
    | Form.name s => s
    | Form.list fs => "[" ++ String.join ((fs.map toString).intersperse " ") ++ "]"

instance : ToString Form := ⟨toString⟩


private def sortSplitTokens (splitTokens : List String) : List String :=
  -- sort tokens so that if s2 is a prefix of s1, s1 should come first
  let lessEqual (s1: String) (s2: String): Bool :=
    if (s2.length < s1.length) && (s2.isPrefixOf s1) then true else
    if (s1.length < s2.length) && (s1.isPrefixOf s2) then false else
    s1 < s2

  splitTokens.mergeSort lessEqual


private partial def stringIndexOf (s: String) (substring: String): Option Nat :=
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


private partial def splitPart (sortedSplitTokens : List String) (part : String) : List String :=
  match sortedSplitTokens with
    | [] => [part]
    | s :: ss =>
      match stringIndexOf part s with
        | some n =>
          let before := part.take n
          let after := part.drop (n + s.length)
          let beforeParts := if before.isEmpty then [] else splitPart sortedSplitTokens before
          let afterParts := if after.isEmpty then [] else splitPart sortedSplitTokens after
          beforeParts ++ [s] ++ afterParts
        | none => splitPart ss part

private def tokenize (sortedSplitTokens : List String) (s : String) : List String :=
  let parts := s.split (fun c => c.isWhitespace)
  parts.flatMap (splitPart sortedSplitTokens)


def parser := List String → Option (List String × Form)


private partial def parseUntil (p: parser)  (closeBlockToken: String) (acc: List Form) (tokens : List String) : Option (List String × List Form) :=
  match tokens with
    | [] => none
    | t :: ts =>
      if t = closeBlockToken then
        some ⟨ts, acc⟩
      else
        match p tokens with
          | some ⟨remainingTokens, form⟩ =>
            parseUntil p closeBlockToken (acc ++ [form]) remainingTokens
          | none => none

partial def parse (openBlockToken: String) (closeBlockToken: String) (tokens : List String) : Option (List String × Form) :=
  match tokens with
    | [] => none
    | t :: ts =>
      if t = openBlockToken then
        match parseUntil (parse openBlockToken closeBlockToken) closeBlockToken [] ts with
          | some ⟨ts, forms⟩ => some ⟨ts, Form.list forms⟩
          | none => none
      else
        some ⟨ts, Form.name t⟩

partial def parseAll (openBlockToken: String) (closeBlockToken: String) (tokens : List String) : Option (List Form) :=
  let rec loop (acc: List Form) (tokens : List String) : Option (List Form) :=
    match tokens with
      | [] => some acc
      | _ :: _ =>
        match parse openBlockToken closeBlockToken tokens with
          | some ⟨ts, form⟩ => loop (acc ++ [form]) ts
          | none => none
  loop [] tokens

def openBlockToken := "("
def closeBlockToken := ")"
def sortedSplitTokens := sortSplitTokens ["(", ")", "+", "-", "*", "/", "=", "==", ":="]


#eval tokenize sortedSplitTokens "x:=(3==2)=1"

#eval parse "(" ")" (tokenize sortedSplitTokens "x:=(3==2)=1")

#eval parse "(" ")" [":=", "(", "3", "==", "2", ")", "=", "1"]

#eval parse "(" ")" ["(", "3", "==", "2", ")", "=", "1"]

#eval parse "(" ")" ["(", "3", ")", "=", "1"]

#eval parseAll "(" ")" (tokenize sortedSplitTokens "x:=(3==2)=1")

end Form
