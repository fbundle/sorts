namespace EL2.ParserCombinator

def Parser χ α  := χ → Option (χ × α)

def fail : Parser χ α := λ _ => none

def pure (a: α): Parser χ α := λ xs =>
  some (xs, a)

def Parser.bind (p: Parser χ α) (f: α → Parser χ β): Parser χ β := λ xs => do
  let (xs, a) ← p xs
  f a xs

instance: Monad (Parser α) where
  pure := pure
  bind := Parser.bind

def Parser.filter (p: Parser χ α) (f: α → Bool): Parser χ α :=
  p.bind (λ a =>
    if f a then
      pure a
    else
      fail
  )

def Parser.map (p: Parser χ α) (f: α → β): Parser χ β :=
  p.bind (λ a => pure (f a))

def Parser.concat (p1: Parser χ α) (p2: Parser χ β): Parser χ (α × β) := λ xs => do
  let (xs, a) ← p1 xs
  let (xs, b) ← p2 xs
  some (xs, (a, b))

infixr: 60 " ++ " => Parser.concat

def Parser.either (p1: Parser χ α) (p2: Parser χ α): Parser χ α := λ xs =>
  match p1 xs with
    | some (xs, a) => some (xs, a)
    | none => p2 xs

infixr: 50 " || " => Parser.either -- lower precedence than concat

partial def Parser.many (p: Parser χ α): Parser χ (List α) := λ xs =>
  let rec loop (as: Array α) (xs: χ): Option (χ × List α) :=
    match p xs with
      | none => some (xs, as.toList)
      | some (rest, a) => loop (as.push a) rest
  loop #[] xs

def Parser.many1 (p: Parser χ α): Parser χ (List α) := λ xs => do
  let (xs, as) ← p.many xs
  if as.length = 0 then
    none
  else
    some (xs, as)

def Parser.transpose (ps: List (Parser χ α)): Parser χ (List α) := λ xs =>
  let rec loop (ys: Array α) (ps: List (Parser χ α)) (xs: χ): Option (χ × List α) :=
    match ps with
      | [] => some (xs, ys.toList)
      | p :: ps =>
        match p xs with
          | none => none
          | some (xs, y) =>
            loop (ys.push y) ps xs
  loop #[] ps xs


def pred (p: χ → Bool): Parser (List χ) χ := λ xs =>
  match xs with
    | [] => none
    | x :: xs =>
      if p x then
        some (xs, x)
      else
        none

def exact [BEq χ] (y: χ): Parser (List χ) χ := pred (· == y)

def exactList [BEq χ] (ys: List χ): Parser (List χ) (List χ) :=
  Parser.transpose (ys.map exact)

namespace String

def toStringParser (p: Parser (List Char) (List Char)): Parser (List Char) String :=
  p.map (String.mk ·)

def whitespaceWeak : Parser (List Char) String :=
  -- parse any whitespace
  -- empty whitespace is ok
  toStringParser (pred (·.isWhitespace)).many

def whiteSpaceWithoutNewLineWeak : Parser (List Char) String :=
  -- parse any whitespace but not new line
  -- empty whitespace is ok
  toStringParser (pred (λ c => c.isWhitespace ∧ (¬ c = '\n'))).many


def whitespace : Parser (List Char) String :=
  -- parse some whitespace
  -- empty whitespace is not ok
  toStringParser (pred (·.isWhitespace)).many1

def exact (ys: String): Parser (List Char) String :=
  toStringParser (exactList ys.toList)

end String

end EL2.ParserCombinator
