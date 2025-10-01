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

#eval sortSplitTokens ["hello", "world", "hello world", "world hello"]

private def splitPart (sortedSplitTokens : List String) (part : String) : List String :=
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
