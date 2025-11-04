namespace Parser.Combinator

def Parser χ α := List χ → Option (α × List χ)

def Parser.filterMap (p: Parser χ α) (f: α → Option β): Parser χ β := λ xs => do
  let (a, xs) ← p xs
  let b ← f a
  (b, xs)

def Parser.map (p: Parser χ α) (f: α → β): Parser χ β := p.filterMap (λ a => some (f a))

def Parser.concat (p1: Parser χ α) (p2: Parser χ β): Parser χ (α × β) := λ xs => do
  let (a, xs) ← p1 xs
  let (b, xs) ← p2 xs
  ((a, b), xs)

infixr: 60 " ++ " => Parser.concat

def Parser.sum (p1: Parser χ α) (p2: Parser χ β): Parser χ (α ⊕ β) := λ xs =>
  match p1 xs with
    | some (a, xs) => some (Sum.inl a, xs)
    | none => match p2 xs with
      | some (b, xs) => some (Sum.inr b, xs)
      | none => none

def Parser.either (p1: Parser χ α) (p2: Parser χ α): Parser χ α :=
  (p1.sum p2).map $ Sum.elim id id

infixr: 50 " || " => Parser.either -- lower precedence than concat

infixr: 40 " ||| " => Parser.sum -- lower precedence than either


partial def Parser.list (p: Parser χ α): Parser χ (List α) := λ xs =>
  let rec loop (as: Array α) (xs: List χ): Option (List α × List χ) :=
    match p xs with
      | none => some (as.toList, xs)
      | some (a, rest) => loop (as.push a) rest
  loop #[] xs

def empty: Parser χ Unit := λ xs => some ((), xs)
def fail: Parser χ Unit := λ _ => none

def pred (p: χ → Bool): Parser χ χ := λ xs =>
  match xs with
    | [] => none
    | x :: xs =>
      if p x then
        some (x, xs)
      else
        none

def exact [BEq χ] (y: χ): Parser χ χ := pred (· == y)

def exactList [BEq χ] (ys: List χ): Parser χ (List χ) := λ xs => do
  let rest ← ys.isPrefixOf? xs
  pure (ys, rest)

def nonEmpty (p: Parser χ (List α)): Parser χ (List α) := λ xs => do
  let (as, xs) ← p xs
  if as.length = 0 then
    none
  else
    pure (as, xs)

#eval exactList "hehe".toList "hehea123".toList

namespace String

def toString (p: Parser Char (List Char)): Parser Char String :=
  p.map (String.mk ·)

def whitespaceWeak : Parser Char String :=
  -- parse any whitespace
  -- empty whitespace is ok
  toString (pred (·.isWhitespace)).list

def whiteSpaceWithoutNewLineWeak : Parser Char String :=
  -- parse any whitespace but not new line
  -- empty whitespace is ok
  toString (pred (λ c => c.isWhitespace ∧ (¬ c = '\n'))).list


def whitespace : Parser Char String :=
  -- parse some whitespace
  -- empty whitespace is not ok
  whitespaceWeak.filterMap (λ s => if s.length = 0 then none else s)

def nameWeak: Parser Char String :=
  -- parse a non-whitespace string
  -- empty name is ok
  toString (pred (¬ ·.isWhitespace)).list

def name: Parser Char String :=
  -- parse a non-whitespace string
  -- empty name is not ok
  nameWeak.filterMap (λ s => if s.length = 0 then none else s)



def parseExactString (ys: String): Parser Char String :=
  (parseExactList ys.toList).map (λ ys => String.mk ys)

#eval parseName "abc123  ".toList
#eval parseWhiteSpace "abc123  ".toList
#eval parseName "   abc123".toList
#eval parseWhiteSpace "   abc123".toList


end String

end Parser.Combinator
