namespace Form

universe u

inductive Form where
  | name: String → Form
  | list: List Form → Form
  deriving Repr

#eval (Form.name "hello" : Form)
#eval (Form.list [Form.name "hello", Form.name "world"] : Form)


private def sortSplitTokens (splitTokens : List String) : List String :=
  -- sort tokens so that if s2 is a prefix of s1, s1 should come first
  let lessEqual (s1: String) (s2: String): Bool :=
    if (s2.length < s1.length) && (s2.isPrefixOf s1) then true else
    if (s1.length < s2.length) && (s1.isPrefixOf s2) then false else
    s1 < s2

  splitTokens.mergeSort lessEqual

#eval sortSplitTokens ["=", "==", ":="]


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

-- Test stringIndexOf
#eval stringIndexOf "hello world" "world"  -- should be some 6
#eval stringIndexOf "hello world" "hello"  -- should be some 0
#eval stringIndexOf "hello world" "xyz"    -- should be none
#eval stringIndexOf "hello world" ""       -- should be some 0


private partial def splitPart (sortedSplitTokens : List String) (part : String) : List String :=
  match sortedSplitTokens with
    | [] => [part]
    | s :: ss =>
      match stringIndexOf part s with
        | some i =>
          [part.take i] ++ [s] ++
            splitPart ss (part.drop (i + s.length))
        | none => splitPart ss part


#eval splitPart (sortSplitTokens ["=", "==", ":="]) "x:=3==2=1"


def tokenize (sortedSplitTokens : List String) (s : String) : List String :=
  let parts := s.split (fun c => c.isWhitespace)
  parts.flatMap (splitPart sortedSplitTokens)


end Form
