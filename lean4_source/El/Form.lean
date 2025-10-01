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
  sorry

private def tokenize (sortedSplitTokens : List String) (s : String) : String :=



  sorry



end Form
