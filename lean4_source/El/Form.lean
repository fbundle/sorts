namespace Form

universe u

inductive Form where
  | name: String → Form
  | list: List Form → Form
  deriving Repr

#eval (Form.name "hello" : Form)
#eval (Form.list [Form.name "hello", Form.name "world"] : Form)


def sortSplitTokens (splitTokens : List String) : List String :=
  -- sort tokens so that if s2 is a prefix of s1, s1 should come first
  let lessEqual (s1: String) (s2: String): Bool :=
    if (s2.length < s1.length) && (s2.isPrefixOf s1) then true else
    if (s1.length < s2.length) && (s1.isPrefixOf s2) then false else
    s1 < s2

  splitTokens.mergeSort lessEqual

#eval sortSplitTokens ["=", "==", ":="]


def stringIndexOf (s: String) (substring: String): Option Int :=
  -- return the starting position of substring in s if exists
  if substring.isEmpty then
    some 0
  else
    let rec helper (s: String) (substring: String) (startIdx: Nat) : Option Int :=
      if startIdx + substring.length > s.length then
        none
      else if s.extract ⟨startIdx⟩ ⟨startIdx + substring.length⟩ = substring then
        some startIdx
      else
        helper s substring (startIdx + 1)
      decreasing_by all_goals sorry
    helper s substring 0

-- Test stringIndexOf
#eval! stringIndexOf "hello world" "world"  -- should be some 6
#eval! stringIndexOf "hello world" "hello"  -- should be some 0
#eval! stringIndexOf "hello world" "xyz"    -- should be none
#eval! stringIndexOf "hello world" ""       -- should be some 0


private def splitPart (sortedSplitTokens : List String) (part : String) : List String :=
  let rec loop (acc: List String) (sortedSplitTokens: List String) (part: String) : List String :=
    match sortedSplitTokens with
      | [] => acc
      | s :: ss =>
        if s.isPrefixOf part then
          loop (acc ++ [s]) sortedSplitTokens (part.drop s.length)
        else
          loop acc ss part
    decreasing_by all_goals sorry


    match sortedSplitTokens with
    | [] => [part]
    | s :: ss =>
      if s.isPrefixOf part then
        s :: (splitPart sortedSplitTokens (part.drop s.length))
      else
        splitPart ss part
    decreasing_by all_goals sorry

#eval! splitPart (sortSplitTokens ["=", "==", ":="]) "x:=3==2=1"


def tokenize (sortedSplitTokens : List String) (s : String) : List String :=
  let parts := s.split (fun c => c.isWhitespace)
  parts.flatMap (splitPart sortedSplitTokens)


end Form
