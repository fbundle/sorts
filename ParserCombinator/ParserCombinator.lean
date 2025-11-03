namespace ParserCombinator

def Parser α := List Char → Option (List Char × α)

def Parser.mapOption  (p: Parser α) (f: α → Option β): Parser β := λ tokens => do
  let (tokens, a) ← p tokens
  let b ← f a
  pure (tokens, b)

def Parser.map (p: Parser α) (f: α → β): Parser β := p.mapOption (λ a => some (f a))

def Parser.concat (p1: Parser α) (p2: Parser β): Parser (α × β) := λ tokens => do
  let (tokens, a) ← p1 tokens
  let (tokens, b) ← p2 tokens
  (tokens, (a, b))

infixr: 60 " ++ " => Parser.concat

def Parser.either (p1: Parser α) (p2: Parser α): Parser α := λ tokens => do
  match p1 tokens with
    | some (rest, a) => some (rest, a)
    | none => p2 tokens

infixr: 50 " || " => Parser.either -- lower precedence than concat

partial def Parser.many (p: Parser α): Parser (List α) :=
  let rec loop (acc: Array α) (tokens: List Char): Option (List Char × List α) :=
    match p tokens with
      | none => some (tokens, acc.toList)
      | some (rest, e) => loop (acc.push e) rest
  loop #[]

-- TOOLS

def parseEmpty: Parser Unit := λ tokens => some (tokens, ())
def parseFail: Parser α := λ _ => none

def parseSingle: Parser Char := λ tokens =>
  match tokens with
    | [] => none
    | head :: rest => some (rest, head)

def parseExact (char: Char): Parser Unit :=
  parseSingle.mapOption (λ head =>
    if head = char then
      some ()
    else
      none
  )

def parseExactString (s: String): Parser Unit :=
  (s.toList.map parseExact).foldl (λ p1 p2 =>
    (p1 ++ p2).map (λ _ => ())
  ) parseEmpty

#eval parseExactString "hehe" "heheh123".toList


end ParserCombinator
