import EL.EL

def source := "
  (inductive {Bool: U_2}
    {true: Bool}
    {false: Bool}
  )

  (inductive {Nat: U_2}
    {zero: Nat}
    {succ: {{_: Nat} => Nat}}
  )

  (let one (succ zero))
  (let two (succ one))

  (let is_pos
    (lambda {n: Nat} (match n
      { (zero)   -> false}
      { (succ x) -> true}
    ))
  )

  (is_pos zero)
  (is_pos two)

  (inductive {{ {T: U_2} => (List T) } : U_2}
    {nil: (List T)}
    {cons: {{init: (List T)} {tail: T} => (List T)} }
  )


  {x := 3}
  {y := 4}

  {x + y * z}

"


-- #eval Util.parseAll EL.parse (EL.tokenize source)

#eval Form.defaultParser.tokenize "{}"
#eval Form.defaultParser.parse (Form.defaultParser.tokenize "{a b (1, 2) => x}")

def main : IO Unit := do

  let tokens := EL.tokenize source
  let result := Util.parseAll EL.parse tokens

  if result.remaining.length ≠ 0 then
    IO.println s!"error at {repr result.remaining}"
  else
    IO.println s!"{repr result.items}"
