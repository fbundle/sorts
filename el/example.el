(let
    Bool U_0                undef
    True Bool               undef
    False Bool              undef

    Nat U_0                 undef
    0 Nat                   undef
    succ {Nat -> Nat}       undef

    1 Nat (succ 0)
    2 Nat (succ 1)
    3 Nat (succ 2)

    Any U_0                 undef
    x Any "example string with \"quotes\" and spaces"
    x Any {1 ⊕ 2 ⊕ 3}
    x Any {1 ⊗ 2 ⊗ 3 ⊗ 4}

    is_two {Nat -> Bool} {x => (match x
        (succ 1) True
        False
    )}

    add {Nat -> Nat -> Nat} {x => {y => (match x
        (succ z) (succ ((add z) y))
        y
    )}}

    (is_two 3)
)