import EL.EL

def source := "
  (let Nat (*U_2))
  (let n0 (*Nat))
  (let succ (*(-> Nat)))

  (let n1 (succ n0))
  (let n2 (succ n0))

  (let x 3)
  (let y 4)

  (+ x y)

  (ind (: Nat U_2)
    (: zero (=> Nat))
    (: succ (=> (: _ Nat) Nat))
  )

  (let one (succ zero))
  (let two (succ one))

  (let is_pos
    (=> (: n Nat) (match n
      (-> (Nat.zero) false)
      (-> (Nat.succ x) true)
    ))
  )
"

#eval Util.parseAll EL.parse (EL.tokenize source)

def main : IO Unit := do
  let tokens := EL.tokenize source
  let x := Util.parseAll EL.parse tokens

  IO.println s!"{repr x}"
