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

  {one := (succ zero)}
  {two := (succ one)}

  {is_pos :=
    {{n: Nat} => (match n
      { (zero)   -> false}
      { (succ x) -> true}
    )}
  }

  (is_pos zero)
  (is_pos two)




  (inductive [{T: U_2} => (List T) : U_2]

  )


  {x := 3}
  {y := 4}

  {x => y => z}

"




def main : IO Unit := do

  let tokens := EL.tokenize source
  let result := Util.parseAll EL.parse tokens

  if result.remaining.length ≠ 0 then
    IO.println s!"{repr result.items}
    error at {repr result.remaining}"
  else
    IO.println s!"{repr result.items}"
