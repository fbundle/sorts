namespace ParserCombinator

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
      | some (a, xs) => loop (as.push a) xs
  loop #[] xs


end ParserCombinator
