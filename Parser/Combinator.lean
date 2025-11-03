namespace Parser.Combinator

def Parser χ α := List χ → Option (α × List χ)

def Parser.mapOption (p: Parser χ α) (f: α → Option β): Parser χ β := λ xs => do
  let (a, xs) ← p xs
  let b ← f a
  (b, xs)

def Parser.map (p: Parser χ α) (f: α → β): Parser χ β := p.mapOption (λ a => some (f a))

def Parser.concat (p1: Parser χ α) (p2: Parser χ β): Parser χ (α × β) := λ xs => do
  let (a, xs) ← p1 xs
  let (b, xs) ← p2 xs
  ((a, b), xs)

infixr: 60 " ++ " => Parser.concat

def Parser.either (p1: Parser χ α) (p2: Parser χ α): Parser χ α := λ xs =>
  match p1 xs with
    | some (a, xs) => some (a, xs)
    | none => p2 xs

infixr: 50 " || " => Parser.either -- lower precedence than concat

partial def Parser.list (p: Parser χ α): Parser χ (List α) := λ xs =>
  let rec loop (as: Array α) (xs: List χ): Option (List α × List χ) :=
    match p xs with
      | none => some (as.toList, xs)
      | some (a, rest) => loop (as.push a) rest
  loop #[] xs

def parseEmpty: Parser χ Unit := λ xs => some ((), xs)
def parseFail: Parser χ Unit := λ _ => none

def parseExact [BEq χ] (y: χ): Parser χ χ := λ xs =>
  match xs with
    | [] => none
    | x :: xs =>
      if y == x then
        some (x, xs)
      else
        none

def parseExactList [BEq χ] (ys: List χ): Parser χ (List χ) := λ xs => do
  let rest ← ys.isPrefixOf? xs
  pure (ys, rest)

#eval parseExactList "hehe".toList "hehea123".toList

-- STRING

def parseNewLine : Parser Char Unit := parseExact '\n'

def parseWhiteSpaceWeak : Parser Char String := λ xs =>
  -- parse any whitespace
  -- empty whitespace is ok
  let rec loop (ys: Array Char) (xs: List Char): Option (String × List Char) :=
    match xs with
      | [] => (String.mk ys.toList, xs)
      | x :: rest =>
        if x.isWhitespace then
          loop (ys.push x) rest
        else
          (String.mk ys.toList, xs)

  loop #[] xs

def parseWhiteSpace : Parser Char String :=
  -- parse some whitespace
  -- empty whitespace is not ok
  parseWhiteSpaceWeak.mapOption (λ s => if s.length = 0 then none else s)

def parseNameWeak: Parser Char String := λ xs =>
  -- parse a non-whitespace string
  -- empty name is ok
  let rec loop (ys: Array Char) (xs: List Char): Option (String × List Char) :=
    match xs with
      | [] => (String.mk ys.toList, xs)
      | x :: rest =>
        if ¬ x.isWhitespace then
          loop (ys.push x) rest
        else
          (String.mk ys.toList, xs)
  loop #[] xs

def parseName: Parser Char String :=
  -- parse a non-whitespace string
  -- empty name is not ok
  parseNameWeak.mapOption (λ s => if s.length = 0 then none else s)

def parseExactString (ys: String): Parser Char Unit := parseExactList ys.toList

#eval parseName "abc123  ".toList
#eval parseWhiteSpace "abc123  ".toList
#eval parseName "   abc123".toList
#eval parseWhiteSpace "   abc123".toList


def parseDigit


end Parser.Combinator
