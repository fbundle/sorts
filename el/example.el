[
    Bool U_0                undefined
    True Bool               undefined
    False Bool              undefined

    Nat U_0                 undefined
    0 Nat                   undefined
    succ {Nat -> Nat}       undefined

    1 Nat (succ 0)
    2 Nat (succ 1)
    3 Nat (succ 2)

    x Any "example string with \"quotes\" and spaces"
    x Any {1 + 2 + 3}
    x Any {1 ⊕ 2 ⊕ 3}

    is_two {Nat -> Bool} {x => (match x
        (succ 1) True
        False
    )}

    (is_two 3)
]