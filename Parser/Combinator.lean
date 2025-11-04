namespace Parser.Combinator

structure Error χ where
  xs: List χ
  deriving Repr

def err (xs: List χ): Except (Error χ) α :=
  Except.error {xs := xs}

def Parser χ α := List χ → Except (Error χ) (α × List χ)

def Parser.fail : Parser χ α := λ xs =>
  err xs

def Parser.pure (a: α): Parser χ α := λ xs =>
  Except.ok (a, xs)

def Parser.bind (p: Parser χ α) (f: α → Parser χ β): Parser χ β := λ xs => do
  let (a, xs) ← p xs
  f a xs

instance: Monad (Parser χ) where
  pure := Parser.pure
  bind := Parser.bind


def Parser.filterMap (p: Parser χ α) (f: α → Option β): Parser χ β := λ xs => do
  let (a, xs) ← p xs
  match f a with
    | none => err xs
    | some b => Except.ok (b, xs)

def Parser.map (p: Parser χ α) (f: α → β): Parser χ β :=
  p.filterMap (λ a => some (f a))


def Parser.concat (p1: Parser χ α) (p2: Parser χ β): Parser χ (α × β) := λ xs => do
  let (a, xs) ← p1 xs
  let (b, xs) ← p2 xs
  Except.ok ((a, b), xs)

infixr: 60 " ++ " => Parser.concat

def Parser.sum (p1: Parser χ α) (p2: Parser χ β): Parser χ (α ⊕ β) := λ xs =>
  match p1 xs with
    | Except.ok (a, xs) => Except.ok (Sum.inl a, xs)
    | Except.error _ => match p2 xs with
      | Except.ok (b, xs) => Except.ok (Sum.inr b, xs)
      | Except.error err => Except.error err

def Parser.either (p1: Parser χ α) (p2: Parser χ α): Parser χ α :=
  (p1.sum p2).map $ Sum.elim id id

infixr: 50 " || " => Parser.either -- lower precedence than concat

infixr: 40 " ||| " => Parser.sum -- lower precedence than either


partial def Parser.many (p: Parser χ α): Parser χ (List α) := λ xs =>
  let rec loop (as: Array α) (xs: List χ): Except (Error χ) (List α × List χ) :=
    match p xs with
      | Except.error _ => Except.ok (as.toList, xs)
      | Except.ok (a, rest) => loop (as.push a) rest
  loop #[] xs

def nonEmpty (p: Parser χ (List α)): Parser χ (List α) := λ xs => do
  let (as, xs) ← p xs
  if as.length = 0 then
    err xs
  else
    pure (as, xs)

def Parser.many1 (p: Parser χ α): Parser χ (List α) := nonEmpty p.many

def Parser.transpose (ps: List (Parser χ α)): Parser χ (List α) := λ xs =>
  let rec loop (ys: Array α) (ps: List (Parser χ α)) (xs: List χ): Except (Error χ) (List α × List χ) :=
    match ps with
      | [] => Except.ok (ys.toList, xs)
      | p :: ps =>
        match p xs with
          | Except.error err => Except.error err
          | Except.ok (y, xs) =>
            loop (ys.push y) ps xs
  loop #[] ps xs


def pred (p: χ → Bool): Parser χ χ := λ xs =>
  match xs with
    | [] => err xs
    | x :: xs =>
      if p x then
        Except.ok (x, xs)
      else
        err xs

def exact [BEq χ] (y: χ): Parser χ χ := pred (· == y)

def exactList [BEq χ] (ys: List χ): Parser χ (List χ) :=
  Parser.transpose (ys.map exact)

#eval exactList "hehe".toList "hehea123".toList

namespace String

def toStringParser (p: Parser Char (List Char)): Parser Char String :=
  p.map (String.mk ·)

def whitespaceWeak : Parser Char String :=
  -- parse any whitespace
  -- empty whitespace is ok
  toStringParser (pred (·.isWhitespace)).many

def whiteSpaceWithoutNewLineWeak : Parser Char String :=
  -- parse any whitespace but not new line
  -- empty whitespace is ok
  toStringParser (pred (λ c => c.isWhitespace ∧ (¬ c = '\n'))).many


def whitespace : Parser Char String :=
  -- parse some whitespace
  -- empty whitespace is not ok
  toStringParser (pred (·.isWhitespace)).many1

def exact (ys: String): Parser Char String :=
  toStringParser (exactList ys.toList)

#eval whitespace "abc123  ".toList
#eval whitespace "   abc123".toList
#eval exact "let" "let123".toList


end String

end Parser.Combinator
