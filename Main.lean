import EL.EL

def source := "
  (let Nat (inh U_2))
  (let n0 (inh Nat))
  (let succ (inh (lambda (: _ Nat) Nat)))

  (let n1 (succ n0))
  (let n2 (succ n0))

  (let x 3)
  (let y 4)

  (+ x y)

  (inductive (: Nat U_2)
    (: zero (lambda Nat))
    (: succ (lambda (: _ Nat) Nat))
  )

  (let one (succ zero))
  (let two (succ one))

  (let is_pos
    (lambda (: n Nat) (match n
      (case (Nat.zero) false)
      (case (Nat.succ x) true)
    ))
  )
"

-- #eval Util.parseAll EL.parse (EL.tokenize source)

def main : IO Unit := do

  let tokens := EL.tokenize source
  let result := Util.parseAll EL.parse tokens

  if result.remaining.length ≠ 0 then
    IO.println s!"error at {repr result.remaining}"
  else
    IO.println s!"{repr result.items}"
