namespace Parser.Combinator

def Parser χ α := List χ → Option (α × List χ)

def Parser.fail : Parser χ α := λ _ => none
def Parser.pure (a: α): Parser χ α := λ xs => some (a, xs)

def Parser.bind (p: Parser χ α) (f: α → Parser χ β): Parser χ β := λ xs =>
  match p xs with
    | none => none
    | some (a, xs) => f a xs

instance: Monad (Parser χ) where
  pure := Parser.pure
  bind := Parser.bind


def Parser.filterMap (p: Parser χ α) (f: α → Option β): Parser χ β := λ xs => do
  let (a, xs) ← p xs
  let b ← f a
  (b, xs)

def Parser.map (p: Parser χ α) (f: α → β): Parser χ β :=
  p.filterMap (λ a => some (f a))


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

def Parser.fromList (ps: List (Parser χ α)): Parser χ (List α) := λ xs =>
  let rec loop (ys: Array α) (ps: List (Parser χ α)) (xs: List χ): Option (List α × List χ) :=
    match ps with
      | [] => some (ys.toList, xs)
      | p :: ps =>
        match p xs with
          | none => none
          | some (y, xs) =>
            loop (ys.push y) ps xs

  loop #[] ps xs

def pred (p: χ → Bool): Parser χ χ := λ xs =>
  match xs with
    | [] => none
    | x :: xs =>
      if p x then
        some (x, xs)
      else
        none

def exact [BEq χ] (y: χ): Parser χ χ := pred (· == y)

def exactList [BEq χ] (ys: List χ): Parser χ (List χ) :=
  Parser.fromList (ys.map exact)


def nonEmpty (p: Parser χ (List α)): Parser χ (List α) := λ xs => do
  let (as, xs) ← p xs
  if as.length = 0 then
    none
  else
    pure (as, xs)

#eval exactList "hehe".toList "hehea123".toList

namespace String

def toStringParser (p: Parser Char (List Char)): Parser Char String :=
  p.map (String.mk ·)

def whitespaceWeak : Parser Char String :=
  -- parse any whitespace
  -- empty whitespace is ok
  toStringParser (pred (·.isWhitespace)).list

def whiteSpaceWithoutNewLineWeak : Parser Char String :=
  -- parse any whitespace but not new line
  -- empty whitespace is ok
  toStringParser (pred (λ c => c.isWhitespace ∧ (¬ c = '\n'))).list


def whitespace : Parser Char String :=
  -- parse some whitespace
  -- empty whitespace is not ok
  whitespaceWeak.filterMap (λ s => if s.length = 0 then none else s)

def exact (ys: String): Parser Char String :=
  toStringParser (exactList ys.toList)

#eval whitespace "abc123  ".toList
#eval whitespace "   abc123".toList
#eval exact "let" "let123".toList


end String

end Parser.Combinator
