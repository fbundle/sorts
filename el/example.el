#infix_left ⊕ ⊗
#infix_right => -> : , = :=

(let
    (let Bool U_0)
    (let True Bool)
    (let False Bool)

    (let Nat U_0)
    (let 0 Nat)
    (let succ {Nat -> Nat})

    (let 1 Nat (succ 0))
    (let 2 Nat (succ 1))
    (let 3 Nat (succ 2))

    (let is_two {Nat -> Bool} {x => (match x
        (succ 1) True
        False
    )})

    (is_two 3)
)