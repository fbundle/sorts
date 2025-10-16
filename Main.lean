import EL.EL

def source := "
  (inductive (: (lambda (Bool)) U_2)
    (: true (lambda (Bool)))
    (: false (lambda (Bool)))
  )

  (inductive (: (lambda (Nat)) U_2)
    (: zero (lambda (Nat)))
    (: succ (lambda (: _ (Nat)) (Nat)))
  )

  (let one (succ zero))
  (let two (succ one))

  (let is_pos
    (lambda (: n (Nat)) (match n
      (case (zero) false)
      (case (succ x) true)
    ))
  )

  (is_pos zero)
  (is_pos two)

  (inductive (: (lambda (: T U_2) (List T)) U_2)
    (: nil (lambda (List T)))
    (: cons (lambda (: init (List T)) (: tail T) (List T)))
  )


  (let x 3)
  (let y 4)

  (+ x y)

"

-- #eval Util.parseAll EL.parse (EL.tokenize source)

def main : IO Unit := do

  let tokens := EL.tokenize source
  let result := Util.parseAll EL.parse tokens

  if result.remaining.length ≠ 0 then
    IO.println s!"error at {repr result.remaining}"
  else
    IO.println s!"{repr result.items}"
